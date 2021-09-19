package main

import (
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/store"
)

func main() {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)

	db.Migrate()
}
