package emotehistory

import (
	"net/http"
	"strconv"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)
	helixClient := helixclient.NewClient(cfg, db)
	auth := auth.NewAuth(cfg, db, helixClient)
	userAdmin := user.NewUserAdmin(cfg, db, helixClient, nil)

	username := r.URL.Query().Get("username")
	userID := ""

	if username == "" {
		authResult, _, err := auth.AttemptAuth(r, w)
		if err != nil {
			return
		}
		userID = authResult.Data.UserID

		if r.URL.Query().Get("managing") != "" {
			userID, err = userAdmin.CheckEditor(r, userAdmin.GetUserConfig(userID))
			if err != nil {
				http.Error(w, err.Error(), err.Status())
				return
			}
		}
	} else {
		user, err := helixClient.GetUserByUsername(username)
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

	api.WriteJson(w, db.GetEmoteHistory(r.Context(), userID, pageNumber, api.EMOTEHISTORY_PAGE_SIZE), http.StatusOK)
}
