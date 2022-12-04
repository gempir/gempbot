package election

import (
	"context"
	"fmt"
	"time"

	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

type ElectionManager struct {
	db          store.Store
	helixclient helixclient.Client
	cpm         *channelpoint.ChannelPointManager
}

func NewElectionManager(db store.Store, helixClient helixclient.Client, cpm *channelpoint.ChannelPointManager) *ElectionManager {

	return &ElectionManager{
		db:          db,
		helixclient: helixClient,
		cpm:         cpm,
	}
}

func (em *ElectionManager) StartElectionManagerRoutine() {
	for range time.NewTicker(1 * time.Minute).C {
		em.checkElections()
	}
}

func (em *ElectionManager) checkElections() {
	elections, err := em.db.GetAllElections(context.Background())
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, election := range elections {
		if election.LastRunAt == nil || time.Since(*election.LastRunAt) > time.Duration(election.Hours)*time.Hour {
			log.Infof("stopping any previous election and starting election for channel %s", election.ChannelTwitchID)
			em.runElection(election)
			time.Sleep(1 * time.Second)
		}
	}
}

func (em *ElectionManager) runElection(election store.Election) {
	token, err := em.db.GetUserAccessToken(election.ChannelTwitchID)
	if err != nil {
		log.Error(err.Error())
		return
	}

	reward := channelpoint.TwitchRewardConfig{
		Title:                             "Nominate a 7TV Emote",
		Prompt:                            fmt.Sprintf("Nominate a 7TV Emote for the next election. Every %d hours a new emote will be added to the channel. Each election will reset the nominations. The most voted one will be added to the channel.", election.Hours),
		Cost:                              election.NominationCost,
		IsUserInputRequired:               true,
		BackgroundColor:                   "#29D8F6",
		IsMaxPerStreamEnabled:             false,
		IsMaxPerUserPerStreamEnabled:      false,
		MaxPerStream:                      0,
		MaxPerUserPerStream:               0,
		IsGlobalCooldownEnabled:           false,
		ShouldRedemptionsSkipRequestQueue: false,
	}

	newReward, err := em.cpm.CreateOrUpdateChannelPointReward(election.ChannelTwitchID, reward, election.ChannelPointRewardID)

	// a.eventsubSubscriptionManager.SubscribeRewardRedemptionAdd(userID, config.ID)
	// if config.ApproveOnly {
	// 	a.eventsubSubscriptionManager.SubscribeRewardRedemptionUpdate(userID, config.ID)
	// }

	// err = em.db.SaveReward(channelpoint.CreateStoreRewardFromReward(election.ChannelPointRewardID, newReward))
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return
	// }

	election.ChannelPointRewardID = newReward.ID
	err = em.db.CreateOrUpdateElection(context.Background(), election)
	if err != nil {
		log.Error(err.Error())
		return
	}
}

func (em *ElectionManager) nominate(channelTwitchID string, userTwitchID string, input string) {
	log.Infof("nominate %s channel %s user: %s", input, channelTwitchID, userTwitchID)
}
