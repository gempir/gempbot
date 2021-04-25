package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

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

	var event helix.EventSubChannelPointsCustomRewardRedemptionEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.emotechief.SetEmote(event.BroadcasterUserID, event.UserInput)
	log.Info(event.UserInput)
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
