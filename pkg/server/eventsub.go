package server

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/emotechief"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/nicklaw5/helix/v2"
)

func (a *Api) EventSubHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("New event sub request %s %s", r.Method, r.URL.Path)
	emoteChief := emotechief.NewEmoteChief(a.cfg, a.db, a.helixClient, a.bot.ChatClient)
	eventSubManager := eventsub.NewEventSubManager(a.cfg, a.helixClient, a.db, emoteChief, a.bot.ChatClient)

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
