package blocks

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

const PAGE_SIZE = 20

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)
	helixClient := helix.NewClient(cfg, db)
	userAdmin := user.NewUserAdmin(cfg, db, helixClient, nil)
	auth := auth.NewAuth(cfg, db, helixClient)

	authResp, _, apiErr := auth.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID := authResp.Data.UserID

	if r.URL.Query().Get("managing") != "" {
		userID, apiErr = userAdmin.CheckEditor(r, userAdmin.GetUserConfig(userID))
		if apiErr != nil {
			http.Error(w, apiErr.Error(), apiErr.Status())
			return
		}
	}

	if r.Method == http.MethodGet {
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}

		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		blocks := db.GetEmoteBlocks(userID, pageNumber, PAGE_SIZE)
		api.WriteJson(w, blocks, http.StatusOK)
		return
	}
}
