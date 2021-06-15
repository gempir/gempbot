package main

import (
	"flag"

	"github.com/gempir/bitraft/pkg/config"
	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/gempir/bitraft/services/commander/predictions"
)

func main() {
	configFile := flag.String("config", "../../config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)
	rStore := store.NewRedis()
	db := store.NewDatabase(cfg)

	helixClient := helix.NewClient(cfg.ClientID, cfg.ClientSecret, cfg.ApiBaseUrl+"/api/callback", cfg.Secret)
	go helixClient.StartRefreshTokenRoutine()

	predictionsListener := predictions.NewListener(db, rStore)
	predictionsListener.RegisterDefaultCommands()
	predictionsListener.StartListener()
}
