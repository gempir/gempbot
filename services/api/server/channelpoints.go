package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gempir/bitraft/pkg/log"
	"github.com/gempir/bitraft/pkg/slice"
	"github.com/labstack/echo/v4"
	nickHelix "github.com/nicklaw5/helix"
)

type Redemptions struct {
	Bttv Redemption
}

type Redemption struct {
	Title  string
	Active bool
}

const (
	bttvPrompt = "Add a BetterTTV emote! In the text field, send a link to the BetterTTV emote. powered by bitraft.gempir.com"
)

type channelPointRedemption struct {
	Subscription struct {
		ID        string `json:"id"`
		Status    string `json:"status"`
		Type      string `json:"type"`
		Version   string `json:"version"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
			RewardID          string `json:"reward_id"`
		} `json:"condition"`
		Transport struct {
			Method   string `json:"method"`
			Callback string `json:"callback"`
		} `json:"transport"`
		CreatedAt time.Time `json:"created_at"`
		Cost      int       `json:"cost"`
	} `json:"subscription"`
	Event struct {
		BroadcasterUserID    string    `json:"broadcaster_user_id"`
		BroadcasterUserLogin string    `json:"broadcaster_user_login"`
		BroadcasterUserName  string    `json:"broadcaster_user_name"`
		ID                   string    `json:"id"`
		UserID               string    `json:"user_id"`
		UserLogin            string    `json:"user_login"`
		UserName             string    `json:"user_name"`
		UserInput            string    `json:"user_input"`
		Status               string    `json:"status"`
		RedeemedAt           time.Time `json:"redeemed_at"`
		Reward               struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			Prompt string `json:"prompt"`
			Cost   int    `json:"cost"`
		} `json:"reward"`
	} `json:"event"`
}

var bttvRegex = regexp.MustCompile(`https?:\/\/betterttv.com\/emotes\/(\w*)`)

func (s *Server) subscribeChannelPoints(userID string) error {
	response, err := s.helixClient.CreateChannelPointsRewardAdd(userID, s.cfg.WebhookApiBaseUrl+"/api/redemption")
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return nil
	}

	if response.StatusCode == http.StatusForbidden {
		return errors.New("forbidden")
	}

	log.Infof("[%d] created subscription %s", response.StatusCode, response.Error)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new subscription for %s id: %s", userID, sub.ID)
		s.db.AddEventSubSubscription(userID, sub.ID, sub.Version)
	}

	return nil
}

func (s *Server) removeEventSubSubscription(userID string, subscriptionID string, reason string) error {
	response, err := s.helixClient.Client.RemoveEventSubSubscription(subscriptionID)
	if err != nil {
		return err
	}

	log.Infof("[%d] removed EventSubSubscription for %s reason: %s", response.StatusCode, userID, reason)
	s.db.RemoveEventSubSubscription(userID, subscriptionID)

	return nil
}

func (s *Server) syncSubscriptions() {
	resp, err := s.helixClient.Client.GetEventSubSubscriptions(&nickHelix.EventSubSubscriptionsParams{})
	if err != nil {
		log.Errorf("Failed to get subscriptions: %s", err)
		return
	}

	log.Infof("Found %d total subscriptions, syncing to DB", resp.Data.Total)
	subscribed := []string{}

	for _, sub := range resp.Data.EventSubSubscriptions {
		if !strings.Contains(sub.Transport.Callback, s.cfg.WebhookApiBaseUrl) || sub.Status == nickHelix.EventSubStatusFailed {
			err := s.removeEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, "bad EventSub subscription, unsubscribing")
			if err != nil {
				log.Errorf("Failed to unsubscribe %s error: %s", sub.Condition.BroadcasterUserID, err.Error())
			}
			continue
		}

		_, err = s.db.GetEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID)
		if err != nil {
			log.Infof("Found unknown subscription, adding %s", sub.Condition.BroadcasterUserID)
			s.db.AddEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, sub.Version)
		}

		subscribed = append(subscribed, sub.Condition.BroadcasterUserID)
	}

	rewards := s.db.GetDistinctRewardsPerUser()
	log.Infof("Found %d total distinct rewards, checking missing subscriptions", len(rewards))

	for _, dbReward := range rewards {
		if !slice.Contains(subscribed, dbReward.OwnerTwitchID) {
			log.Infof("Found no subscription for existing reward, creating subscription %s", dbReward.OwnerTwitchID)
			err := s.subscribeChannelPoints(dbReward.OwnerTwitchID)
			if err != nil {
				log.Infof("Removing reward for user %s because we didn't get permission to subscribe eventsub", dbReward.OwnerTwitchID)
				s.db.DeleteChannelPointReward(dbReward.OwnerTwitchID, dbReward.Type)
			}
		}
	}
}

func (s *Server) handleChannelPointsRedemption(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed reading body")
	}

	verified := nickHelix.VerifyEventSubNotification(s.cfg.Secret, c.Request().Header, string(body))
	if !verified {
		log.Errorf("Failed verification %s", c.Request().Header.Get("Twitch-Eventsub-Message-Id"))
		return echo.NewHTTPError(http.StatusPreconditionFailed, "failed verfication")
	}

	if c.Request().Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		return s.handleChallenge(c, body)
	}

	var redemption channelPointRedemption
	err = json.Unmarshal(body, &redemption)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed decoding body: "+err.Error())
	}

	err = s.handleRedemption(redemption)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.String(http.StatusOK, "success")
}

func (s *Server) handleRedemption(redemption channelPointRedemption) error {
	exists, err := s.store.Client.Exists("redemption:" + redemption.Event.ID).Result()
	if err != nil {
		log.Error(err)
		return err
	}
	if exists == 1 {
		log.Infof("Redemption already handled before, ignoring retry %s", redemption.Event.ID)
		return nil
	}

	s.store.Client.SetNX("redemption:"+redemption.Event.ID, "1", time.Minute*10)

	if redemption.Subscription.Version != "1" {
		log.Errorf("Unknown subscription version found %s %s", redemption.Subscription.Version, redemption.Subscription.ID)
		return nil
	}

	reward, err := s.db.GetEnabledChannelPointRewardByID(redemption.Event.Reward.ID)
	if err != nil {
		// no redemption found
		return nil
	}

	// Err is only returned when it's worth responding with a bad response code
	if reward.Type == TYPE_BTTV {
		err = s.handleBttvRedemption(redemption)
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleBttvRedemption(redemption channelPointRedemption) error {
	success := false

	matches := bttvRegex.FindAllStringSubmatch(redemption.Event.UserInput, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		emoteAdded, emoteRemoved, err := s.emotechief.SetEmote(redemption.Event.BroadcasterUserID, matches[0][1], redemption.Event.BroadcasterUserLogin)
		if err != nil {
			log.Warnf("Bttv error %s %s", redemption.Event.BroadcasterUserLogin, err)
			s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add emote from: @%s error: %s", redemption.Event.UserName, err.Error()))
		} else if emoteAdded != nil && emoteRemoved != nil {
			success = true
			s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserLogin, fmt.Sprintf("✅ Added new emote: %s redeemed by @%s removed: %s", emoteAdded.Code, redemption.Event.UserName, emoteRemoved.Code))
		} else if emoteAdded != nil {
			success = true
			s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserLogin, fmt.Sprintf("✅ Added new emote: %s redeemed by @%s", emoteAdded.Code, redemption.Event.UserName))
		} else {
			success = true
			s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserLogin, fmt.Sprintf("✅ Added new emote: [unknown] redeemed by @%s", redemption.Event.UserName))
		}
	} else {
		s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add emote from @%s error: no bttv link found in message", redemption.Event.UserName))
	}

	token, err := s.db.GetUserAccessToken(redemption.Event.BroadcasterUserID)
	if err != nil {
		log.Errorf("Failed to get userAccess token to update redemption status for %s", redemption.Event.BroadcasterUserID)
		return nil
	} else {
		err := s.helixClient.UpdateRedemptionStatus(redemption.Event.BroadcasterUserID, token.AccessToken, redemption.Event.Reward.ID, redemption.Event.ID, success)
		if err != nil {
			log.Errorf("Failed to update redemption status %s", err.Error())
			return nil
		}
	}

	return nil
}

func (s *Server) handleChallenge(c echo.Context, body []byte) error {
	var event struct {
		Challenge string `json:"challenge"`
	}
	err := json.Unmarshal(body, &event)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to handle challenge: "+err.Error(), http.StatusBadRequest))
	}

	log.Infof("Challenge success: %s", event.Challenge)
	return c.String(http.StatusOK, event.Challenge)
}
