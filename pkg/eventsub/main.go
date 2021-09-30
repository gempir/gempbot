package eventsub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/log"
	nickHelix "github.com/nicklaw5/helix"
)

type EventSubManager struct {
	cfg *config.Config
}

func NewEventSubManager(cfg *config.Config) *EventSubManager {
	return &EventSubManager{cfg: cfg}
}

func (s *EventSubManager) HandleWebhook(w http.ResponseWriter, r *http.Request, response interface{}) (done bool, apiErr api.Error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return true, api.NewApiError(http.StatusBadRequest, err)
	}

	verified := nickHelix.VerifyEventSubNotification(s.cfg.Secret, r.Header, string(body))
	if !verified {
		log.Errorf("Failed verification %s", r.Header.Get("Twitch-Eventsub-Message-Id"))
		return true, api.NewApiError(http.StatusPreconditionFailed, fmt.Errorf("failed verfication"))
	}

	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		return true, s.handleChallenge(w, r, body)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return true, api.NewApiError(http.StatusPreconditionFailed, fmt.Errorf("failed decoding body"+err.Error()))
	}

	api.WriteText(w, "ok", http.StatusOK)

	return false, nil
}

func (s *EventSubManager) handleChallenge(w http.ResponseWriter, r *http.Request, body []byte) api.Error {
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

func (s *EventSubManager) HandleChannelPointsCustomRewardRedemption(redemption nickHelix.EventSubChannelPointsCustomRewardRedemptionEvent) {
	log.Infof("Channel points custom reward redemption: %s", redemption.ID)
}
