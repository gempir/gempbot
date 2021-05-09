package server

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *Server) handleRewards(w http.ResponseWriter, r *http.Request) {
	ok, auth, token := s.authenticate(r)
	if !ok {
		http.Error(w, "bad authentication", http.StatusUnauthorized)
		return
	}

	rewardsUserId := auth.Data.UserID
	userAccessToken := token.AccessToken

	managing := r.URL.Query().Get("managing")
	if managing != "" {
		userConfig, err, isNew := s.getUserConfig(auth.Data.UserID)
		if err != nil || isNew {
			http.Error(w, "no config found", http.StatusUnauthorized)
			return
		}

		ownerUserID, err := s.checkEditor(r, userConfig)
		if err != nil {
			http.Error(w, "not editor: "+err.Error(), http.StatusUnauthorized)
			return
		}

		rewardsUserId = ownerUserID

		val, err := s.store.Client.HGet("userAccessTokensData", rewardsUserId).Result()
		if err != nil {
			http.Error(w, "no access token found: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var ownerToken userAcessTokenData
		if err := json.Unmarshal([]byte(val), &ownerToken); err != nil {
			http.Error(w, "no access token found: "+err.Error(), http.StatusInternalServerError)
			return
		}

		userAccessToken = ownerToken.AccessToken
	}

	rewards, err := s.helixClient.GetRewards(rewardsUserId, userAccessToken)
	if err != nil {
		log.Errorf("Failed to get rewards: %s", err)
	}

	writeJSON(w, rewards, http.StatusOK)
}
