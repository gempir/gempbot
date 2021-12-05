package subscriptions

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/emotechief"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
)

type SubscribtionStatus struct {
	Predictions bool `json:"predictions"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)
	helixClient := helixclient.NewClient(cfg, db)
	auth := auth.NewAuth(cfg, db, helixClient)
	userAdmin := user.NewUserAdmin(cfg, db, helixClient, nil)
	chatClient := chat.NewClient(cfg)
	go chatClient.Connect(func() {})
	emoteChief := emotechief.NewEmoteChief(cfg, db, helixClient, chatClient)
	eventSubManager := eventsub.NewEventSubManager(cfg, helixClient, db, emoteChief, chatClient)

	authResp, _, apiErr := auth.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID := authResp.Data.UserID

	if r.URL.Query().Get("managing") != "" {
		userID, apiErr = userAdmin.CheckEditor(r, userAdmin.GetUserConfig(userID))
		if apiErr != nil {
			http.Error(w, apiErr.Error(), apiErr.Status())
			return
		}
	}

	if r.Method == http.MethodPut {
		eventSubManager.SubscribePredictions(userID)

		api.WriteJson(w, "ok", http.StatusOK)
	} else if r.Method == http.MethodDelete {
		for _, sub := range db.GetAllPredictionSubscriptions(userID) {
			log.Infof("Removing subscribtion on request %s from %s", sub.SubscriptionID, sub.TargetTwitchID)
			err := eventSubManager.RemoveEventSubSubscription(sub.SubscriptionID)
			if err != nil {
				log.Error(err)
			}
		}

		api.WriteJson(w, "ok", http.StatusOK)
	} else if r.Method == http.MethodGet {
		subs := db.GetAllPredictionSubscriptions(userID)
		log.Info(subs)

		hasPredictions := len(subs) > 0

		api.WriteJson(w, SubscribtionStatus{Predictions: hasPredictions}, http.StatusOK)
	}
}
