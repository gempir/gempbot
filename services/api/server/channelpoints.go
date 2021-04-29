package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
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

func (s *Server) subscribeChannelPoints(channelId string) {
	log.Infof("Subscribing webhooks for: %s", channelId)

	// Twitch doesn't need a user token here, always an app token eventhough the user has to authenticate beforehand.
	// Internally they check if the app token has authenticated users
	// s.helixUserClient.Client.SetUserAccessToken(s.store.Client.HGet("accessToken", "77829817").Val())
	response, err := s.helixUserClient.Client.CreateEventSubSubscription(
		&helix.EventSubSubscription{
			Condition: helix.EventSubCondition{BroadcasterUserID: channelId},
			Transport: helix.EventSubTransport{Method: "webhook", Callback: s.cfg.ApiBaseUrl + "/api/redemption", Secret: s.cfg.Secret},
			Type:      "channel.channel_points_custom_reward_redemption.add",
			Version:   "1",
		},
	)
	if err != nil {
		log.Error(err)
	}

	log.Info(response)
}

func (s *Server) handleChannelPointsRedemption(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed reading body", http.StatusBadRequest)
	}

	verified := helix.VerifyEventSubNotification(s.cfg.Secret, r.Header, string(body))
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

	matches := bttvRegex.FindAllStringSubmatch(redemption.Event.UserInput, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		err = s.emotechief.SetEmote(redemption.Event.BroadcasterUserID, matches[0][1], redemption.Event.BroadcasterUserLogin)
		if err != nil {
			log.Warn(err)
		}
		fmt.Fprint(w, "success")
		return
	}

	log.Warnf("Could not find emote in message: %s", redemption.Event.UserInput)
	http.Error(w, "Could not find emote", http.StatusBadRequest)
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
