package prediction

import (
	"encoding/json"
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
	nickHelix "github.com/nicklaw5/helix/v2"
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

	log.Infof("[helix] %d Created Prediction %s Error: %s", resp.StatusCode, userID, resp.ErrorMessage)

	api.WriteJson(w, resp.ErrorMessage, resp.StatusCode)
}
