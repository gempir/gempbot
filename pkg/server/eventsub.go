package server

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/log"
	"github.com/nicklaw5/helix/v2"
)

func (a *Api) EventSubHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("New event sub request %s %s", r.Method, r.URL.Path)

	event, err := a.eventsubManager.HandleWebhook(w, r)
	if err != nil || len(event) == 0 {
		if err != nil {
			http.Error(w, err.Error(), err.Status())
		}
		return
	}

	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd {
		a.eventsubManager.HandleChannelPointsCustomRewardRedemption(event)
		return
	}
	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPointsCustomRewardRedemptionUpdate {
		a.eventsubManager.HandleChannelPointsCustomRewardRedemption(event)
		return
	}
	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPredictionBegin {
		a.eventsubManager.HandlePredictionBegin(event)
		return
	}
	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPredictionLock {
		a.eventsubManager.HandlePredictionLock(event)
		return
	}
	if r.URL.Query().Get("type") == helix.EventSubTypeChannelPredictionEnd {
		a.eventsubManager.HandlePredictionEnd(event)
		return
	}

	http.Error(w, "Invalid event type", http.StatusBadRequest)
}
