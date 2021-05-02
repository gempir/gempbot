package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

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
		log.Errorf("failed to request userAccessToken: %s %v", err, resp)
		s.dashboardRedirect(w, r, http.StatusForbidden, "")
		return
	}

	marshalled, err := json.Marshal(userAcessTokenData{resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " ")})
	if err != nil {
		log.Errorf("failed to marshal userAcessToken in callback %s", err)
		s.dashboardRedirect(w, r, http.StatusInternalServerError, "")
	}

	err = s.store.Client.HSet("userAccessTokens", code, marshalled).Err()
	if err != nil {
		log.Errorf("failed to set userAccessToken in callback: %s", err)
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

	val, err := s.store.Client.HGet("userAccessTokens", scToken).Result()
	if err != nil {
		log.Errorf("found no accessToken: %s", err)
		return false, nil
	}

	var token userAcessTokenData
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		log.Errorf("failed to unmarshal token: %s", err)
		return false, nil
	}

	success, resp, err := s.helixClient.Client.ValidateToken(token.AccessToken)
	if !success || err != nil {
		if err != nil {
			log.Errorf("token did not validate: %s", err)
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
					log.Errorf("refreshed Token did not validate: %s", err)
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
		log.Errorf("failed refresh userAcessToken: %s", err)
		return false, nil
	}

	newToken := userAcessTokenData{resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " ")}
	marshalled, err := json.Marshal(newToken)
	if err != nil {
		log.Errorf("failed marshal refreshed userAcessToken: %s", err)
		return false, nil
	}

	err = s.store.Client.HSet("userAccessTokens", scToken, marshalled).Err()
	if err != nil {
		log.Errorf("failed to set userAccessTokenData to redis: %s", err)
		return false, nil
	}
	log.Infof("refreshed a token")

	return true, &newToken
}

func (s *Server) tokenRefreshRoutine() {
	for {
		time.Sleep(time.Hour)

		tokens, err := s.store.Client.HGetAll("userAccessTokens").Result()
		if err != nil {
			log.Errorf("tried refreshing tokens: %s", err)
			continue
		}

		log.Infof("starting refresh of %d tokens", len(tokens))

		for scToken, tokenDataString := range tokens {
			var userToken userAcessTokenData
			if err := json.Unmarshal([]byte(tokenDataString), &userToken); err != nil {
				log.Errorf("failed unmarshal userAccessTokenData in tokenRefreshRoutine %s", err)
				continue
			}

			s.refreshToken(scToken, userToken)
			time.Sleep(time.Millisecond * 500)
		}

		log.Infof("finished refresh of %d tokens", len(tokens))
	}
}
