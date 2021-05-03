package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-redis/redis/v7"
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
	ok, auth, _ := s.authenticate(r)
	if !ok {
		http.Error(w, "bad authentication", http.StatusUnauthorized)
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

		newConfig := false

		_, err = s.store.Client.HGet("userConfig", auth.Data.UserID).Result()
		if err == redis.Nil {
			newConfig = true
		} else if err != nil {
			log.Error(err)
		}

		_, err = s.store.Client.HSet("userConfig", auth.Data.UserID, body).Result()
		if err != nil {
			log.Error(err)
			http.Error(w, "Failed updating: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if newConfig {
			log.Info("Created new config for: ", auth.Data.Login)
			s.subscribeChannelPoints(auth.Data.UserID)
		}

		writeJSON(w, "", http.StatusOK)
	} else if r.Method == http.MethodDelete {
		_, err := s.store.Client.HDel("userConfig", auth.Data.UserID).Result()
		if err != nil {
			log.Error(err)
			http.Error(w, "Failed deleting: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.unsubscribeChannelPoints(auth.Data.UserID)
		if err != nil {
			log.Error(err)
			http.Error(w, "Failed to unsubscribe"+err.Error(), http.StatusInternalServerError)
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
	_, err = w.Write(js)
	if err != nil {
		log.Errorf("Faile to writeJSON: %s", err)
	}
}
