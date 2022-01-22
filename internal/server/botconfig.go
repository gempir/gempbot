package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

func (a *Api) BotConfigHandler(w http.ResponseWriter, r *http.Request) {
	authResp, _, apiErr := a.authClient.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID := authResp.Data.UserID
	ownerLogin := authResp.Data.Login

	if r.URL.Query().Get("managing") != "" {
		userID, apiErr = a.userAdmin.CheckEditor(r, a.userAdmin.GetUserConfig(userID))
		if apiErr != nil {
			http.Error(w, apiErr.Error(), apiErr.Status())
			return
		}

		uData, err := a.helixClient.GetUserByUserID(userID)
		if err != nil {
			api.WriteJson(w, fmt.Errorf("could not find managing user in helix"), http.StatusBadRequest)
			return
		}
		ownerLogin = uData.Login
	}

	if r.Method == http.MethodGet {
		cfg, err := a.db.GetBotConfig(userID)
		if err != nil {
			log.Error(err)
		}

		api.WriteJson(w, cfg, http.StatusOK)
		return
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Failed reading update body: %s", err)
			api.WriteJson(w, fmt.Errorf("Failure unmarshalling config "+err.Error()), http.StatusInternalServerError)
			return
		}

		var botCfg store.BotConfig
		if err := json.Unmarshal(body, &botCfg); err != nil {
			log.Errorf("Failed unmarshalling botConfig: %s", err)
			api.WriteJson(w, fmt.Errorf("Failure unmarshalling config "+err.Error()), http.StatusInternalServerError)
			return
		}
		botCfg.OwnerTwitchID = userID

		dbErr := a.db.SaveBotConfig(context.Background(), botCfg)
		if dbErr != nil {
			log.Error(dbErr)
			api.WriteJson(w, fmt.Errorf("failed to save bot config"), http.StatusInternalServerError)
			return
		}
		if botCfg.JoinBot {
			a.bot.Join(ownerLogin)
		} else {
			a.bot.Part(ownerLogin)
		}

		return
	}

	http.Error(w, "unknown method", http.StatusMethodNotAllowed)
}
