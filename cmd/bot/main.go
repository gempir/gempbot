package main

import (
	"github.com/gempir/gempbot/cmd/bot/collector"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/store"
)

func main() {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)

	helixClient := helix.NewClient(cfg, db)
	go helixClient.StartRefreshTokenRoutine()

	bot := collector.NewBot(cfg, db, helixClient)
	go bot.Connect()

	<-bot.Done
}
