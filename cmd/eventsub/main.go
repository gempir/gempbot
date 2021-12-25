package main

import (
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/nicklaw5/helix/v2"
)

func main() {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)
	hClient := helixclient.NewClient(cfg, db)
	// subscriptionManager := eventsub.NewSubscriptionManager(cfg, db, hClient)

	subscriptionIds := map[string]bool{}

	for _, sub := range hClient.GetAllSubscriptions(helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd) {
		// log.Info(sub.Transport.Callback)
		// if strings.HasPrefix(sub.Transport.Callback, "https://bot.gempir.com/api/eventsub") {
		// 	err := subscriptionManager.RemoveSubscription(sub.ID)
		// 	if err != nil {
		// 		log.Error(err)
		// 	}
		// 	log.Info(sub)
		// 	subscriptionManager.SubscribeRewardRedemptionAdd(sub.Condition.BroadcasterUserID, sub.Condition.RewardID)
		// }

		// subscriptionManager.SubscribeRewardRedemptionAdd(sub.Condition.BroadcasterUserID, sub.Condition.RewardID)
		subscriptionIds[sub.ID] = true
		// db.AddEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, sub.Version, sub.Type, sub.Condition.RewardID)
	}

	// for _, sub := range hClient.GetAllSubscriptions(helix.EventSubTypeChannelPointsCustomRewardRedemptionUpdate) {
	// 	log.Info(sub)
	// 	// subscriptionManager.SubscribeRewardRedemptionUpdate(sub.Condition.BroadcasterUserID, sub.Condition.RewardID)
	// 	subscriptionIds[sub.ID] = true
	// 	db.AddEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, sub.Version, sub.Type, sub.Condition.RewardID)
	// }

	// for _, sub := range hClient.GetAllSubscriptions(helix.EventSubTypeChannelPredictionBegin) {
	// 	subscriptionIds[sub.ID] = true
	// 	db.AddEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, sub.Version, sub.Type, "")
	// }

	// for _, sub := range hClient.GetAllSubscriptions(helix.EventSubTypeChannelPredictionEnd) {
	// 	subscriptionIds[sub.ID] = true
	// 	db.AddEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, sub.Version, sub.Type, "")
	// }

	// for _, sub := range hClient.GetAllSubscriptions(helix.EventSubTypeChannelPredictionLock) {
	// 	subscriptionIds[sub.ID] = true
	// 	db.AddEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, sub.Version, sub.Type, "")
	// }

	// for _, sub := range db.GetAllSubscriptions() {
	// 	if !subscriptionIds[sub.SubscriptionID] {
	// 		log.Infof("Removing subscription, not found in helix %v", sub)
	// 		db.RemoveEventSubSubscription(sub.SubscriptionID)
	// 	}
	// }
}
