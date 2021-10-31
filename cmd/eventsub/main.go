package main

import (
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
)

var (
	cfg         *config.Config
	db          *store.Database
	helixClient *helix.Client
)

func main() {
	cfg = config.FromEnv()
	db = store.NewDatabase(cfg)
	helixClient = helix.NewClient(cfg, db)
	subscriptionManager := eventsub.NewSubscriptionManager(cfg, db, helixClient)

	for _, sub := range db.GetAllSubscriptions() {
		err := subscriptionManager.RemoveSubscription(sub.SubscriptionID)
		if err != nil {
			log.Error(err)
		}
	}
}
