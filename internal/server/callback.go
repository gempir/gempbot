package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/gempbot/internal/auth"
	"github.com/gempir/gempbot/internal/log"
)

func (a *Api) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	resp, err := a.helixClient.RequestUserAccessToken(code)
	if err != nil || resp.StatusCode >= 400 {
		if err != nil {
			log.Errorf("failed to request user access token: %s %s", err.Error(), resp.ErrorMessage)
		} else {
			log.Errorf("failed to request userAccessToken: %s", resp.ErrorMessage)
		}
		a.dashboardRedirect(w, r, "")
		return
	}

	// validate
	success, validateResp, err := a.helixClient.ValidateToken(resp.Data.AccessToken)
	if !success || err != nil {
		fmt.Fprintf(w, "failed to veryify new Token %s", err)
		return
	}

	token := auth.CreateApiToken(a.cfg.Secret, validateResp)

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
