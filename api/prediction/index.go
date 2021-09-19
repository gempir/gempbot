package prediction

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
	nickHelix "github.com/nicklaw5/helix"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)
	helixClient := helix.NewClient(cfg, db)
	auth := auth.NewAuth(cfg, db, helixClient)
	userAdmin := user.NewUserAdmin(cfg, db, helixClient, nil)

	userID := ""
	authResp, _, apiErr := auth.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID = authResp.Data.UserID

	if r.URL.Query().Get("managing") != "" {
		userID, apiErr = userAdmin.CheckEditor(r, userAdmin.GetUserConfig(userID))
		if apiErr != nil {
			http.Error(w, apiErr.Error(), apiErr.Status())
			return
		}
	}

	prediction := &nickHelix.CreatePredictionParams{}
	err := json.NewDecoder(r.Body).Decode(prediction)
	if err != nil {
		http.Error(w, "invalid create prediction body", http.StatusBadRequest)
		return
	}

	token, err := db.GetUserAccessToken(userID)
	if err != nil {
		http.Error(w, "no api token, broadcaster needs to login again in dashboard", http.StatusUnauthorized)
		return
	}

	helixClient.Client.SetUserAccessToken(token.AccessToken)
	resp, err := helixClient.Client.CreatePrediction(prediction)
	if err != nil {
		http.Error(w, resp.ErrorMessage, resp.StatusCode)
		return
	}

	for _, prediction := range resp.Data.Predictions {
		err := db.SavePrediction(store.PredictionLog{ID: prediction.ID, OwnerTwitchID: userID, Title: prediction.Title, StartedAt: prediction.CreatedAt.Time, LockedAt: nil, EndedAt: nil, WinningOutcomeID: prediction.WinningOutcomeID})
		if err != nil {
			log.Error(err)
		}

		for _, outcome := range prediction.Outcomes {
			err := db.SaveOutcome(store.PredictionLogOutcome{ID: outcome.ID, PredictionID: prediction.ID, Title: outcome.Title, Color: strings.ToLower(outcome.Color), Users: outcome.Users, ChannelPoints: outcome.ChannelPoints})
			if err != nil {
				log.Error(err)
			}
		}
	}
	log.Infof("[helix] %d Created Prediction %s", resp.StatusCode, userID)

	api.WriteJson(w, resp.ErrorMessage, resp.StatusCode)
}
