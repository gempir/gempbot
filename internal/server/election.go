package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

func (a *Api) ElectionHandler(w http.ResponseWriter, r *http.Request) {
	var userID string
	if r.Method != http.MethodGet || a.authClient.HasAuth(r) {
		authResp, _, apiErr := a.authClient.AttemptAuth(r, w)
		if apiErr != nil {
			return
		}
		userID = authResp.Data.UserID

		if r.URL.Query().Get("managing") != "" {
			userID, apiErr = a.userAdmin.CheckEditor(r, a.userAdmin.GetUserConfig(userID))
			if apiErr != nil {
				http.Error(w, apiErr.Error(), apiErr.Status())
				return
			}
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

		user, err := a.helixClient.GetUserByUserID(newElection.ChannelTwitchID)
		if err != nil {
			log.Errorf("Failed to get user %s", err.Error())
		}

		reward := channelpoint.TwitchRewardConfig{
			Enabled:                           true,
			Title:                             "Nominate a 7TV Emote",
			Prompt:                            fmt.Sprintf("Nominate a 7TV Emote (Link) for the next election. Top voted nominations will be added. https://bot.gempir.com/nominations/%s", user.Login),
			Cost:                              newElection.NominationCost,
			IsUserInputRequired:               true,
			BackgroundColor:                   "#29D8F6",
			IsMaxPerStreamEnabled:             false,
			IsMaxPerUserPerStreamEnabled:      false,
			MaxPerStream:                      0,
			MaxPerUserPerStream:               0,
			IsGlobalCooldownEnabled:           false,
			ShouldRedemptionsSkipRequestQueue: false,
		}

		prevReward, err := a.db.GetChannelPointReward(newElection.ChannelTwitchID, dto.REWARD_ELECTION)
		if err == nil {
			log.Infof("Found previous reward %s", prevReward.RewardID)
		}

		newReward, err := a.channelPointManager.CreateOrUpdateChannelPointReward(newElection.ChannelTwitchID, reward, prevReward.RewardID)
		if err != nil {
			log.Errorf("Failed to create/updated reward %s", err.Error())
			if strings.Contains(err.Error(), "The broadcaster doesn't have partner or affiliate status") {
				err := a.db.DeleteElection(context.Background(), newElection.ChannelTwitchID)
				if err != nil {
					log.Errorf("Failed to delete election %s", err.Error())
					http.Error(w, fmt.Sprintf("Failed to delete election %s", err.Error()), http.StatusInternalServerError)
					return
				}
				http.Error(w, fmt.Sprintf("Deleted election because channel is not partner/affiliate: %s", newElection.ChannelTwitchID), http.StatusBadRequest)
				log.Infof("Deleted election because channel is not partner/affiliate: %s", newElection.ChannelTwitchID)
				return
			}
			http.Error(w, "Failed to save reward "+err.Error(), http.StatusInternalServerError)
			return
		}

		electionReward := &channelpoint.ElectionReward{TwitchRewardConfig: newReward, ElectionRewardAdditionalOptions: channelpoint.ElectionRewardAdditionalOptions{}}
		err = a.db.SaveReward(channelpoint.CreateStoreRewardFromReward(newElection.ChannelTwitchID, electionReward))
		if err != nil {
			log.Errorf("Failed to save reward %s", err.Error())
			http.Error(w, "Failed to save reward "+err.Error(), http.StatusInternalServerError)
			return
		}
		a.eventsubManager.SubscribeRewardRedemptionAdd(newElection.ChannelTwitchID, newReward.ID)

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
