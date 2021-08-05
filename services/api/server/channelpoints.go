package server

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gempir/bitraft/pkg/dto"
	"github.com/gempir/bitraft/pkg/log"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	nickHelix "github.com/nicklaw5/helix"
)

type Redemptions struct {
	Bttv Redemption
}

type Redemption struct {
	Title  string
	Active bool
}

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
	response, err := s.helixClient.CreateEventSubSubscription(userID, s.cfg.WebhookApiBaseUrl+"/api/redemption", nickHelix.EventSubTypeChannelPointsCustomRewardRedemptionAdd)
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return nil
	}

	if response.StatusCode == http.StatusForbidden {
		return errors.New("forbidden")
	}

	log.Infof("[%d] subscription %s %s", response.StatusCode, response.Error, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new subscription for %s id: %s", userID, sub.ID)
		s.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type)
	}

	return nil
}

func (s *Server) removeEventSubSubscription(userID string, subscriptionID string, subType string, reason string) error {
	response, err := s.helixClient.Client.RemoveEventSubSubscription(subscriptionID)
	if err != nil {
		return err
	}

	log.Infof("[%d] removed EventSubSubscription for %s reason: %s", response.StatusCode, userID, reason)
	s.db.RemoveEventSubSubscription(userID, subscriptionID, subType)

	return nil
}

func (s *Server) handleChannelPointsRedemption(c echo.Context) error {
	var redemption channelPointRedemption
	done, err := s.handleWebhook(c, &redemption)
	if err != nil || done {
		return err
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

	if redemption.Subscription.Version != "1" && redemption.Subscription.Version != "" {
		log.Errorf("Unknown subscription version found %s %s", redemption.Subscription.Version, redemption.Subscription.ID)
		return nil
	}

	reward, err := s.db.GetEnabledChannelPointRewardByID(redemption.Event.Reward.ID)
	if err != nil {
		// no redemption found
		return nil
	}

	// Err is only returned when it's worth responding with a bad response code
	if reward.Type == dto.REWARD_BTTV {
		err = s.handleBttvRedemption(reward, redemption)
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleBttvRedemption(reward store.ChannelPointReward, redemption channelPointRedemption) error {
	var opts BttvAdditionalOptions
	err := mapstructure.Decode(reward.AdditionalOptions, &opts)
	if err != nil {
		log.Errorf("Error decoding additional options: %s", err)
		s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserID, redemption.Event.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add emote from: @%s error: %s", redemption.Event.UserName, "failed to read reward options from broadcaster"))
		return nil
	}
	success := false

	matches := bttvRegex.FindAllStringSubmatch(redemption.Event.UserInput, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		emoteAdded, emoteRemoved, err := s.emotechief.SetEmote(redemption.Event.BroadcasterUserID, matches[0][1], redemption.Event.BroadcasterUserLogin, opts.Slots)
		if err != nil {
			log.Warnf("Bttv error %s %s", redemption.Event.BroadcasterUserLogin, err)
			s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserID, redemption.Event.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add emote from: @%s error: %s", redemption.Event.UserName, err.Error()))
		} else if emoteAdded != nil && emoteRemoved != nil {
			success = true
			s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserID, redemption.Event.BroadcasterUserLogin, fmt.Sprintf("✅ Added new emote: %s redeemed by @%s removed: %s", emoteAdded.Code, redemption.Event.UserName, emoteRemoved.Code))
		} else if emoteAdded != nil {
			success = true
			s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserID, redemption.Event.BroadcasterUserLogin, fmt.Sprintf("✅ Added new emote: %s redeemed by @%s", emoteAdded.Code, redemption.Event.UserName))
		} else {
			success = true
			s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserID, redemption.Event.BroadcasterUserLogin, fmt.Sprintf("✅ Added new emote: [unknown] redeemed by @%s", redemption.Event.UserName))
		}
	} else {
		s.store.PublishSpeakerMessage(redemption.Event.BroadcasterUserID, redemption.Event.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add emote from @%s error: no bttv link found in message", redemption.Event.UserName))
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
