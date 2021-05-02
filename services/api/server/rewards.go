package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *Server) handleRewards(w http.ResponseWriter, r *http.Request) {
	ok, auth, token := s.authenticate(r)
	if !ok {
		http.Error(w, "bad authentication", http.StatusForbidden)
		return
	}

	rewards, err := s.helixClient.GetRewards(auth.Data.UserID, token.AccessToken)
	if err != nil {
		log.Errorf("Failed to get rewards: %s", err)
	}

	writeJSON(w, rewards, http.StatusOK)
}
