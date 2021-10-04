package eventsub

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	nickHelix "github.com/nicklaw5/helix/v2"
)

type SubscriptionManager struct {
	cfg         *config.Config
	db          *store.Database
	helixClient *helix.Client
}

func NewSubscriptionManager(cfg *config.Config, db *store.Database, helixClient *helix.Client) *SubscriptionManager {
	return &SubscriptionManager{
		helixClient: helixClient,
		db:          db,
		cfg:         cfg,
	}
}

func (esm *SubscriptionManager) SubscribeRewardRedemptionAdd(userID, rewardId string) {
	response, err := esm.helixClient.CreateRewardEventSubSubscription(
		userID,
		esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+nickHelix.EventSubTypeChannelPointsCustomRewardRedemptionAdd,
		nickHelix.EventSubTypeChannelPointsCustomRewardRedemptionAdd,
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

	log.Infof("[%d] subscription %s %s", response.StatusCode, response.Error, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new subscription for %s id: %s", userID, sub.ID)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type, rewardId)
	}
}

func (esm *SubscriptionManager) RemoveSubscription(subscriptionID string) error {
	response, err := esm.helixClient.Client.RemoveEventSubSubscription(subscriptionID)
	if err != nil {
		return err
	}

	log.Infof("[%d] removed EventSubSubscription", response.StatusCode)
	esm.db.RemoveEventSubSubscription(subscriptionID)

	return nil
}
