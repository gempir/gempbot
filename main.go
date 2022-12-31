package main

import (
	"net/http"
	"os"

	"github.com/gempir/gempbot/internal/auth"
	"github.com/gempir/gempbot/internal/bot"
	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/election"
	"github.com/gempir/gempbot/internal/emotechief"
	"github.com/gempir/gempbot/internal/emoteservice"
	"github.com/gempir/gempbot/internal/eventsubmanager"
	"github.com/gempir/gempbot/internal/eventsubsubscription"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/media"
	"github.com/gempir/gempbot/internal/server"
	"github.com/gempir/gempbot/internal/store"
	"github.com/gempir/gempbot/internal/user"
	"github.com/gempir/gempbot/internal/ws"
	"github.com/rs/cors"
)

func main() {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 1 && argsWithoutProg[0] == "migrate" {
		db.Migrate()
		os.Exit(0)
		return
	}

	helixClient := helixclient.NewClient(cfg, db)
	go helixClient.StartRefreshTokenRoutine()

	userAdmin := user.NewUserAdmin(cfg, db, helixClient, nil)
	authClient := auth.NewAuth(cfg, db, helixClient)

	bot := bot.NewBot(cfg, db, helixClient)
	go bot.Connect()

	seventvClient := emoteservice.NewSevenTvClient(db)

	emoteChief := emotechief.NewEmoteChief(cfg, db, helixClient, bot.ChatClient, seventvClient)
	eventsubSubscriptionManager := eventsubsubscription.NewSubscriptionManager(cfg, db, helixClient)
	channelPointManager := channelpoint.NewChannelPointManager(cfg, helixClient, db)
	electionManager := election.NewElectionManager(db, helixClient, channelPointManager, eventsubSubscriptionManager, bot, seventvClient)
	go electionManager.StartElectionManagerRoutine()

	go func() {
		go eventsubSubscriptionManager.RemoveSubscription("e587f4cd-ee4a-4326-8c02-0d85abcc194d")
		go eventsubSubscriptionManager.RemoveSubscription("753a10fc-3f84-450c-b84c-4afd33996091")
		go eventsubSubscriptionManager.RemoveSubscription("4ce501c7-7e98-430a-858b-95ce83023ac2")
		go eventsubSubscriptionManager.RemoveSubscription("dac637ba-9dc4-4144-9918-b97e14d9df26")
		go eventsubSubscriptionManager.RemoveSubscription("7a4ee1af-220d-4ece-abdf-87d59938c364")
		go eventsubSubscriptionManager.RemoveSubscription("09fc7c36-186f-48bb-981b-ed99fbc6f9ce")
		go eventsubSubscriptionManager.RemoveSubscription("fc58ee4e-3875-496c-ab25-2baae42d13a4")
		go eventsubSubscriptionManager.RemoveSubscription("42c1eba4-fde1-4ab2-8075-67b9abd666c8")
		go eventsubSubscriptionManager.RemoveSubscription("237d892f-286e-404d-a2cc-630781c5fe4b")
		go eventsubSubscriptionManager.RemoveSubscription("0362eb95-1182-45e7-b93e-da1151ec83fc")
		go eventsubSubscriptionManager.RemoveSubscription("b5e22f9d-2468-4454-a74e-c649f198fa2a")
		go eventsubSubscriptionManager.RemoveSubscription("29bf821d-4d83-4b30-b8d5-7f18ba1d1384")
	}()

	mediaManager := media.NewMediaManager(db, helixClient, bot)
	wsHandler := ws.NewWsHandler(authClient, mediaManager)
	eventsubManager := eventsubmanager.NewEventsubManager(cfg, helixClient, db, emoteChief, bot.ChatClient, electionManager)

	apiHandlers := server.NewApi(cfg, db, helixClient, userAdmin, authClient, bot, emoteChief, eventsubManager, eventsubSubscriptionManager, channelPointManager, seventvClient, wsHandler)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<img src=\"data:image/webp;base64,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\" />"))
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 page not found", http.StatusNotFound)
	})
	mux.HandleFunc("/api/blocks", apiHandlers.BlocksHandler)
	mux.HandleFunc("/api/botconfig", apiHandlers.BotConfigHandler)
	mux.HandleFunc("/api/callback", apiHandlers.CallbackHandler)
	mux.HandleFunc("/api/emotehistory", apiHandlers.EmoteHistoryHandler)
	mux.HandleFunc("/api/eventsub", apiHandlers.EventSubHandler)
	mux.HandleFunc("/api/reward", apiHandlers.RewardHandler)
	mux.HandleFunc("/api/election", apiHandlers.ElectionHandler)
	mux.HandleFunc("/api/nominations", apiHandlers.NominationsHandler)
	mux.HandleFunc("/api/nominations/vote", apiHandlers.NominationVoteHandler)
	mux.HandleFunc("/api/subscriptions", apiHandlers.SubscriptionsHandler)
	mux.HandleFunc("/api/userconfig", apiHandlers.UserConfigHandler)
	mux.HandleFunc("/api/ws", wsHandler.HandleWs)

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
