package eventsubsubscription

import (
	"net/http"

	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/nicklaw5/helix/v2"
)

type SubscriptionManager struct {
	cfg         *config.Config
	db          *store.Database
	helixClient helixclient.Client
}

func NewSubscriptionManager(cfg *config.Config, db *store.Database, helixClient helixclient.Client) *SubscriptionManager {
	return &SubscriptionManager{
		helixClient: helixClient,
		db:          db,
		cfg:         cfg,
	}
}

func (esm *SubscriptionManager) SubscribeRewardRedemptionAdd(userID, rewardId string) {
	response, err := esm.helixClient.CreateRewardEventSubSubscription(
		userID,
		esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd,
		helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd,
		rewardId,
	)
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	if response.StatusCode == http.StatusForbidden {
		log.Errorf("Forbidden subscription %s", response.ErrorMessage)
		return
	}

	log.Infof("[%d] SubscribeRewardRedemptionAdd %s %s", response.StatusCode, response.Error, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new subscription for %s id: %s", userID, sub.ID)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type, rewardId)
	}
}

func (esm *SubscriptionManager) SubscribeRewardRedemptionUpdate(userID, rewardId string) {
	response, err := esm.helixClient.CreateRewardEventSubSubscription(
		userID,
		esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+helix.EventSubTypeChannelPointsCustomRewardRedemptionUpdate,
		helix.EventSubTypeChannelPointsCustomRewardRedemptionUpdate,
		rewardId,
	)
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	if response.StatusCode == http.StatusForbidden {
		log.Errorf("Forbidden subscription %s", response.ErrorMessage)
		return
	}

	log.Infof("[%d] SubscribeRewardRedemptionUpdate %s %s", response.StatusCode, response.Error, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new subscription for %s id: %s", userID, sub.ID)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type, rewardId)
	}
}

func (esm *SubscriptionManager) RemoveSubscription(subscriptionID string) error {
	response, err := esm.helixClient.RemoveEventSubSubscription(subscriptionID)
	if err != nil {
		return err
	}

	log.Infof("[%d] removed EventSubSubscription", response.StatusCode)
	esm.db.RemoveEventSubSubscription(subscriptionID)

	return nil
}
