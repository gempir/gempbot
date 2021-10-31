package reward

import (
	"fmt"
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/channelpoint"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/dto"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
)

func HandlerBttv(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)
	helixClient := helix.NewClient(cfg, db)
	auth := auth.NewAuth(cfg, db, helixClient)
	userAdmin := user.NewUserAdmin(cfg, db, helixClient, nil)
	cpm := channelpoint.NewChannelPointManager(cfg, helixClient, db)
	subscriptionManager := eventsub.NewSubscriptionManager(cfg, db, helixClient)

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

	if r.Method == http.MethodGet {
		reward, err := db.GetChannelPointReward(userID, dto.RewardType(r.URL.Query().Get("type")))
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
		reward, err := db.GetChannelPointReward(userID, newReward.GetType())
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

		err = db.SaveReward(channelpoint.CreateStoreRewardFromReward(userID, newReward))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed saving reward: %s", err), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodDelete {

		reward, err := db.GetChannelPointReward(userID, dto.RewardType(r.URL.Query().Get("type")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		token, err := db.GetUserAccessToken(userID)
		if err != nil {
			http.Error(w, "no accessToken to edit reward", http.StatusNotFound)
		}

		db.DeleteChannelPointReward(userID, dto.RewardType(r.URL.Query().Get("type")))

		err = helixClient.DeleteReward(userID, token.AccessToken, reward.RewardID)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
