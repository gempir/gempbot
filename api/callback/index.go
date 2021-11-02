package callback

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
)

var (
	cfg         *config.Config
	db          *store.Database
	helixClient *helixclient.Client
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg = config.FromEnv()
	db = store.NewDatabase(cfg)
	helixClient = helixclient.NewClient(cfg, db)

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

	token := auth.CreateApiToken(cfg.Secret, validateResp.Data.UserID)

	err = db.SaveUserAccessToken(r.Context(), validateResp.Data.UserID, resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "))
	if err != nil {
		fmt.Fprintf(w, "failed to set userAccessToken in callback: %s", err)
		return
	}

	dashboardRedirect(w, r, token)
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
