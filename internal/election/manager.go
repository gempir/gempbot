package election

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gempir/gempbot/internal/bot"
	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/dto"
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

const CHECK_INTERVAL_SECONDS = 30

// Will check minutes before and after the specific time
const MINUTE_ROOM = 3

func (em *ElectionManager) StartElectionManagerRoutine() {
	for range time.NewTicker(CHECK_INTERVAL_SECONDS * time.Second).C {
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
			if election.SpecificTime != nil {
				now := time.Now()
				specificTime := time.Date(now.Year(), now.Month(), now.Day(), election.SpecificTime.Hour(), election.SpecificTime.Minute(), 0, 0, election.SpecificTime.Location())
				if specificTime.Add(time.Minute*MINUTE_ROOM).After(now) || specificTime.Sub(now) < -MINUTE_ROOM*time.Minute {
					continue
				}
			}
			log.Infof("Starting election for channel %s", election.ChannelTwitchID)
			em.startElection(election)
			time.Sleep(1 * time.Second)
		} else if election.StartedRunAt != nil && time.Since(*election.StartedRunAt) > (time.Duration(election.Hours)*time.Hour)-(time.Second*CHECK_INTERVAL_SECONDS) {
			if election.SpecificTime != nil {
				now := time.Now()
				specificTime := time.Date(now.Year(), now.Month(), now.Day(), election.SpecificTime.Hour(), election.SpecificTime.Minute(), 0, 0, election.SpecificTime.Location())
				if specificTime.Add(time.Minute*MINUTE_ROOM).After(now) || specificTime.Sub(now) < -MINUTE_ROOM*time.Minute {
					continue
				}
			}

			log.Infof("Stopping election for channel %s", election.ChannelTwitchID)
			em.stopElection(election)
			time.Sleep(1 * time.Second)
		}
	}
}

func (em *ElectionManager) stopElection(election store.Election) {
	log.Infof("Stopping election %v", election)
	nominations, err := em.db.GetNominations(context.Background(), election.ChannelTwitchID)
	if err != nil {
		log.Errorf("Failed to get top voted nomination %s", err.Error())
	}

	log.Infof("Nominations %v", nominations)

	nominationsAdded := []store.Nomination{}
	nominatedByList := []string{}
	for _, nomination := range nominations {
		isBlocked := em.db.IsEmoteBlocked(election.ChannelTwitchID, nomination.EmoteID, dto.REWARD_SEVENTV)
		if isBlocked {
			log.Errorf("Emote %s is blocked in channel %s", nomination.EmoteCode, election.ChannelTwitchID)
			continue
		}

		err = em.sevenTvClient.AddEmote(election.ChannelTwitchID, nomination.EmoteID)
		if err != nil {
			log.Errorf("Failed to add emote %s", err.Error())
			continue
		}
		em.db.AddEmoteLogEntry(context.Background(), store.EmoteLog{CreatedAt: time.Now(), EmoteID: nomination.EmoteID, AddedBy: nomination.NominatedBy, Type: dto.REWARD_ELECTION, EmoteCode: nomination.EmoteCode, ChannelTwitchID: election.ChannelTwitchID})
		nominationsAdded = append(nominationsAdded, nomination)
		nominatedByList = append(nominatedByList, nomination.NominatedBy)

		if (len(nominationsAdded)) >= election.EmoteAmount {
			break
		}
	}

	election.StartedRunAt = nil
	err = em.db.CreateOrUpdateElection(context.Background(), election)
	if err != nil {
		log.Errorf("Failed to create/update election %s", err.Error())
		return
	}

	err = em.db.ClearNominations(context.Background(), election.ChannelTwitchID)
	if err != nil {
		log.Errorf("Failed to clear nominations %s", err.Error())
	}

	users, err := em.helixclient.GetUsersByUserIds(nominatedByList)
	if err != nil {
		log.Errorf("Failed to get user %s", err.Error())
	}

	nominationStrings := []string{}
	for _, nomination := range nominationsAdded {
		var text string
		text += fmt.Sprintf("[⬆️%d|%d⬇️] %s ", len(nomination.Votes), len(nomination.Downvotes), nomination.EmoteCode)

		if val, ok := users[nomination.NominatedBy]; ok {
			text += fmt.Sprintf(" by %s", val.DisplayName)
		}
		nominationStrings = append(nominationStrings, text)
	}

	em.bot.SayByChannelID(election.ChannelTwitchID, fmt.Sprintf("🗳️ Election round over. Adding Emotes: %s", strings.Join(nominationStrings, ", ")))
}

func (em *ElectionManager) startElection(election store.Election) {
	time := time.Now()
	election.StartedRunAt = &time
	err := em.db.CreateOrUpdateElection(context.Background(), election)
	if err != nil {
		log.Errorf("Failed to create/update election %s", err.Error())
		return
	}

	userData, err := em.helixclient.GetUserByUserID(election.ChannelTwitchID)
	if err != nil {
		log.Errorf("Failed to get user data %s", err.Error())
		return
	}

	em.bot.Say(userData.Login, fmt.Sprintf("🗳️ A new Election has begun. Nominate a 7TV Emote with channel points. Top voted nominations will be added. Each election will reset the nominations. The most voted one will be added to the channel. Checkout bot.gempir.com/nominations/%s", userData.Login))
}

func (em *ElectionManager) Nominate(reward store.ChannelPointReward, redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent) {
	election, err := em.db.GetElection(context.Background(), reward.OwnerTwitchID)
	if err != nil {
		log.Errorf("failed to find election, refunding and deleting reward. %s", err.Error())
		err = em.helixclient.UpdateRedemptionStatus(reward.OwnerTwitchID, reward.RewardID, redemption.ID, false)
		if err != nil {
			log.Error(err.Error())
		}

		err = em.cpm.DeleteChannelPointReward(reward.OwnerTwitchID, reward.RewardID)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	count, err := em.db.CountNominations(context.Background(), reward.OwnerTwitchID, redemption.UserID)
	if err != nil {
		log.Errorf("failed to count nominations, refunding. %s", err.Error())
		err = em.helixclient.UpdateRedemptionStatus(reward.OwnerTwitchID, reward.RewardID, redemption.ID, false)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	if count >= election.MaxNominationPerUser {
		log.Infof("Max nominations %d reached, refunding", election.MaxNominationPerUser)
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

	isBlocked := em.db.IsEmoteBlocked(redemption.BroadcasterUserID, emoteID, dto.REWARD_SEVENTV)
	if isBlocked {
		log.Infof("Emote is blocked, refunding")
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

	err = em.helixclient.UpdateRedemptionStatus(reward.OwnerTwitchID, reward.RewardID, redemption.ID, true)
	if err != nil {
		log.Error(err.Error())
	}
}
