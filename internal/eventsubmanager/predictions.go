package eventsubmanager

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/humanize"
	"github.com/gempir/gempbot/internal/log"
	"github.com/nicklaw5/helix/v2"
)

func (esm *EventsubManager) SubscribePredictions(userID string) {
	esm.SubscribePredictionsBegin(userID)
	esm.SubscribePredictionsLock(userID)
	esm.SubscribePredictionsEnd(userID)
}

func (esm *EventsubManager) SubscribePredictionsBegin(userID string) {
	response, err := esm.helixClient.CreateEventSubSubscription(userID, esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+helix.EventSubTypeChannelPredictionBegin, "channel.prediction.begin")
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	log.Infof("[%d] created subscription %s", response.StatusCode, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new sub in %s %s", userID, sub.Type)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type, "")
	}
}

func (esm *EventsubManager) SubscribePredictionsLock(userID string) {
	response, err := esm.helixClient.CreateEventSubSubscription(userID, esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+helix.EventSubTypeChannelPredictionLock, "channel.prediction.lock")
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	log.Infof("[%d] created subscription %s", response.StatusCode, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new sub in %s %s", userID, sub.Type)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type, "")
	}
}

func (esm *EventsubManager) SubscribePredictionsEnd(userID string) {
	response, err := esm.helixClient.CreateEventSubSubscription(userID, esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+helix.EventSubTypeChannelPredictionEnd, "channel.prediction.end")
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	log.Infof("[%d] created subscription %s", response.StatusCode, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new sub in %s %s", userID, sub.Type)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type, "")
	}
}

func (esm *EventsubManager) HandlePredictionBegin(event []byte) {
	var data helix.EventSubChannelPredictionBeginEvent
	err := json.Unmarshal(event, &data)
	if err != nil {
		log.Errorf("Failed to decode event: %s", err)
		return
	}

	log.Infof("predictionBegin %s", data.StartedAt)
	if data.ID == "" {
		return
	}

	lockedTime := &data.LocksAt.Time
	if lockedTime.IsZero() {
		lockedTime = nil
	}

	titles := []string{}
	for _, outcome := range data.Outcomes {
		titles = append(titles, outcome.Title)
	}

	esm.helixClient.SendChatMessage(
		data.BroadcasterUserID,
		fmt.Sprintf("PogChamp prediction: %s [ %s ] ending in %s",
			data.Title,
			strings.Join(titles, " | "),
			humanize.TimeUntil(data.StartedAt.Time, data.LocksAt.Time),
		),
	)
}

func (esm *EventsubManager) HandlePredictionLock(event []byte) {
	var data helix.EventSubChannelPredictionLockEvent
	err := json.Unmarshal(event, &data)
	if err != nil {
		log.Errorf("Failed to decode event: %s", err)
		return
	}

	log.Infof("predictionLock %s", data.LockedAt)
	if data.ID == "" {
		return
	}

	lockedTime := &data.LockedAt.Time
	if lockedTime.IsZero() {
		lockedTime = nil
	}

	esm.helixClient.SendChatMessage(
		data.BroadcasterUserID,
		fmt.Sprintf("FBtouchdown locked submissions for: %s",
			data.Title,
		),
	)
}

func (esm *EventsubManager) HandlePredictionEnd(event []byte) {
	var data helix.EventSubChannelPredictionEndEvent
	err := json.Unmarshal(event, &data)
	if err != nil {
		log.Errorf("Failed to decode event: %s", err)
		return
	}

	log.Infof("predictionEnd %s", data.Status)
	if data.ID == "" {
		return
	}

	endTime := &data.EndedAt.Time
	if endTime.IsZero() {
		endTime = nil
	}

	var winningOutcome helix.EventSubOutcome

	for _, outcome := range data.Outcomes {
		if data.WinningOutcomeID == outcome.ID {
			winningOutcome = outcome
		}
	}

	if strings.ToUpper(data.Status) == dto.PredictionStatusCanceled {
		esm.helixClient.SendChatMessage(
			data.BroadcasterUserID,
			fmt.Sprintf("NinjaGrumpy canceled prediction: %s",
				data.Title,
			),
		)
	} else {
		esm.helixClient.SendChatMessage(
			data.BroadcasterUserID,
			fmt.Sprintf("PogChamp ended prediction: %s Winner: %s %s",
				data.Title,
				getColorEmoji(winningOutcome),
				winningOutcome.Title,
			),
		)
	}
}

func getColorEmoji(outcome helix.EventSubOutcome) string {
	if outcome.Color == dto.Outcome_First {
		return "🟦"
	}

	return "🟪"
}
