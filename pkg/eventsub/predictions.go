package eventsub

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gempir/gempbot/pkg/dto"
	"github.com/gempir/gempbot/pkg/humanize"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	nickHelix "github.com/nicklaw5/helix/v2"
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

func (esm *EventSubManager) HandlePredictionBegin(event []byte) {
	var data nickHelix.EventSubChannelPredictionBeginEvent
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

	err = esm.db.SavePrediction(store.PredictionLog{ID: data.ID, OwnerTwitchID: data.BroadcasterUserID, Title: data.Title, StartedAt: data.StartedAt.Time, LockedAt: lockedTime})
	if err != nil {
		log.Error(err)
	}

	for _, outcome := range data.Outcomes {
		err = esm.db.SaveOutcome(store.PredictionLogOutcome{ID: outcome.ID, PredictionID: data.ID, Title: outcome.Title, Color: outcome.Color})
		if err != nil {
			log.Error(err)
		}
	}

	esm.chatClient.Say(
		data.BroadcasterUserLogin,
		fmt.Sprintf("PogChamp prediction: %s [ %s | %s ] ending in %s",
			data.Title,
			data.Outcomes[0].Title,
			data.Outcomes[1].Title,
			humanize.TimeUntil(data.StartedAt.Time, data.LocksAt.Time),
		),
	)
}

func (esm *EventSubManager) HandlePredictionLock(event []byte) {
	var data nickHelix.EventSubChannelPredictionLockEvent
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

	err = esm.db.SavePrediction(store.PredictionLog{ID: data.ID, OwnerTwitchID: data.BroadcasterUserID, Title: data.Title, StartedAt: data.StartedAt.Time, LockedAt: lockedTime})
	if err != nil {
		log.Error(err)
	}

	for _, outcome := range data.Outcomes {
		err = esm.db.SaveOutcome(store.PredictionLogOutcome{ID: outcome.ID, PredictionID: data.ID, Title: outcome.Title, Color: outcome.Color, Users: outcome.Users, ChannelPoints: outcome.ChannelPoints})
		if err != nil {
			log.Error(err)
		}
	}

	esm.chatClient.Say(
		data.BroadcasterUserLogin,
		fmt.Sprintf("FBtouchdown locked submissions for: %s",
			data.Title,
		),
	)
}

func (esm *EventSubManager) HandlePredictionEnd(event []byte) {
	var data nickHelix.EventSubChannelPredictionEndEvent
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

	err = esm.db.SavePrediction(store.PredictionLog{ID: data.ID, OwnerTwitchID: data.BroadcasterUserID, Title: data.Title, StartedAt: data.StartedAt.Time, EndedAt: endTime, WinningOutcomeID: data.WinningOutcomeID, Status: data.Status})
	if err != nil {
		log.Error(err)
	}

	var winningOutcome store.PredictionLogOutcome

	for _, outcome := range data.Outcomes {
		outcomeModel := store.PredictionLogOutcome{ID: outcome.ID, PredictionID: data.ID, Title: outcome.Title, Color: outcome.Color, Users: outcome.Users, ChannelPoints: outcome.ChannelPoints}

		if data.WinningOutcomeID == outcome.ID {
			winningOutcome = outcomeModel
		}

		err = esm.db.SaveOutcome(outcomeModel)
		if err != nil {
			log.Error(err)
		}
	}

	if strings.ToUpper(data.Status) == dto.PredictionStatusCanceled {
		esm.chatClient.Say(
			data.BroadcasterUserLogin,
			fmt.Sprintf("NinjaGrumpy canceled prediction: %s",
				data.Title,
			),
		)
	} else {
		esm.chatClient.Say(
			data.BroadcasterUserLogin,
			fmt.Sprintf("PogChamp ended prediction: %s Winner: %s %s",
				data.Title,
				winningOutcome.GetColorEmoji(),
				winningOutcome.Title,
			),
		)
	}
}
