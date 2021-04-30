package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

func createDefaultUserConfig() UserConfig {
	return UserConfig{
		Redemptions: Redemptions{
			Bttv: Redemption{Title: "Bttv emote", Active: false},
		},
	}
}

func (s *Server) handleUserConfig(w http.ResponseWriter, r *http.Request) {
	ok, auth := s.authenticate(r)
	if !ok {
		http.Error(w, "bad authentication", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodGet {
		val, err := s.store.Client.HGet("userConfig", auth.Data.UserID).Result()
		if err != nil {
			writeJSON(w, createDefaultUserConfig(), http.StatusOK)
			return
		}
		var userConfig UserConfig
		if err := json.Unmarshal([]byte(val), &userConfig); err != nil {
			log.Error(err)
			return
		}

		writeJSON(w, userConfig, http.StatusOK)
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Failed reading update body: %s", err)
			http.Error(w, "failed reading body"+err.Error(), http.StatusBadRequest)
			return
		}

		_, err = s.store.Client.HSet("userConfig", auth.Data.UserID, body).Result()
		if err != nil {
			log.Error(err)
			http.Error(w, "Failed updating: "+err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, "", http.StatusOK)
	}

}

func writeJSON(w http.ResponseWriter, data interface{}, code int) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(js)
}

func (s *Server) authenticate(r *http.Request) (bool, *helix.ValidateTokenResponse) {
	token := r.Header.Get("accessToken")
	if token == "" {
		return false, nil
	}

	success, resp, err := s.helixClient.Client.ValidateToken(token)
	if err != nil {
		log.Errorf("Failed to authenticate: %s", err)
	}

	return success, resp
}
