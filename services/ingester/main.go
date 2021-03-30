package main

import (
	"flag"

	"github.com/gempir/spamchamp/pkg/config"
	"github.com/gempir/spamchamp/pkg/helix"
	"github.com/gempir/spamchamp/pkg/store"
	"github.com/gempir/spamchamp/services/ingester/collector"
)

func main() {
	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)
	rStore := store.NewStore()
	helixClient := helix.NewClient(cfg.ClientID, cfg.ClientSecret)

	bot := collector.NewBot(cfg, &helixClient, rStore)
	go bot.Connect()

	helixClient.StartRefreshTokenRoutine()
}
