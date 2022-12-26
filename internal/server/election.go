package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

func (a *Api) ElectionHandler(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodGet {
		channel := r.URL.Query().Get("channel")
		if channel != "" {
			user, err := a.helixClient.GetUserByUsername(channel)
			if err != nil {
				http.Error(w, "user not found", http.StatusBadRequest)
				return
			}
			userID = user.ID
		}

		election, err := a.db.GetElection(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		api.WriteJson(w, election, http.StatusOK)
	} else if r.Method == http.MethodPost {

		newElection, err := readElectionBody(r)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		prevElection, err := a.db.GetElection(r.Context(), userID)
		if err == nil {
			newElection.StartedRunAt = prevElection.StartedRunAt
			newElection.CreatedAt = prevElection.CreatedAt
			newElection.UpdatedAt = time.Now()
		} else {
			newElection.StartedRunAt = nil
			newElection.CreatedAt = time.Now()
			newElection.UpdatedAt = time.Now()
		}

		newElection.ChannelTwitchID = userID
		err = a.db.CreateOrUpdateElection(r.Context(), newElection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		api.WriteJson(w, nil, http.StatusOK)
	} else if r.Method == http.MethodDelete {
		err := a.db.DeleteElection(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = a.db.ClearNominations(r.Context(), userID)
		if err != nil {
			log.Warnf("failed to clear nominations, probably ok: %s", err.Error())
		}

		err = a.channelPointManager.DeleteElectionReward(userID)
		if err != nil {
			log.Warnf("failed to delete election reward, probably ok: %s", err.Error())
		}

		api.WriteJson(w, nil, http.StatusOK)
	}
}

func readElectionBody(r *http.Request) (store.Election, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return store.Election{}, err
	}

	var data store.Election
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return store.Election{}, err
	}

	return data, nil
}
