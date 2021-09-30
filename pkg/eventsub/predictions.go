package eventsub

import (
	"github.com/gempir/gempbot/pkg/log"
	nickHelix "github.com/nicklaw5/helix"
)

func (esm *EventSubManager) SubscribePredictions(userID string) {
	response, err := esm.helixClient.CreateEventSubSubscription(userID, esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+nickHelix.EventSubTypeChannelPredictionBegin, "channel.prediction.begin")
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	log.Infof("[%d] created subscription %s", response.StatusCode, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new sub in %s %s", userID, sub.Type)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type)
	}

	response, err = esm.helixClient.CreateEventSubSubscription(userID, esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+nickHelix.EventSubTypeChannelPredictionLock, "channel.prediction.lock")
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	log.Infof("[%d] created subscription %s", response.StatusCode, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new sub in %s %s", userID, sub.Type)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type)
	}

	response, err = esm.helixClient.CreateEventSubSubscription(userID, esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+nickHelix.EventSubTypeChannelPredictionEnd, "channel.prediction.end")
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	log.Infof("[%d] created subscription %s", response.StatusCode, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new sub in %s %s", userID, sub.Type)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type)
	}
}
