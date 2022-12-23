package server

import (
	"encoding/json"
	"io"
	"net/http"

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
		reward, err := a.db.GetElection(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		api.WriteJson(w, reward, http.StatusOK)
	} else if r.Method == http.MethodPost {
		newElection, err := readElectionBody(r)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newElection.ChannelTwitchID = userID
		err = a.db.CreateOrUpdateElection(r.Context(), newElection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		api.WriteJson(w, nil, http.StatusOK)
	} else if r.Method == http.MethodDelete {
		election, err := a.db.GetElection(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = a.db.DeleteElection(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if election.ChannelTwitchID != "" {
			err = a.channelPointManager.DeleteElectionReward(userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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