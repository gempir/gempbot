package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gempir/spamchamp/pkg/humanize"

	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

var bttvRegex = regexp.MustCompile(`https?:\/\/betterttv.com\/emotes\/(\w*)`)

func (s *Server) subscribeChannelPoints() {
	// s.helixUserClient.Client.SetUserAccessToken(s.store.Client.HGet("accessToken", "77829817").Val())
	response, err := s.helixUserClient.Client.CreateEventSubSubscription(
		&helix.EventSubSubscription{
			Condition: helix.EventSubCondition{BroadcasterUserID: "77829817"},
			Transport: helix.EventSubTransport{Method: "webhook", Callback: s.cfg.HttpBaseUrl + "/api/redemption", Secret: s.cfg.Secret},
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
	// verified := helix.VerifyEventSubNotification(s.cfg.Secret, r.Header, "")
	// if !verified {
	// 	http.Error(w, "failed verfication", http.StatusPreconditionFailed)
	// }

	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		s.handleChallenge(w, r)
		return
	}
	log.Info(humanize.FormatRequest(r))

	var event helix.EventSubChannelPointsCustomRewardRedemptionEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matches := bttvRegex.FindAllString(event.UserInput, -1)
	log.Infof("%v", matches)
	for _, match := range matches {
		_ = s.emotechief.SetEmote(event.BroadcasterUserID, match)
		fmt.Fprint(w, "success")
		return
	}

	log.Warnf("Could not find emote in message: %s", event.UserInput)
	http.Error(w, "Could not find emote", http.StatusBadRequest)
}

func (s *Server) handleChallenge(w http.ResponseWriter, r *http.Request) {
	var event struct{ Challenge string }
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, event.Challenge)
}
