package main

import (
	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/store"
)

func main() {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)

	db.Migrate()
}
