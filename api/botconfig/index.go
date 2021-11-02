package userconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	chatClient := chat.NewClient(cfg)
	go chatClient.Connect()
	db := store.NewDatabase(cfg)
	helixClient := helixclient.NewClient(cfg, db)
	auth := auth.NewAuth(cfg, db, helixClient)
	userAdmin := user.NewUserAdmin(cfg, db, helixClient, chatClient)

	authResp, _, apiErr := auth.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID := authResp.Data.UserID
	ownerLogin := authResp.Data.Login

	if r.URL.Query().Get("managing") != "" {
		userID, apiErr = userAdmin.CheckEditor(r, userAdmin.GetUserConfig(userID))
		if apiErr != nil {
			http.Error(w, apiErr.Error(), apiErr.Status())
			return
		}

		uData, err := helixClient.GetUserByUserID(userID)
		if err != nil {
			api.WriteJson(w, fmt.Errorf("could not find managing user in helix"), http.StatusBadRequest)
			return
		}
		ownerLogin = uData.Login
	}

	if r.Method == http.MethodGet {
		cfg, err := db.GetBotConfig(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
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

		dbErr := db.SaveBotConfig(context.Background(), botCfg)
		if dbErr != nil {
			log.Error(dbErr)
			api.WriteJson(w, fmt.Errorf("failed to save bot config"), http.StatusInternalServerError)
			return
		}
		log.Info("waiting for bot connection")
		chatClient.WaitForConnect()
		if botCfg.JoinBot {
			chatClient.JoinBot(ownerLogin)
		} else {
			chatClient.PartBot(ownerLogin)
		}

		return
	}

	http.Error(w, "unknown method", http.StatusMethodNotAllowed)
}
