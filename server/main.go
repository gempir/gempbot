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
	"github.com/rs/cors"
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

	mux := http.NewServeMux()

	mux.HandleFunc("/api/blocks", apiHandlers.BlocksHandler)
	mux.HandleFunc("/api/botconfig", apiHandlers.BotConfigHandler)
	mux.HandleFunc("/api/callback", apiHandlers.CallbackHandler)
	mux.HandleFunc("/api/emotehistory", apiHandlers.EmoteHistoryHandler)
	mux.HandleFunc("/api/eventsub", apiHandlers.EventSubHandler)
	mux.HandleFunc("/api/reward", apiHandlers.RewardHandler)
	mux.HandleFunc("/api/subscriptions", apiHandlers.SubscriptionsHandler)
	mux.HandleFunc("/api/userconfig", apiHandlers.UserConfigHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.WebBaseUrl},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)
	err := http.ListenAndServe("localhost:3010", handler)
	if err != nil {
		log.Fatal(err)
	}
}
