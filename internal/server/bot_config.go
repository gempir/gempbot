package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

type BotConfigRequest struct {
	ChannelID               string `json:"channelId"`
	PredictionAnnouncements *bool  `json:"predictionAnnouncements,omitempty"`
}

type BotConfigResponse struct {
	PredictionAnnouncements bool `json:"predictionAnnouncements"`
}

func (a *Api) BotConfigHandler(w http.ResponseWriter, r *http.Request) {
	authResp, _, err := a.authClient.AttemptAuth(r, w)
	if err != nil {
		return
	}

	userID := authResp.Data.UserID

	// Check if managing another channel
	if r.URL.Query().Get("managing") != "" {
		userConfig := a.userAdmin.GetUserConfig(authResp.Data.UserID)
		ownerUserID, err := a.userAdmin.CheckEditor(r, userConfig)
		if err != nil {
			http.Error(w, err.Error(), err.Status())
			return
		}
		userID = ownerUserID
	}

	if r.Method == http.MethodGet {
		botConfig, err := a.db.GetBotConfig(userID)
		if err != nil {
			// Return default config if not found
			botConfig = store.BotConfig{
				OwnerTwitchID:           userID,
				PredictionAnnouncements: false,
			}
		}

		response := BotConfigResponse{
			PredictionAnnouncements: botConfig.PredictionAnnouncements,
		}

		api.WriteJson(w, response, http.StatusOK)
		return

	} else if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Failed reading bot config body: %s", err)
			http.Error(w, "Failure reading body "+err.Error(), http.StatusInternalServerError)
			return
		}

		var req BotConfigRequest
		if err := json.Unmarshal(body, &req); err != nil {
			log.Errorf("Failed unmarshalling bot config: %s", err)
			http.Error(w, "Failure unmarshalling config "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get existing config or create new one
		botConfig, err := a.db.GetBotConfig(userID)
		if err != nil {
			botConfig = store.BotConfig{
				OwnerTwitchID: userID,
			}
		}

		previousPredictionAnnouncements := botConfig.PredictionAnnouncements

		// Update fields
		if req.PredictionAnnouncements != nil {
			botConfig.PredictionAnnouncements = *req.PredictionAnnouncements
		}

		err = a.db.SaveBotConfig(r.Context(), botConfig)
		if err != nil {
			log.Errorf("failed saving bot config: %s", err)
			http.Error(w, "Failed to save config", http.StatusInternalServerError)
			return
		}

		// Handle prediction announcements subscription
		if req.PredictionAnnouncements != nil && previousPredictionAnnouncements != *req.PredictionAnnouncements {
			if *req.PredictionAnnouncements {
				// Subscribe to prediction events
				log.Infof("Subscribing to predictions for user %s", userID)
				a.eventsubManager.SubscribePredictions(userID)
			} else {
				// Unsubscribe from prediction events
				log.Infof("Unsubscribing from predictions for user %s", userID)
				subscriptions := a.db.GetAllPredictionSubscriptions(userID)
				for _, sub := range subscriptions {
					a.eventsubManager.RemoveSubscription(sub.SubscriptionID)
				}
			}
		}

		response := BotConfigResponse{
			PredictionAnnouncements: botConfig.PredictionAnnouncements,
		}

		api.WriteJson(w, response, http.StatusOK)
		return
	}

	http.Error(w, "unknown method", http.StatusMethodNotAllowed)
}
