package election

import (
	"context"
	"fmt"
	"time"

	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

type ElectionManager struct {
	db          store.Store
	helixclient helixclient.Client
}

func NewElectionManager(db store.Store, helixClient helixclient.Client) *ElectionManager {

	return &ElectionManager{
		db:          db,
		helixclient: helixClient,
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

	reward := helixclient.CreateCustomRewardRequest{
		Title:                             "Nominate a 7TV Emote",
		Prompt:                            fmt.Sprintf("Nominate a 7TV Emote for the next election. Every %d hours a new emote will be added to the channel. Each election will reset the nominations. The most voted one will be added to the channel.", election.Hours),
		Cost:                              election.NominationCost,
		IsEnabled:                         true,
		IsUserInputRequired:               true,
		BackgroundColor:                   "#29D8F6",
		IsMaxPerStreamEnabled:             false,
		IsMaxPerUserPerStreamEnabled:      false,
		MaxPerStream:                      0,
		MaxPerUserPerStream:               0,
		IsGlobalCooldownEnabled:           false,
		GlobalCoolDownSeconds:             0,
		ShouldRedemptionsSkipRequestQueue: false,
	}

	resp, err := em.helixclient.CreateOrUpdateReward(election.ChannelTwitchID, token.AccessToken, reward, election.ChannelPointRewardID)
	if err != nil {
		log.Error(err.Error())
		return
	}

	election.ChannelPointRewardID = resp.ID
	err = em.db.CreateOrUpdateElection(context.Background(), election)
	if err != nil {
		log.Error(err.Error())
		return
	}
}
