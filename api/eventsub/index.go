package eventsub

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/emotechief"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/nicklaw5/helix/v2"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)
	helixClient := helixclient.NewClient(cfg, db)
	chatClient := chat.NewClient(cfg)
	go chatClient.Connect(func() {})
	emoteChief := emotechief.NewEmoteChief(cfg, db, helixClient, chatClient)
	eventSubManager := eventsub.NewEventSubManager(cfg, helixClient, db, emoteChief, chatClient)

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
