package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	resp, err := s.helixUserClient.Client.RequestUserAccessToken(code)
	if err != nil || resp.StatusCode >= 400 {
		log.Error(err, resp)
		s.dashboardRedirect(w, r, http.StatusForbidden, "")
		return
	}

	marshalled, err := json.Marshal(resp.Data)
	if err != nil {
		log.Error(err)
		s.dashboardRedirect(w, r, http.StatusInternalServerError, "")
	}

	err = s.store.Client.HSet("accessTokens", code, marshalled).Err()
	if err != nil {
		log.Error(err)
		s.dashboardRedirect(w, r, http.StatusInternalServerError, "")
	}

	s.dashboardRedirect(w, r, http.StatusOK, code)
}

func (s *Server) dashboardRedirect(w http.ResponseWriter, r *http.Request, status int, scToken string) {
	params := url.Values{
		"result": {fmt.Sprint(status)},
	}
	if scToken != "" {
		params.Add("scToken", scToken)
	}

	http.Redirect(w, r, s.cfg.WebBaseUrl+"/dashboard"+"?"+params.Encode(), http.StatusFound)
}
