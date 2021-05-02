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

	success, _, err := s.helixClient.Client.ValidateToken(data.AccessToken)
	if !success || err != nil {
		if err != nil {
			log.Error(err)
		}
		http.Error(w, "token did not validate", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "success")
}
