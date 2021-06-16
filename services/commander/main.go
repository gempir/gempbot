package main

import (
	"flag"

	"github.com/gempir/bitraft/pkg/config"
	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/gempir/bitraft/services/commander/predictions"
	"github.com/gempir/bitraft/services/commander/reader"
)

func main() {
	configFile := flag.String("config", "../../config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)
	rStore := store.NewRedis()
	db := store.NewDatabase(cfg)

	helixClient := helix.NewClient(cfg.ClientID, cfg.ClientSecret, cfg.ApiBaseUrl+"/api/callback", cfg.Secret)
	go helixClient.StartRefreshTokenRoutine()

	redis := store.NewRedis()

	predictionHandler := predictions.NewHandler(helixClient, redis, db)

	predictionsListener := reader.NewListener(db, rStore, predictionHandler)
	predictionsListener.RegisterDefaultCommands()
	predictionsListener.StartListener()
}
