package predictionhistory

import (
	"net/http"
	"strconv"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)
	helixClient := helix.NewClient(cfg, db)
	auth := auth.NewAuth(cfg, db, helixClient)
	userAdmin := user.NewUserAdmin(db, helixClient)

	username := r.URL.Query().Get("username")
	userID := ""

	if username == "" {
		auth, _, err := auth.AttemptAuth(r, w)
		if err != nil {
			return
		}
		userID = auth.Data.UserID

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

	api.WriteJson(w, db.GetPredictions(r.Context(), userID, pageNumber, api.PREDICTIONS_PAGE_SIZE), http.StatusOK)
}
