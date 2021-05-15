package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	nickHelix "github.com/nicklaw5/helix"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type userAcessTokenData struct {
	AccessToken  string
	RefreshToken string
	Scope        string
}

type tokenClaims struct {
	UserID         string
	StandardClaims jwt.StandardClaims
}

func (t *tokenClaims) Valid() error {
	return nil
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	resp, err := s.helixUserClient.Client.RequestUserAccessToken(code)
	if err != nil || resp.StatusCode >= 400 {
		log.Errorf("failed to request userAccessToken: %s %v", err, resp)
		s.dashboardRedirect(w, r, http.StatusBadRequest, "")
		return
	}

	// validate
	success, validateResp, err := s.helixClient.Client.ValidateToken(resp.Data.AccessToken)
	if !success || err != nil {
		log.Errorf("failed to veryify new Token %s", err)
		s.dashboardRedirect(w, r, http.StatusInternalServerError, "")
	}

	token, err := s.createApiToken(validateResp.Data.UserID)
	if err != nil {
		log.Errorf("failed to create jwt token in callback %s", err)
		s.dashboardRedirect(w, r, http.StatusInternalServerError, "")
	}

	marshalled, err := json.Marshal(userAcessTokenData{resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " ")})
	if err != nil {
		log.Errorf("failed to marshal userAcessToken in callback %s", err)
		s.dashboardRedirect(w, r, http.StatusInternalServerError, "")
	}

	err = s.store.Client.HSet("userAccessTokensData", validateResp.Data.UserID, marshalled).Err()
	if err != nil {
		log.Errorf("failed to set userAccessToken in callback: %s", err)
		s.dashboardRedirect(w, r, http.StatusInternalServerError, "")
	}

	s.dashboardRedirect(w, r, http.StatusOK, token)
}

func (s *Server) createApiToken(userID string) (string, error) {
	expirationTime := time.Now().Add(365 * 24 * time.Hour)
	claims := &tokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.Secret))

	return tokenString, err
}

func (s *Server) dashboardRedirect(w http.ResponseWriter, r *http.Request, status int, scToken string) {
	params := url.Values{
		"result": {fmt.Sprint(status)},
	}

	cookie := http.Cookie{
		Name:    "scToken",
		Value:   scToken,
		Domain:  s.cfg.CookieDomain,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, s.cfg.WebBaseUrl+"/dashboard"+"?"+params.Encode(), http.StatusFound)
}

func (s *Server) getUserAccessToken(userID string) (userAcessTokenData, error) {
	val, err := s.store.Client.HGet("userAccessTokensData", userID).Result()
	if err != nil {
		return userAcessTokenData{}, err
	}

	var token userAcessTokenData
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		return userAcessTokenData{}, err
	}

	return token, nil
}

func (s *Server) authenticate(r *http.Request) (bool, *nickHelix.ValidateTokenResponse, *userAcessTokenData) {
	scToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// Initialize a new instance of `Claims`
	claims := &tokenClaims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(scToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Secret), nil
	})
	if err != nil || !tkn.Valid {
		log.Errorf("found to validate jwt: %s", err)
		return false, nil, nil
	}

	token, err := s.getUserAccessToken(claims.UserID)
	if err != nil {
		log.Errorf("Failed to get userAccessTokenData: %s", err.Error())
		return false, nil, nil
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
				return false, nil, nil
			}

			success, resp, err = s.helixClient.Client.ValidateToken(refreshResp.AccessToken)
			if !success || err != nil {
				if err != nil {
					log.Errorf("refreshed Token did not validate: %s", err)
				}

				return success, resp, refreshResp
			}
		}

		return false, nil, nil
	}

	return success, resp, &token
}

func (s *Server) refreshToken(userID string, token userAcessTokenData) (bool, *userAcessTokenData) {
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

	err = s.store.Client.HSet("userAccessTokensData", userID, marshalled).Err()
	if err != nil {
		log.Errorf("failed to set userAccessTokenData to redis: %s", err)
		return false, nil
	}

	return true, &newToken
}

func (s *Server) tokenRefreshRoutine() {
	for {
		time.Sleep(time.Hour)

		tokens, err := s.store.Client.HGetAll("userAccessTokensData").Result()
		if err != nil {
			log.Errorf("tried refreshing tokens: %s", err)
			continue
		}

		log.Infof("starting refresh of %d tokens", len(tokens))

		for userID, tokenDataString := range tokens {
			var userToken userAcessTokenData
			if err := json.Unmarshal([]byte(tokenDataString), &userToken); err != nil {
				log.Errorf("failed unmarshal userAccessTokenData in tokenRefreshRoutine %s", err)
				continue
			}

			s.refreshToken(userID, userToken)
			time.Sleep(time.Millisecond * 500)
		}

		log.Infof("finished refresh of %d tokens", len(tokens))
	}
}
