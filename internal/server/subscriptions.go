package server

import (
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/log"
)

type SubscribtionStatus struct {
	Predictions bool `json:"predictions"`
}

func (a *Api) SubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	authResp, _, apiErr := a.authClient.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID := authResp.Data.UserID

	if r.URL.Query().Get("managing") != "" {
		userID, apiErr = a.userAdmin.CheckEditor(r, a.userAdmin.GetUserConfig(userID))
		if apiErr != nil {
			http.Error(w, apiErr.Error(), apiErr.Status())
			return
		}
	}

	if r.Method == http.MethodPut {
		a.eventsubManager.SubscribePredictions(userID)

		api.WriteJson(w, "ok", http.StatusOK)
	} else if r.Method == http.MethodDelete {
		for _, sub := range a.db.GetAllPredictionSubscriptions(userID) {
			log.Infof("Removing subscribtion on request %s from %s", sub.SubscriptionID, sub.TargetTwitchID)
			err := a.eventsubManager.RemoveEventSubSubscription(sub.SubscriptionID)
			if err != nil {
				log.Error(err)
			}
		}

		api.WriteJson(w, "ok", http.StatusOK)
	} else if r.Method == http.MethodGet {
		subs := a.db.GetAllPredictionSubscriptions(userID)
		log.Info(subs)

		hasPredictions := len(subs) > 0

		api.WriteJson(w, SubscribtionStatus{Predictions: hasPredictions}, http.StatusOK)
	}
}
