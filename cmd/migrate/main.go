package main

import (
	"github.com/gempir/bot/pkg/config"
	"github.com/gempir/bot/pkg/store"
)

func main() {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)

	db.Migrate()
}
