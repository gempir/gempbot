package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *Server) handleOauth(w http.ResponseWriter, r *http.Request) {
	var data struct {
		AccessToken string `json:"accessToken"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	success, resp, err := s.helixClient.Client.ValidateToken(data.AccessToken)
	if !success || err != nil {
		if err != nil {
			log.Error(err)
		}
		http.Error(w, "token did not validate", http.StatusBadRequest)
		return
	}

	s.store.Client.HSet("accessToken", resp.Data.UserID, data.AccessToken)
	log.Infof("Stored new accessToken for %s %v", resp.Data.Login, resp.Data.Scopes)

	// by default subscribe channelPoints
	sub, err := json.Marshal(&subscription{[]redemption{redemption{"bttv emote"}}})
	if err != nil {
		log.Infof("Failed to add subscription to redis %s", err)
	}
	err = s.store.Client.HSet("subscriptions", resp.Data.UserID, sub).Err()
	if err != nil {
		log.Infof("Failed to add subscription to redis %s", err)
	}
	go s.subscribeChannelPoints(resp.Data.UserID)

	fmt.Fprint(w, "success")
}
