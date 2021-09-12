package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/bot/pkg/config"
	"github.com/gempir/bot/pkg/helix"
	"github.com/gempir/bot/pkg/log"
	"github.com/gempir/bot/pkg/store"
	"github.com/golang-jwt/jwt"
)

var (
	cfg         *config.Config
	db          *store.Database
	helixClient *helix.Client
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg = config.FromEnv()
	db = store.NewDatabase(cfg)
	helixClient = helix.NewClient(cfg)

	code := r.URL.Query().Get("code")

	resp, err := helixClient.Client.RequestUserAccessToken(code)
	if err != nil || resp.StatusCode >= 400 {
		log.Errorf("failed to request userAccessToken: %s %s", err, resp.ErrorMessage)
		dashboardRedirect(w, r, "")
		return
	}

	// validate
	success, validateResp, err := helixClient.Client.ValidateToken(resp.Data.AccessToken)
	if !success || err != nil {
		fmt.Fprintf(w, "failed to veryify new Token %s", err)
		return
	}

	token, err := createApiToken(validateResp.Data.UserID)
	if err != nil {
		fmt.Fprintf(w, "failed to create jwt token in callback %s", err)
		return
	}

	err = db.SaveUserAccessToken(validateResp.Data.UserID, resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "))
	if err != nil {
		fmt.Fprintf(w, "failed to set userAccessToken in callback: %s", err)
		return
	}

	err = db.SaveBotConfig(store.BotConfig{OwnerTwitchID: validateResp.Data.UserID, JoinBot: true})
	if err != nil {
		log.Error(err)
	}

	dashboardRedirect(w, r, token)
}

type tokenClaims struct {
	UserID         string
	StandardClaims jwt.StandardClaims
}

func (t *tokenClaims) Valid() error {
	return nil
}

func createApiToken(userID string) (string, error) {
	expirationTime := time.Now().Add(365 * 24 * time.Hour)
	claims := &tokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Secret))

	return tokenString, err
}

func dashboardRedirect(w http.ResponseWriter, r *http.Request, scToken string) {
	cookie := http.Cookie{
		Name:    "scToken",
		Value:   scToken,
		Domain:  cfg.CookieDomain,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, cfg.WebBaseUrl, http.StatusFound)
}
