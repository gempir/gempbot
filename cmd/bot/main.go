package main

import (
	"github.com/gempir/gempbot/cmd/bot/collector"
	"github.com/gempir/gempbot/cmd/bot/commander"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/store"
)

func main() {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)

	helixClient := helix.NewClient(cfg)
	go helixClient.StartRefreshTokenRoutine()

	handler := commander.NewHandler(helixClient, db)

	listener := commander.NewListener(db, handler)
	listener.RegisterDefaultCommands()

	bot := collector.NewBot(cfg, db, helixClient, listener)
	bot.Connect()

	<-bot.Done
}
