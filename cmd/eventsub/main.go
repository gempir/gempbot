package main

import (
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	nickHelix "github.com/nicklaw5/helix/v2"
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
		if sub.Type == nickHelix.EventSubTypeChannelPointsCustomRewardRedemptionAdd {
			if sub.ForeignID == "" {
				log.Warnf("Found subscription reward without foreignID channel: %s id: %s", sub.TargetTwitchID, sub.SubscriptionID)
				err := subscriptionManager.RemoveSubscription(sub.SubscriptionID)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
