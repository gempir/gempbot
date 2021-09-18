package main

import (
	"flag"

	"github.com/gempir/gempbot/cmd/bot/commander/predictions"
	"github.com/gempir/gempbot/cmd/bot/commander/reader"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/store"
)

func main() {
	configFile := flag.String("config", "../../config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)
	rStore := store.NewRedis()
	db := store.NewDatabase(cfg)

	helixClient := helix.NewClient(cfg)
	go helixClient.StartRefreshTokenRoutine()

	redis := store.NewRedis()

	predictionHandler := predictions.NewHandler(helixClient, redis, db)

	predictionsListener := reader.NewListener(db, rStore, predictionHandler)
	predictionsListener.RegisterDefaultCommands()
	predictionsListener.StartListener()
}
