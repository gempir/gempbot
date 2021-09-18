package main

import (
	"fmt"

	"github.com/gempir/gempbot/cmd/bot/collector"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/store"
)

func main() {

	cfg := config.FromEnv()
	fmt.Printf("ENV: '%s'\n", cfg.DbName)
	db := store.NewDatabase(cfg)

	helixClient := helix.NewClient(cfg)
	go helixClient.StartRefreshTokenRoutine()

	bot := collector.NewBot(cfg, db, helixClient)
	bot.Connect()
}
