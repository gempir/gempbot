package main

import (
	"flag"

	"github.com/gempir/bitraft/pkg/config"
	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/gempir/bitraft/services/ingester/collector"
)

func main() {
	configFile := flag.String("config", "../../config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)
	redis := store.NewRedis()
	db := store.NewDatabase(cfg)

	helixClient := helix.NewClient(cfg.ClientID, cfg.ClientSecret, cfg.ApiBaseUrl+"/api/callback", cfg.Secret)
	go helixClient.StartRefreshTokenRoutine()

	bot := collector.NewBot(cfg, redis, db, helixClient)
	bot.Connect()
}
