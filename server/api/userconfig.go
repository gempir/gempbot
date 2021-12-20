package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/user"
)

func (a *Api) UserConfigHandler(w http.ResponseWriter, r *http.Request) {
	chatClient := chat.NewClient(a.cfg)
	go chatClient.Connect(func() {})

	authResp, _, err := a.authClient.AttemptAuth(r, w)
	if err != nil {
		return
	}

	if r.Method == http.MethodGet {
		userConfig := a.userAdmin.GetUserConfig(authResp.Data.UserID)

		if r.URL.Query().Get("managing") != "" {
			ownerUserID, err := a.userAdmin.CheckEditor(r, userConfig)
			if err != nil {
				http.Error(w, err.Error(), err.Status())
				return
			}

			editorFor := userConfig.Protected.EditorFor
			userConfig = a.userAdmin.GetUserConfig(ownerUserID)

			userConfig.Protected.EditorFor = editorFor
		}

		userConfig, err := a.userAdmin.ConvertUserConfig(userConfig, true)
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

		err = a.userAdmin.ProcessConfig(r.Context(), authResp.Data.UserID, authResp.Data.Login, newConfig, r.URL.Query().Get("managing"))
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
