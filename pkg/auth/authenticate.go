package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/golang-jwt/jwt"
	"github.com/nicklaw5/helix/v2"
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

func NewAuth(cfg *config.Config, db *store.Database, helixClient *helixclient.Client) *Auth {
	return &Auth{
		cfg:         cfg,
		db:          db,
		helixClient: helixClient,
	}
}

type Auth struct {
	helixClient *helixclient.Client
	db          *store.Database
	cfg         *config.Config
}

func (a *Auth) AttemptAuth(r *http.Request, w http.ResponseWriter) (helix.ValidateTokenResponse, store.UserAccessToken, api.Error) {
	resp, token, err := a.Authenticate(r)
	if err != nil {
		a.WriteDeleteCookieResponse(w, err)
		return helix.ValidateTokenResponse{}, store.UserAccessToken{}, err
	}

	return resp, token, nil
}

func (a *Auth) Authenticate(r *http.Request) (helix.ValidateTokenResponse, store.UserAccessToken, api.Error) {
	scToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	for _, cookie := range r.Cookies() {
		if cookie.Name == "scToken" {
			scToken = cookie.Value
		}
	}

	if scToken == "" {
		return helix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("no scToken cookie set"))
	}

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
		return helix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("bad authentication"))
	}

	token, err := a.db.GetUserAccessToken(claims.UserID)
	if err != nil {
		log.Errorf("Failed to get userAccessTokenData: %s", err.Error())
		return helix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("Failed to get userAccessTokenData: %s", err.Error()))
	}

	success, resp, err := a.helixClient.Client.ValidateToken(token.AccessToken)
	if !success || err != nil {
		if err != nil {
			log.Errorf("token did not validate: %s", err)
		}

		// Token might be expired, let's try refreshing
		if resp.Error == "Unauthorized" {
			err := a.helixClient.RefreshToken(token)
			if err != nil {
				return helix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("failed to refresh token"))
			}

			refreshedToken, err := a.db.GetUserAccessToken(claims.UserID)
			if err != nil {
				log.Errorf("Failed to get userAccessTokenData: %s", err.Error())
				return helix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("Failed to get userAccessTokenData: %s", err.Error()))
			}

			success, resp, err = a.helixClient.Client.ValidateToken(refreshedToken.AccessToken)
			if !success || err != nil {
				if err != nil {
					log.Errorf("refreshed Token did not validate: %s", err)
				}

				return helix.ValidateTokenResponse{}, refreshedToken, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("refreshed token did not validate"))
			}

			return *resp, refreshedToken, nil
		}

		return helix.ValidateTokenResponse{}, store.UserAccessToken{}, api.NewApiError(http.StatusUnauthorized, fmt.Errorf("token not valid: %s", resp.ErrorMessage))
	}

	return *resp, token, nil
}

func (a *Auth) WriteDeleteCookieResponse(w http.ResponseWriter, err api.Error) {
	cookie := &http.Cookie{
		Name:     "scToken",
		Value:    "",
		Path:     "/",
		Domain:   a.cfg.CookieDomain,
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	http.Error(w, err.Error(), err.Status())
}
