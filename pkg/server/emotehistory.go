package server

import (
	"net/http"
	"strconv"

	"github.com/gempir/gempbot/pkg/api"
)

func (a *Api) EmoteHistoryHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	userID := ""

	if username == "" {
		authResult, _, err := a.authClient.AttemptAuth(r, w)
		if err != nil {
			return
		}
		userID = authResult.Data.UserID

		if r.URL.Query().Get("managing") != "" {
			userID, err = a.userAdmin.CheckEditor(r, a.userAdmin.GetUserConfig(userID))
			if err != nil {
				http.Error(w, err.Error(), err.Status())
				return
			}
		}
	} else {
		user, err := a.helixClient.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID = user.ID
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	api.WriteJson(w, a.db.GetEmoteHistory(r.Context(), userID, pageNumber, api.EMOTEHISTORY_PAGE_SIZE, r.URL.Query().Has("added")), http.StatusOK)
}
