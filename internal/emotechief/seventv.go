package emotechief

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"

	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/nicklaw5/helix/v2"
)

var sevenTvRegex = regexp.MustCompile(`https?:\/\/(?:next\.)?7tv.app\/emotes\/(\w*)`)

func (ec *EmoteChief) VerifySetSevenTvEmote(channelUserID, emoteId, channel, redeemedByUsername string, slots int) (emoteAddType dto.EmoteChangeType, removalTargetEmoteId string, err error) {
	if ec.db.IsEmoteBlocked(channelUserID, emoteId, dto.REWARD_SEVENTV) {
		return dto.EMOTE_ADD_ADD, "", errors.New("emote is blocked")
	}

	nextEmote, err := ec.sevenTvClient.GetEmote(emoteId)
	if err != nil {
		return
	}

	user, err := ec.sevenTvClient.GetUser(channelUserID)
	if err != nil {
		return
	}

	for _, emote := range user.Emotes {
		if emote.Code == nextEmote.Code {
			return dto.EMOTE_ADD_ADD, "", fmt.Errorf("emote code \"%s\" already added", nextEmote.Code)
		}
	}
	log.Infof("Current 7TV emotes: %d/%d", len(user.Emotes), user.EmoteSlots)

	emotesAdded := ec.db.GetEmoteAdded(channelUserID, dto.REWARD_SEVENTV, slots)
	log.Infof("Total Previous emotes %d in %s", len(emotesAdded), channelUserID)

	if len(emotesAdded) > 0 {
		oldestEmote := emotesAdded[len(emotesAdded)-1]
		if !oldestEmote.Blocked {
			for _, sharedEmote := range user.Emotes {
				if oldestEmote.EmoteID == sharedEmote.ID {
					removalTargetEmoteId = oldestEmote.EmoteID
					log.Infof("Found removal target %s in %s", removalTargetEmoteId, channelUserID)
				}
			}
		} else {
			log.Infof("Removal target %s is already blocked, so already removed, skipping removal", oldestEmote.EmoteID)
		}
	}

	emoteAddType = dto.EMOTE_ADD_REMOVED_PREVIOUS
	if removalTargetEmoteId == "" && len(user.Emotes) >= user.EmoteSlots {
		if len(user.Emotes) == 0 {
			return dto.EMOTE_ADD_ADD, "", errors.New("emotes limit reached and can't find amount of emotes added to choose random")
		}

		emoteAddType = dto.EMOTE_ADD_REMOVED_RANDOM
		log.Infof("Didn't find previous emote history of %d emotes and limit reached, choosing random in %s", slots, channelUserID)
		removalTargetEmoteId = user.Emotes[rand.Intn(len(user.Emotes))].ID
	}

	return
}

func (ec *EmoteChief) setSevenTvEmote(channelUserID, emoteId, channel, redeemedByUsername string, slots int) (addedEmoteId string, removedEmoteID string, err error) {
	emoteAddType, removalTargetEmoteId, err := ec.VerifySetSevenTvEmote(channelUserID, emoteId, channel, redeemedByUsername, slots)
	if err != nil {
		return "", "", err
	}

	// do we need to remove the emote?
	if removalTargetEmoteId != "" {
		err := ec.sevenTvClient.RemoveEmote(channelUserID, removalTargetEmoteId)
		if err != nil {
			return "", "", err
		}

		ec.db.CreateEmoteAdd(channelUserID, dto.REWARD_SEVENTV, removalTargetEmoteId, emoteAddType)
	}

	err = ec.sevenTvClient.AddEmote(channelUserID, emoteId)
	if err != nil {
		return "", removalTargetEmoteId, err
	}

	ec.db.CreateEmoteAdd(channelUserID, dto.REWARD_SEVENTV, emoteId, dto.EMOTE_ADD_ADD)

	return emoteId, removalTargetEmoteId, nil
}

func GetSevenTvEmoteId(message string) (string, error) {
	matches := sevenTvRegex.FindAllStringSubmatch(message, -1)

	if len(matches) == 1 && len(matches[0]) == 2 {
		return matches[0][1], nil
	}

	return "", errors.New("no 7TV emote link found")
}

func (ec *EmoteChief) VerifySeventvRedemption(reward store.ChannelPointReward, redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent) bool {
	opts := channelpoint.UnmarshallSevenTvAdditionalOptions(reward.AdditionalOptions)

	emoteID, err := GetSevenTvEmoteId(redemption.UserInput)
	if err == nil {
		_, _, err := ec.VerifySetSevenTvEmote(redemption.BroadcasterUserID, emoteID, redemption.BroadcasterUserLogin, redemption.UserLogin, opts.Slots)
		if err != nil {
			log.Warnf("7TV error %s %s", redemption.BroadcasterUserLogin, err)
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add 7TV emote from @%s error: %s", redemption.UserName, err.Error()))
			return false
		}

		return true
	}

	ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add 7TV emote from @%s error: %s", redemption.UserName, err.Error()))
	return false
}

func (ec *EmoteChief) HandleSeventvRedemption(reward store.ChannelPointReward, redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent, updateStatus bool) {
	opts := channelpoint.UnmarshallSevenTvAdditionalOptions(reward.AdditionalOptions)
	success := false

	emoteID, err := GetSevenTvEmoteId(redemption.UserInput)
	if err == nil {
		log.Infof("Seen 7TV emote link %s", emoteID)
		added, removed, settingErr := ec.setSevenTvEmote(redemption.BroadcasterUserID, emoteID, redemption.BroadcasterUserLogin, redemption.UserName, opts.Slots)
		addedEmote, err := ec.sevenTvClient.GetEmote(added)
		if err != nil && len(added) > 0 {
			log.Error("Error fetching added emote: " + err.Error())
		}
		removedEmote, err := ec.sevenTvClient.GetEmote(removed)
		if err != nil && len(removed) > 0 {
			log.Error("Error fetching removed emote: " + err.Error())
		}

		if settingErr != nil {
			log.Warnf("7TV error %s %s", redemption.BroadcasterUserLogin, settingErr)
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add 7TV emote from @%s %s", redemption.UserName, settingErr.Error()))
		} else if addedEmote.Code != "" && removedEmote.Code != "" {
			success = true
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new 7TV emote %s redeemed by @%s removed %s", addedEmote.Code, redemption.UserName, removedEmote.Code))
		} else if addedEmote.Code != "" {
			success = true
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new 7TV emote %s redeemed by @%s", addedEmote.Code, redemption.UserName))
		} else {
			success = true
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new 7TV emote [unknown] redeemed by @%s", redemption.UserName))
		}
	} else {
		ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add 7TV emote from @%s %s", redemption.UserName, err.Error()))
	}

	if redemption.UserID == dto.GEMPIR_USER_ID {
		return
	}

	if updateStatus {
		err := ec.helixClient.UpdateRedemptionStatus(redemption.BroadcasterUserID, redemption.Reward.ID, redemption.ID, success)
		if err != nil {
			log.Errorf("Failed to update redemption status %s", err.Error())
			return
		}
	}
}
