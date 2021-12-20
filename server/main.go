package main

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
	"github.com/gempir/gempbot/server/api"
	"github.com/gempir/gempbot/server/bot"
)

func main() {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)

	helixClient := helixclient.NewClient(cfg, db)
	go helixClient.StartRefreshTokenRoutine()

	userAdmin := user.NewUserAdmin(cfg, db, helixClient, nil)
	authClient := auth.NewAuth(cfg, db, helixClient)

	bot := bot.NewBot(cfg, db, helixClient)
	go bot.Connect()

	apiHandlers := api.NewApi(cfg, db, helixClient, userAdmin, authClient)

	http.HandleFunc("/blocks", apiHandlers.BlocksHandler)
	http.HandleFunc("/botconfig", apiHandlers.BotConfigHandler)
	http.HandleFunc("/callback", apiHandlers.CallbackHandler)
	http.HandleFunc("/emotehistory", apiHandlers.EmoteHistoryHandler)
	http.HandleFunc("/eventsub", apiHandlers.EventSubHandler)
	http.HandleFunc("/reward", apiHandlers.RewardHandler)
	http.HandleFunc("/subscriptions", apiHandlers.SubscriptionsHandler)
	http.HandleFunc("/userconfig", apiHandlers.UserConfigHandler)

	err := http.ListenAndServe("127.0.0.1:3010", nil)
	if err != nil {
		log.Fatal(err)
	}
}
