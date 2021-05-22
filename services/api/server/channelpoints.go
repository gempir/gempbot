package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gempir/spamchamp/pkg/helix"
	"github.com/gempir/spamchamp/pkg/slice"
	nickHelix "github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
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

func (s *Server) subscribeChannelPoints(userID string) {
	// Twitch doesn't need a user token here, always an app token eventhough the user has to authenticate beforehand.
	// Internally they check if the app token has authenticated users
	response, err := s.helixUserClient.Client.CreateEventSubSubscription(
		&nickHelix.EventSubSubscription{
			Condition: nickHelix.EventSubCondition{BroadcasterUserID: userID},
			Transport: nickHelix.EventSubTransport{Method: "webhook", Callback: s.cfg.WebhookApiBaseUrl + "/api/redemption", Secret: s.cfg.Secret},
			Type:      "channel.channel_points_custom_reward_redemption.add",
			Version:   "1",
		},
	)
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("[%d] New subscription for %s id: %s", response.StatusCode, userID, sub.ID)
		s.store.Client.HSet("subscriptions", userID, sub.ID)
	}
}

func (s *Server) unsubscribeChannelPoints(userID string, reason string) error {
	subId, err := s.store.Client.HGet("subscriptions", userID).Result()
	if err != nil {
		return err
	}

	return s.removeEventSubSubscription(userID, subId, reason)
}

func (s *Server) removeEventSubSubscription(userID string, subscriptionID string, reason string) error {
	response, err := s.helixUserClient.Client.RemoveEventSubSubscription(subscriptionID)
	if err != nil {
		return err
	}

	log.Infof("[%d] removed EventSubSubscription for %s reason: %s", response.StatusCode, userID, reason)

	_, err = s.store.Client.HDel("subscriptions", userID).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) syncSubscriptions() {
	resp, err := s.helixUserClient.Client.GetEventSubSubscriptions(&nickHelix.EventSubSubscriptionsParams{})
	if err != nil {
		log.Errorf("Failed to get subscriptions: %s", err)
		return
	}

	log.Infof("Found %d total subscriptions, syncing to Redis", resp.Data.Total)
	subscribed := []string{}

	for _, sub := range resp.Data.EventSubSubscriptions {
		exists, err := s.store.Client.HExists("userConfig", sub.Condition.BroadcasterUserID).Result()
		if err != nil {
			log.Errorf("Failed to get userConfig while syncing %s", err.Error())
			continue
		}
		if !exists {
			err := s.removeEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, "no userConfig found")
			if err != nil {
				log.Errorf("Failed to unsubscribe %s error: %s", sub.Condition.BroadcasterUserID, err.Error())
			}
			continue
		}

		if !strings.Contains(sub.Transport.Callback, s.cfg.WebhookApiBaseUrl) {
			err := s.removeEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, "unknown transport, unsubscribing")
			if err != nil {
				log.Errorf("Failed to unsubscribe %s error: %s", sub.Condition.BroadcasterUserID, err.Error())
			}
			continue
		}

		subscribed = append(subscribed, sub.Condition.BroadcasterUserID)
		s.store.Client.HSet("subscriptions", sub.Condition.BroadcasterUserID, sub.ID)
	}

	userConfigs, err := s.store.Client.HGetAll("userConfig").Result()
	if err != nil {
		log.Errorf("Failed to sync subscriptions with userConfig: %s", err)
		return
	}

	for userID := range userConfigs {
		if !slice.Contains(subscribed, userID) {
			log.Info("Found no subscription for existing userConfig, creating subscription")
			s.subscribeChannelPoints(userID)
		}
	}
}

func (s *Server) handleChannelPointsRedemption(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed reading body", http.StatusBadRequest)
	}

	verified := nickHelix.VerifyEventSubNotification(s.cfg.Secret, r.Header, string(body))
	if !verified {
		log.Errorf("Failed verification: %s", r.Header.Get("Twitch-Eventsub-Message-Id"))
		http.Error(w, "failed verfication", http.StatusPreconditionFailed)
	}

	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		s.handleChallenge(w, body)
		return
	}

	var redemption channelPointRedemption
	err = json.Unmarshal(body, &redemption)
	if err != nil {
		log.Error(err)
		http.Error(w, "Failed decoding body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// get active subscritions for this channel
	val, err := s.store.Client.HGet("userConfig", redemption.Event.BroadcasterUserID).Result()
	if err != nil {
		log.Errorf("Won't handle redemption, no userConfig found %s", err)
		return
	}
	var userCfg UserConfig
	if err := json.Unmarshal([]byte(val), &userCfg); err != nil {
		log.Error(err)
		return
	}

	if userCfg.Rewards.BttvReward != nil && userCfg.Rewards.BttvReward.Enabled && userCfg.Rewards.BttvReward.ID == redemption.Event.Reward.ID {
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

		token, err := s.getUserAccessToken(redemption.Event.BroadcasterUserID)
		if err != nil {
			log.Errorf("Failed to get userAccess token to update redemption status for %s", redemption.Event.BroadcasterUserID)
		} else {
			log.Info(token)
			err := s.helixUserClient.UpdateRedemptionStatus(redemption.Event.BroadcasterUserID, token.AccessToken, redemption.Event.Reward.ID, redemption.Event.ID, success)
			if err != nil {
				log.Errorf("Failed to update redemption status %s", err.Error())
			}
		}
	}

	fmt.Fprint(w, "success")
}

func (s *Server) handleChallenge(w http.ResponseWriter, body []byte) {
	var event struct {
		Challenge string `json:"challenge"`
	}
	err := json.Unmarshal(body, &event)
	if err != nil {
		log.Errorf("Failed to handle challenge: %s", err.Error())
		http.Error(w, "Failed to handle challenge: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Infof("Challenge success: %s", event.Challenge)
	fmt.Fprint(w, event.Challenge)
}

func (s *Server) createOrUpdateChannelPointReward(userID string, request BttvReward, rewardID string) (BttvReward, error) {
	token, err := s.getUserAccessToken(userID)
	if err != nil {
		return BttvReward{}, err
	}

	req := helix.CreateCustomRewardRequest{
		Title:                             request.Title,
		Prompt:                            bttvPrompt,
		Cost:                              request.Cost,
		IsEnabled:                         request.Enabled,
		BackgroundColor:                   request.Backgroundcolor,
		IsUserInputRequired:               true,
		ShouldRedemptionsSkipRequestQueue: false,
		IsMaxPerStreamEnabled:             false,
		IsMaxPerUserPerStreamEnabled:      false,
		IsGlobalCooldownEnabled:           false,
	}

	if request.MaxPerStream != 0 {
		req.IsMaxPerStreamEnabled = true
		req.MaxPerStream = request.MaxPerStream
	}

	if request.MaxPerUserPerStream != 0 {
		req.IsMaxPerUserPerStreamEnabled = true
		req.MaxPerUserPerStream = request.MaxPerUserPerStream
	}

	if request.GlobalCooldownSeconds != 0 {
		req.IsGlobalCooldownEnabled = true
		req.GlobalCoolDownSeconds = request.GlobalCooldownSeconds
	}

	resp, err := s.helixUserClient.CreateOrUpdateReward(userID, token.AccessToken, req, rewardID)
	if err != nil {
		return BttvReward{}, err
	}

	return BttvReward{
		Title:                             resp.Title,
		Prompt:                            resp.Prompt,
		Cost:                              resp.Cost,
		Backgroundcolor:                   resp.BackgroundColor,
		IsMaxPerStreamEnabled:             resp.MaxPerStreamSetting.IsEnabled,
		MaxPerStream:                      resp.MaxPerStreamSetting.MaxPerStream,
		IsMaxPerUserPerStreamEnabled:      resp.MaxPerUserPerStreamSetting.IsEnabled,
		MaxPerUserPerStream:               resp.MaxPerUserPerStreamSetting.MaxPerUserPerStream,
		IsUserInputRequired:               resp.IsUserInputRequired,
		IsGlobalCooldownEnabled:           resp.GlobalCooldownSetting.IsEnabled,
		GlobalCooldownSeconds:             resp.GlobalCooldownSetting.GlobalCooldownSeconds,
		ShouldRedemptionsSkipRequestQueue: resp.ShouldRedemptionsSkipRequestQueue,
		Enabled:                           resp.IsEnabled,
		ID:                                resp.ID,
	}, nil
}
