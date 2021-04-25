package main

import (
	"flag"

	"github.com/gempir/spamchamp/pkg/config"
	"github.com/gempir/spamchamp/pkg/helix"
	"github.com/gempir/spamchamp/pkg/store"
	"github.com/gempir/spamchamp/services/api/server"
	"github.com/gempir/spamchamp/services/api/stats"
)

var broadcastQueue = make(chan server.BroadcastMessage)

func main() {
	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)
	rStore := store.NewStore()
	helixClient := helix.NewClient(cfg.ClientID, cfg.ClientSecret)
	go helixClient.StartRefreshTokenRoutine()

	server := server.NewServer(cfg, &helixClient, rStore, broadcastQueue)
	broadcaster := stats.NewBroadcaster(broadcastQueue, rStore, &helixClient)
	go broadcaster.Start()

	server.Start()
}
