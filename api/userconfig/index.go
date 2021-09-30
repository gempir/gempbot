package userconfig

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	chatClient := chat.NewClient(cfg)
	go chatClient.Connect()
	db := store.NewDatabase(cfg)
	helixClient := helix.NewClient(cfg, db)
	auth := auth.NewAuth(cfg, db, helixClient)
	userAdmin := user.NewUserAdmin(cfg, db, helixClient, chatClient)

	authResp, _, err := auth.AttemptAuth(r, w)
	if err != nil {
		return
	}

	if r.Method == http.MethodGet {
		userConfig := userAdmin.GetUserConfig(authResp.Data.UserID)

		if r.URL.Query().Get("managing") != "" {
			ownerUserID, err := userAdmin.CheckEditor(r, userConfig)
			if err != nil {
				http.Error(w, err.Error(), err.Status())
				return
			}

			editorFor := userConfig.Protected.EditorFor
			userConfig = userAdmin.GetUserConfig(ownerUserID)

			userConfig.Protected.EditorFor = editorFor
		}

		userConfig, err := userAdmin.ConvertUserConfig(userConfig, true)
		if err != nil {
			http.Error(w, err.Error(), err.Status())
			return
		}

		api.WriteJson(w, userConfig, http.StatusOK)
		return

	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Failed reading update body: %s", err)
			http.Error(w, "Failure saving body "+err.Error(), http.StatusInternalServerError)
			return
		}

		var newConfig user.UserConfig
		if err := json.Unmarshal(body, &newConfig); err != nil {
			log.Errorf("Failed unmarshalling userConfig: %s", err)
			http.Error(w, "Failure unmarshalling config "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = userAdmin.ProcessConfig(r.Context(), authResp.Data.UserID, authResp.Data.Login, newConfig, r.URL.Query().Get("managing"))
		if err != nil {
			log.Errorf("failed processing config: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		api.WriteJson(w, nil, http.StatusOK)
		return
	}

	http.Error(w, "unknown method", http.StatusMethodNotAllowed)
}
