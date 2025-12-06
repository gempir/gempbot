package main

import (
	"net/http"
	"os"

	"github.com/gempir/gempbot/internal/auth"
	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/emotechief"
	"github.com/gempir/gempbot/internal/emoteservice"
	"github.com/gempir/gempbot/internal/eventsubmanager"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/server"
	"github.com/gempir/gempbot/internal/store"
	"github.com/gempir/gempbot/internal/user"
	"github.com/gempir/gempbot/internal/ws"
	"github.com/rs/cors"
)

func main() {
	cfg := config.FromEnv()
	log.SetLogLevel(cfg.LogLevel)
	db := store.NewDatabase(cfg)

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 1 && argsWithoutProg[0] == "migrate" {
		db.Migrate()
		os.Exit(0)
		return
	}

	helixClient := helixclient.NewClient(cfg, db)

	userAdmin := user.NewUserAdmin(cfg, db, helixClient)
	authClient := auth.NewAuth(cfg, db, helixClient)

	seventvClient := emoteservice.NewSevenTvClient(db)

	emoteChief := emotechief.NewEmoteChief(cfg, db, helixClient, seventvClient)
	channelPointManager := channelpoint.NewChannelPointManager(cfg, helixClient, db)
	wsHandler := ws.NewWsHandler(authClient)
	eventsubManager := eventsubmanager.NewEventsubManager(cfg, helixClient, db, emoteChief)

	apiHandlers := server.NewApi(cfg, db, helixClient, userAdmin, authClient, emoteChief, eventsubManager, channelPointManager, seventvClient, wsHandler)

	mux := http.NewServeMux()

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 page not found", http.StatusNotFound)
	})
	mux.HandleFunc("/api/blocks", apiHandlers.BlocksHandler)
	mux.HandleFunc("/api/bot/config", apiHandlers.BotConfigHandler)
	mux.HandleFunc("/api/fossabot", apiHandlers.FossabotHandler)
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

	log.Info("Starting server on " + cfg.ListenAddress)
	err := http.ListenAndServe(cfg.ListenAddress, handler)
	if err != nil {
		log.Fatal(err)
	}
}
