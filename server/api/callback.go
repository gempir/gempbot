package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/log"
)

func (a *Api) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	resp, err := a.helixClient.Client.RequestUserAccessToken(code)
	if err != nil || resp.StatusCode >= 400 {
		log.Errorf("failed to request userAccessToken: %s %s", err, resp.ErrorMessage)
		a.dashboardRedirect(w, r, "")
		return
	}

	// validate
	success, validateResp, err := a.helixClient.Client.ValidateToken(resp.Data.AccessToken)
	if !success || err != nil {
		fmt.Fprintf(w, "failed to veryify new Token %s", err)
		return
	}

	token := auth.CreateApiToken(a.cfg.Secret, validateResp.Data.UserID)

	err = a.db.SaveUserAccessToken(r.Context(), validateResp.Data.UserID, resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "))
	if err != nil {
		fmt.Fprintf(w, "failed to set userAccessToken in callback: %s", err)
		return
	}

	a.dashboardRedirect(w, r, token)
}

func (a *Api) dashboardRedirect(w http.ResponseWriter, r *http.Request, scToken string) {
	cookie := http.Cookie{
		Name:    "scToken",
		Value:   scToken,
		Domain:  a.cfg.CookieDomain,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, a.cfg.WebBaseUrl, http.StatusFound)
}
