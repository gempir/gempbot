package election

import (
	"context"
	"fmt"
	"time"

	"github.com/gempir/gempbot/internal/bot"
	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/emotechief"
	"github.com/gempir/gempbot/internal/emoteservice"
	"github.com/gempir/gempbot/internal/eventsubsubscription"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/nicklaw5/helix/v2"
)

type ElectionManager struct {
	db            store.Store
	helixclient   helixclient.Client
	cpm           *channelpoint.ChannelPointManager
	esm           *eventsubsubscription.SubscriptionManager
	bot           *bot.Bot
	sevenTvClient emoteservice.ApiClient
}

func NewElectionManager(db store.Store, helixClient helixclient.Client, cpm *channelpoint.ChannelPointManager, esm *eventsubsubscription.SubscriptionManager, bot *bot.Bot, sevenTvClient emoteservice.ApiClient) *ElectionManager {
	return &ElectionManager{
		db:            db,
		helixclient:   helixClient,
		cpm:           cpm,
		esm:           esm,
		bot:           bot,
		sevenTvClient: sevenTvClient,
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
		if election.StartedRunAt == nil {
			log.Infof("Starting election for channel %s", election.ChannelTwitchID)
			em.startElection(election)
			time.Sleep(1 * time.Second)
		} else if election.StartedRunAt != nil && time.Since(*election.StartedRunAt) > time.Duration(election.Hours)*time.Hour {
			log.Infof("Stopping election for channel %s", election.ChannelTwitchID)
			em.stopElection(election)
			time.Sleep(1 * time.Second)
		}
	}
}

func (em *ElectionManager) stopElection(election store.Election) {
	nomination, err := em.db.GetTopVotedNominated(context.Background(), election.ChannelTwitchID, election.ID)
	if err != nil {
		log.Errorf("Failed to get top voted nomination %s", err.Error())
		return
	}

	err = em.sevenTvClient.AddEmote(election.ChannelTwitchID, nomination.EmoteID)
	if err != nil {
		log.Errorf("Failed to add emote %s", err.Error())
		return
	}

	election.StartedRunAt = nil
	err = em.db.CreateOrUpdateElection(context.Background(), election)
	if err != nil {
		log.Errorf("Failed to create/update election %s", err.Error())
		return
	}

	em.bot.SayByChannelID(election.ChannelTwitchID, fmt.Sprintf("🗳️ The emote %s has won the election with %d votes!", nomination.EmoteCode, nomination.Votes))
}

func (em *ElectionManager) startElection(election store.Election) {
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
		log.Errorf("Failed to create/updated reward %s", err.Error())
	}

	electionReward := &channelpoint.ElectionReward{TwitchRewardConfig: newReward, ElectionRewardAdditionalOptions: channelpoint.ElectionRewardAdditionalOptions{}}
	err = em.db.SaveReward(channelpoint.CreateStoreRewardFromReward(election.ChannelTwitchID, electionReward))
	if err != nil {
		log.Errorf("Failed to save reward %s", err.Error())
		return
	}

	em.esm.SubscribeRewardRedemptionAdd(election.ChannelTwitchID, newReward.ID)

	election.ChannelPointRewardID = newReward.ID
	time := time.Now()
	election.StartedRunAt = &time
	err = em.db.CreateOrUpdateElection(context.Background(), election)
	if err != nil {
		log.Errorf("Failed to create/update election %s", err.Error())
		return
	}

	em.bot.SayByChannelID(election.ChannelTwitchID, fmt.Sprintf("🗳️ A new Election has begun. Nominate a 7TV Emote with channel points. Every %d hours a new emote will be added to the channel. Each election will reset the nominations. The most voted one will be added to the channel.", election.Hours))
}

func (em *ElectionManager) Nominate(reward store.ChannelPointReward, redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent) {
	election, err := em.db.GetElection(context.Background(), reward.OwnerTwitchID)
	if err != nil {
		log.Errorf("failed to find election, refunding. %s", err.Error())
		err = em.helixclient.UpdateRedemptionStatus(reward.OwnerTwitchID, reward.RewardID, redemption.ID, false)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	emoteID, err := emotechief.GetSevenTvEmoteId(redemption.UserInput)
	if err != nil {
		log.Errorf("failed to parse emote, refunding. %s", err.Error())
		err = em.helixclient.UpdateRedemptionStatus(reward.OwnerTwitchID, reward.RewardID, redemption.ID, false)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	emote, err := em.sevenTvClient.GetEmote(emoteID)
	if err != nil {
		log.Errorf("failed to find emote, refunding. %s", err.Error())
		err = em.helixclient.UpdateRedemptionStatus(reward.OwnerTwitchID, reward.RewardID, redemption.ID, false)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	err = em.db.CreateOrIncrementNomination(context.Background(), store.Nomination{
		EmoteID:         emoteID,
		ChannelTwitchID: reward.OwnerTwitchID,
		ElectionID:      election.ID,
		EmoteCode:       emote.Code,
		NominatedBy:     redemption.UserID,
	})
	if err != nil {
		log.Errorf("failed to update nomination, refunding. %s", err.Error())
		err = em.helixclient.UpdateRedemptionStatus(reward.OwnerTwitchID, reward.RewardID, redemption.ID, false)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}
}
