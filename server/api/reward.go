package api

import (
	"fmt"
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/channelpoint"
	"github.com/gempir/gempbot/pkg/dto"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/log"
)

func (a *Api) RewardHandler(w http.ResponseWriter, r *http.Request) {
	cpm := channelpoint.NewChannelPointManager(a.cfg, a.helixClient, a.db)
	subscriptionManager := eventsub.NewSubscriptionManager(a.cfg, a.db, a.helixClient)

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
		reward, err := a.db.GetChannelPointReward(userID, dto.RewardType(r.URL.Query().Get("type")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		api.WriteJson(w, reward, http.StatusOK)
	} else if r.Method == http.MethodPost {
		newReward, err := channelpoint.CreateRewardFromBody(r.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rewardID := ""
		reward, err := a.db.GetChannelPointReward(userID, newReward.GetType())
		if err == nil {
			rewardID = reward.RewardID
		}

		config, err := cpm.CreateOrUpdateChannelPointReward(userID, newReward.GetConfig(), rewardID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed saving reward to twitch: %s", err), http.StatusInternalServerError)
			return
		}

		subscriptionManager.SubscribeRewardRedemptionAdd(userID, config.ID)
		if config.ApproveOnly {
			subscriptionManager.SubscribeRewardRedemptionUpdate(userID, config.ID)
		}

		newReward.SetConfig(config)

		err = a.db.SaveReward(channelpoint.CreateStoreRewardFromReward(userID, newReward))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed saving reward: %s", err), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodDelete {

		reward, err := a.db.GetChannelPointReward(userID, dto.RewardType(r.URL.Query().Get("type")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		token, err := a.db.GetUserAccessToken(userID)
		if err != nil {
			http.Error(w, "no accessToken to edit reward", http.StatusNotFound)
		}

		a.db.DeleteChannelPointReward(userID, dto.RewardType(r.URL.Query().Get("type")))

		err = a.helixClient.DeleteReward(userID, token.AccessToken, reward.RewardID)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
