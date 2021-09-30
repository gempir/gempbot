package eventsub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	nickHelix "github.com/nicklaw5/helix"
)

type EventSubManager struct {
	cfg         *config.Config
	helixClient *helix.Client
	db          *store.Database
}

func NewEventSubManager(cfg *config.Config, helixClient *helix.Client, db *store.Database) *EventSubManager {
	return &EventSubManager{cfg: cfg, helixClient: helixClient, db: db}
}

func (esm *EventSubManager) HandleWebhook(w http.ResponseWriter, r *http.Request, response interface{}) (done bool, apiErr api.Error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return true, api.NewApiError(http.StatusBadRequest, err)
	}

	verified := nickHelix.VerifyEventSubNotification(esm.cfg.Secret, r.Header, string(body))
	if !verified {
		log.Errorf("Failed verification %s", r.Header.Get("Twitch-Eventsub-Message-Id"))
		return true, api.NewApiError(http.StatusPreconditionFailed, fmt.Errorf("failed verfication"))
	}

	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		return true, esm.handleChallenge(w, r, body)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return true, api.NewApiError(http.StatusPreconditionFailed, fmt.Errorf("failed decoding body"+err.Error()))
	}

	api.WriteText(w, "ok", http.StatusOK)

	return false, nil
}

func (esm *EventSubManager) handleChallenge(w http.ResponseWriter, r *http.Request, body []byte) api.Error {
	var event struct {
		Challenge string `json:"challenge"`
	}
	err := json.Unmarshal(body, &event)
	if err != nil {
		return api.NewApiError(http.StatusBadRequest, fmt.Errorf("Failed to handle challenge: "+err.Error()))
	}

	log.Infof("Challenge success: %s", event.Challenge)
	api.WriteText(w, event.Challenge, http.StatusOK)
	return nil
}

func (esm *EventSubManager) HandleChannelPointsCustomRewardRedemption(redemption nickHelix.EventSubChannelPointsCustomRewardRedemptionEvent) {
	log.Infof("Channel points custom reward redemption: %s", redemption.ID)
}

func (esm *EventSubManager) SubscribeChannelPoints(userID string) {
	response, err := esm.helixClient.CreateEventSubSubscription(userID, esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+nickHelix.EventSubTypeChannelPointsCustomRewardRedemptionAdd, nickHelix.EventSubTypeChannelPointsCustomRewardRedemptionAdd)
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
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type)
	}
}

func (esm *EventSubManager) RemoveEventSubSubscription(userID string, subscriptionID string, subType string, reason string) error {
	response, err := esm.helixClient.Client.RemoveEventSubSubscription(subscriptionID)
	if err != nil {
		return err
	}

	log.Infof("[%d] removed EventSubSubscription for %s reason: %s", response.StatusCode, userID, reason)
	esm.db.RemoveEventSubSubscription(userID, subscriptionID, subType)

	return nil
}

func (esm *EventSubManager) RemoveAllEventSubSubscriptions(userID string) {
	// @TODO rework using the DB so we don't need to query literally every sub
	resp, err := esm.helixClient.Client.GetEventSubSubscriptions(&nickHelix.EventSubSubscriptionsParams{})
	if err != nil {
		log.Errorf("Failed to get subscriptions: %s", err)
		return
	}

	subscriptions := resp.Data.EventSubSubscriptions

	for {
		cursor := resp.Data.Pagination.Cursor
		if cursor == "" {
			break
		}
		log.Infof("Getting next subscriptions cursor: %s", cursor)

		nextResp, err := esm.helixClient.Client.GetEventSubSubscriptions(&nickHelix.EventSubSubscriptionsParams{})
		if err != nil {
			log.Errorf("Failed to get subscriptions: %s", err)
		}

		subscriptions = append(subscriptions, nextResp.Data.EventSubSubscriptions...)
	}

	for _, sub := range subscriptions {
		if sub.Condition.BroadcasterUserID != userID {
			continue
		}

		err := esm.RemoveEventSubSubscription(userID, sub.ID, sub.Type, "removed all subscriptions")
		if err != nil {
			log.Error(err)
			return
		}
	}
}
