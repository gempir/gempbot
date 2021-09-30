package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/golang-jwt/jwt"
	nickHelix "github.com/nicklaw5/helix/v2"
)

func CreateApiToken(secret, userID string) string {
	expirationTime := time.Now().Add(365 * 24 * time.Hour)
	claims := &TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString
}

type TokenClaims struct {
	UserID         string
	StandardClaims jwt.StandardClaims
}

func (t *TokenClaims) Valid() error {
	return nil
}

func NewAuth(cfg *config.Config, db *store.Database, helixClient *helix.Client) *Auth {
	return &Auth{
		cfg:         cfg,
		db:          db,
		helixClient: helixClient,
	}
}

type Auth struct {
	helixClient *helix.Client
	db          *store.Database
	cfg         *config.Config
}

func (a *Auth) AttemptAuth(r *http.Request, w http.ResponseWriter) (nickHelix.ValidateTokenResponse, store.UserAccessToken, api.Error) {
	resp, token, err := a.Authenticate(r)
	if err != nil {
		a.WriteDeleteCookieResponse(w, err)
		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, err
	}

	return resp, token, nil
}

func (a *Auth) Authenticate(r *http.Request) (nickHelix.ValidateTokenResponse, store.UserAccessToken, api.Error) {
	scToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// Initialize a new instance of `Claims`
	claims := &TokenClaims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(scToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.Secret), nil
	})
	if err != nil || !tkn.Valid {
		log.Errorf("found to validate jwt: %s", err)
		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("bad authentication"))
	}

	token, err := a.db.GetUserAccessToken(claims.UserID)
	if err != nil {
		log.Errorf("Failed to get userAccessTokenData: %s", err.Error())
		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("Failed to get userAccessTokenData: %s", err.Error()))
	}

	success, resp, err := a.helixClient.Client.ValidateToken(token.AccessToken)
	if !success || err != nil {
		if err != nil {
			log.Errorf("token did not validate: %s", err)
		}

		// Token might be expired, let's try refreshing
		if resp.Error == "Unauthorized" {
			err := a.refreshToken(r.Context(), token)
			if err != nil {
				return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("failed to refresh token"))
			}

			refreshedToken, err := a.db.GetUserAccessToken(claims.UserID)
			if err != nil {
				log.Errorf("Failed to get userAccessTokenData: %s", err.Error())
				return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("Failed to get userAccessTokenData: %s", err.Error()))
			}

			success, resp, err = a.helixClient.Client.ValidateToken(refreshedToken.AccessToken)
			if !success || err != nil {
				if err != nil {
					log.Errorf("refreshed Token did not validate: %s", err)
				}

				return nickHelix.ValidateTokenResponse{}, refreshedToken, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("refreshed token did not validate"))
			}

			return *resp, refreshedToken, nil
		}

		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("token not valid: %s", resp.ErrorMessage))
	}

	return *resp, token, nil
}

func (a *Auth) refreshToken(ctx context.Context, token store.UserAccessToken) error {
	resp, err := a.helixClient.Client.RefreshUserAccessToken(token.RefreshToken)
	if err != nil {
		return err
	}

	err = a.db.SaveUserAccessToken(ctx, token.OwnerTwitchID, resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "))
	if err != nil {
		return err
	}

	return nil
}

func (a *Auth) WriteDeleteCookieResponse(w http.ResponseWriter, err api.Error) {
	cookie := &http.Cookie{
		Name:     "scToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	http.Error(w, err.Error(), err.Status())
}

// func (a *Auth) getUserConfig(userID string) UserConfig {
// 	uCfg := createDefaultUserConfig()

// 	botConfig, err := s.db.GetBotConfig(userID)
// 	if err != nil {
// 		uCfg.BotJoin = false
// 	} else {
// 		uCfg.BotJoin = botConfig.JoinBot
// 	}

// 	uCfg.Protected.CurrentUserID = userID

// 	perms := s.db.GetChannelPermissions(userID)
// 	for _, perm := range perms {
// 		uCfg.Permissions[perm.TwitchID] = Permission{perm.Editor, perm.Prediction}
// 	}

// 	for _, perm := range s.db.GetUserPermissions(userID) {
// 		uCfg.Protected.EditorFor = append(uCfg.Protected.EditorFor, perm.ChannelTwitchId)
// 	}

// 	return uCfg
// }
