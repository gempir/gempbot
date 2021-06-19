package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/bitraft/pkg/store"
	"github.com/labstack/echo/v4"
	nickHelix "github.com/nicklaw5/helix"

	"github.com/dgrijalva/jwt-go"
	"github.com/gempir/bitraft/pkg/log"
)

type tokenClaims struct {
	UserID         string
	StandardClaims jwt.StandardClaims
}

func (t *tokenClaims) Valid() error {
	return nil
}

func (s *Server) handleCallback(c echo.Context) error {
	code := c.QueryParam("code")

	resp, err := s.helixClient.Client.RequestUserAccessToken(code)
	if err != nil || resp.StatusCode >= 400 {
		log.Errorf("failed to request userAccessToken: %s %s", err, resp.ErrorMessage)
		// @TODO redirect to somewhere better
		return s.dashboardRedirect(c, "")
	}

	// validate
	success, validateResp, err := s.helixClient.Client.ValidateToken(resp.Data.AccessToken)
	if !success || err != nil {
		return fmt.Errorf("failed to veryify new Token %s", err)
	}

	token, err := s.createApiToken(validateResp.Data.UserID)
	if err != nil {
		return fmt.Errorf("failed to create jwt token in callback %s", err)
	}

	err = s.db.SaveUserAccessToken(validateResp.Data.UserID, resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "))
	if err != nil {
		return fmt.Errorf("failed to set userAccessToken in callback: %s", err)
	}

	err = s.db.SaveBotConfig(store.BotConfig{OwnerTwitchID: validateResp.Data.UserID, JoinBot: true})
	if err != nil {
		log.Error(err)
	}
	s.store.PublishIngesterMessage(store.IngesterMsgJoin, validateResp.Data.Login)

	go s.subscribePredictions(validateResp.Data.UserID)

	return s.dashboardRedirect(c, token)
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

func (s *Server) dashboardRedirect(c echo.Context, scToken string) error {
	cookie := http.Cookie{
		Name:    "scToken",
		Value:   scToken,
		Domain:  s.cfg.CookieDomain,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	}

	c.SetCookie(&cookie)
	err := c.Redirect(http.StatusFound, s.cfg.WebBaseUrl+"/dashboard")
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *Server) authenticate(c echo.Context) (nickHelix.ValidateTokenResponse, store.UserAccessToken, error) {
	scToken := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")

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
		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "bad authentication")
	}

	token, err := s.db.GetUserAccessToken(claims.UserID)
	if err != nil {
		log.Errorf("Failed to get userAccessTokenData: %s", err.Error())
		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "Failed to get userAccessTokenData: %s", err.Error())
	}

	success, resp, err := s.helixClient.Client.ValidateToken(token.AccessToken)
	if !success || err != nil {
		if err != nil {
			log.Errorf("token did not validate: %s", err)
		}

		// Token might be expired, let's try refreshing
		if resp.Error == "Unauthorized" {
			err := s.refreshToken(token)
			if err != nil {
				return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "failed to refresh token")
			}

			refreshedToken, err := s.db.GetUserAccessToken(claims.UserID)
			if err != nil {
				log.Errorf("Failed to get userAccessTokenData: %s", err.Error())
				return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "Failed to get userAccessTokenData: %s", err.Error())
			}

			success, resp, err = s.helixClient.Client.ValidateToken(refreshedToken.AccessToken)
			if !success || err != nil {
				if err != nil {
					log.Errorf("refreshed Token did not validate: %s", err)
				}

				return nickHelix.ValidateTokenResponse{}, refreshedToken, echo.NewHTTPError(http.StatusUnauthorized, "refreshed token did not validate")
			}

			return *resp, refreshedToken, nil
		}

		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "token not valid")
	}

	return *resp, token, nil
}

func (s *Server) refreshToken(token store.UserAccessToken) error {
	resp, err := s.helixClient.Client.RefreshUserAccessToken(token.RefreshToken)
	if err != nil {
		return err
	}

	err = s.db.SaveUserAccessToken(token.OwnerTwitchID, resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) tokenRefreshRoutine() {
	for {
		time.Sleep(time.Hour)

		tokens := s.db.GetAllUserAccessToken()

		log.Infof("starting refresh of %d tokens", len(tokens))

		for _, token := range tokens {
			err := s.refreshToken(token)
			if err != nil {
				log.Errorf("failed to refresh token %s", err)
			}
			time.Sleep(time.Millisecond * 500)
		}

		log.Infof("finished refresh of %d tokens", len(tokens))
	}
}
