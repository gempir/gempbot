package main

import (
	"flag"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/pkg/config"
	"github.com/gempir/spamchamp/services/bot/api"
	"github.com/gempir/spamchamp/services/bot/helix"
	"github.com/gempir/spamchamp/services/bot/stats"
	"github.com/gempir/spamchamp/services/bot/store"
)

var messageQueue = make(chan twitch.PrivateMessage)
var broadcastQueue = make(chan api.BroadcastMessage)

func main() {
	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)
	rStore := store.NewStore()
	helixClient := helix.NewClient(cfg.ClientID, cfg.ClientSecret)
	go helixClient.StartRefreshTokenRoutine()

	server := api.NewServer(cfg, &helixClient, broadcastQueue)
	broadcaster := stats.NewBroadcaster(messageQueue, broadcastQueue, rStore)
	go broadcaster.Start()

	server.Start()
}
