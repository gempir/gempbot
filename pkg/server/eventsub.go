package server

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/emotechief"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/nicklaw5/helix/v2"
)

func (a *Api) EventSubHandler(w http.ResponseWriter, r *http.Request) {
	chatClient := chat.NewClient(a.cfg)
	go chatClient.Connect(func() {})
	emoteChief := emotechief.NewEmoteChief(a.cfg, a.db, a.helixClient, chatClient)
	eventSubManager := eventsub.NewEventSubManager(a.cfg, a.helixClient, a.db, emoteChief, chatClient)

	event, err := eventSubManager.HandleWebhook(w, r)
	if err != nil || len(event) == 0 {
		if err != nil {
			http.Error(w, err.Error(), err.Status())
		}
		return
	}

	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd {
		eventSubManager.HandleChannelPointsCustomRewardRedemption(event)
		return
	}
	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPointsCustomRewardRedemptionUpdate {
		eventSubManager.HandleChannelPointsCustomRewardRedemption(event)
		return
	}
	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPredictionBegin {
		eventSubManager.HandlePredictionBegin(event)
		return
	}
	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPredictionLock {
		eventSubManager.HandlePredictionLock(event)
		return
	}
	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPredictionEnd {
		eventSubManager.HandlePredictionEnd(event)
		return
	}

	http.Error(w, "Invalid event type", http.StatusBadRequest)
}
