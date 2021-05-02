package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	nickHelix "github.com/nicklaw5/helix"

	log "github.com/sirupsen/logrus"
)

type userAcessTokenData struct {
	AccessToken  string
	RefreshToken string
	Scope        string
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	resp, err := s.helixUserClient.Client.RequestUserAccessToken(code)
	if err != nil || resp.StatusCode >= 400 {
		log.Error(err, resp)
		s.dashboardRedirect(w, r, http.StatusForbidden, "")
		return
	}

	marshalled, err := json.Marshal(userAcessTokenData{resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " ")})
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

func (s *Server) authenticate(r *http.Request) (bool, *nickHelix.ValidateTokenResponse) {
	scToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	val, err := s.store.Client.HGet("accessTokens", scToken).Result()
	if err != nil {
		log.Error(err)
		return false, nil
	}

	var token userAcessTokenData
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		log.Error(err)
		return false, nil
	}

	success, resp, err := s.helixClient.Client.ValidateToken(token.AccessToken)
	if !success || err != nil {
		if err != nil {
			log.Error(err)
		}

		// Token might be expired, let's try refreshing
		if resp.Error == "Unauthorized" {
			success, refreshResp := s.refreshToken(scToken, token)
			if !success {
				return false, nil
			}

			success, resp, err = s.helixClient.Client.ValidateToken(refreshResp.AccessToken)
			if !success || err != nil {
				if err != nil {
					log.Error(err)
				}

				return success, resp
			}
		}

		return false, nil
	}

	return success, resp
}

func (s *Server) refreshToken(scToken string, token userAcessTokenData) (bool, *userAcessTokenData) {
	resp, err := s.helixClient.Client.RefreshUserAccessToken(token.RefreshToken)
	if err != nil {
		log.Error(err)
		return false, nil
	}

	newToken := userAcessTokenData{resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " ")}
	marshalled, err := json.Marshal(newToken)
	if err != nil {
		log.Error(err)
		return false, nil
	}

	err = s.store.Client.HSet("accessTokens", scToken, marshalled).Err()
	if err != nil {
		log.Error(err)
		return false, nil
	}
	log.Infof("Refreshed token for %s", scToken)

	return true, &newToken
}
