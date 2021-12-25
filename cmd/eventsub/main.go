package main

import (
	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/emotechief"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/nicklaw5/helix/v2"
)

var (
	cfg         *config.Config
	db          *store.Database
	helixClient *helixclient.Client
)

func main() {
	cfg = config.FromEnv()
	db = store.NewDatabase(cfg)
	helixClient = helixclient.NewClient(cfg, db)
	subscriptionManager := eventsub.NewSubscriptionManager(cfg, db, helixClient)
	chatClient := chat.NewClient(cfg)
	emotechief := emotechief.NewEmoteChief(cfg, db, helixClient, chatClient)
	eventsubManager := eventsub.NewEventSubManager(cfg, helixClient, db, emotechief, chatClient)

	for _, sub := range db.GetAllSubscriptions() {
		if sub.Type == helix.EventSubTypeChannelPredictionBegin || sub.Type == helix.EventSubTypeChannelPredictionEnd || sub.Type == helix.EventSubTypeChannelPredictionLock {
			err := subscriptionManager.RemoveSubscription(sub.SubscriptionID)
			if err != nil {
				log.Error(err)
			}

			eventsubManager.SubscribePredictions(sub.TargetTwitchID)
		}
	}
}
