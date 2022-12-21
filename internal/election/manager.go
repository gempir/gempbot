package election

import (
	"context"
	"fmt"
	"time"

	"github.com/gempir/gempbot/internal/bot"
	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/eventsubsubscription"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/nicklaw5/helix/v2"
)

type ElectionManager struct {
	db          store.Store
	helixclient helixclient.Client
	cpm         *channelpoint.ChannelPointManager
	esm         *eventsubsubscription.SubscriptionManager
	bot         *bot.Bot
}

func NewElectionManager(db store.Store, helixClient helixclient.Client, cpm *channelpoint.ChannelPointManager, esm *eventsubsubscription.SubscriptionManager, bot *bot.Bot) *ElectionManager {
	return &ElectionManager{
		db:          db,
		helixclient: helixClient,
		cpm:         cpm,
		esm:         esm,
		bot:         bot,
	}
}

func (em *ElectionManager) StartElectionManagerRoutine() {
	em.checkElections()
	// for range time.NewTicker(1 * time.Second).C {
	// 	em.checkElections()
	// }
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
	reward := channelpoint.TwitchRewardConfig{
		Enabled:                           true,
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
	if err != nil {
		log.Error(err.Error())
		return
	}

	electionReward := &channelpoint.ElectionReward{TwitchRewardConfig: newReward, ElectionRewardAdditionalOptions: channelpoint.ElectionRewardAdditionalOptions{}}
	err = em.db.SaveReward(channelpoint.CreateStoreRewardFromReward(election.ChannelPointRewardID, electionReward))
	if err != nil {
		log.Error(err.Error())
		return
	}

	em.esm.SubscribeRewardRedemptionAdd(election.ChannelTwitchID, newReward.ID)

	election.ChannelPointRewardID = newReward.ID
	// debug code, enable later
	// time := time.Now()
	// election.LastRunAt = &time
	err = em.db.CreateOrUpdateElection(context.Background(), election)
	if err != nil {
		log.Error(err.Error())
		return
	}

	em.bot.SayByChannelID(election.ChannelTwitchID, fmt.Sprintf("🗳️ A new Election has begun. Nominate a 7TV Emote with channel points. Every %d hours a new emote will be added to the channel. Each election will reset the nominations. The most voted one will be added to the channel.", election.Hours))
}

func (em *ElectionManager) Nominate(reward store.ChannelPointReward, redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent) {
	log.Infof("nominate %s %s", redemption.UserLogin, redemption.UserInput)
}
