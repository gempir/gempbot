package channelpoint

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	nickHelix "github.com/nicklaw5/helix/v2"
)

var bttvRegex = regexp.MustCompile(`https?:\/\/betterttv.com\/emotes\/(\w*)`)
var sevenTvRegex = regexp.MustCompile(`https?:\/\/7tv.app\/emotes\/(\w*)`)

func (cpm *ChannelPointManager) HandleBttvRedemption(reward store.ChannelPointReward, redemption nickHelix.EventSubChannelPointsCustomRewardRedemptionEvent) {
	opts := UnmarshallBttvAdditionalOptions(reward.AdditionalOptions)
	success := false

	matches := bttvRegex.FindAllStringSubmatch(redemption.UserInput, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		emoteAdded, emoteRemoved, err := cpm.emotechief.SetBttvEmote(redemption.BroadcasterUserID, matches[0][1], redemption.BroadcasterUserLogin, opts.Slots)
		if err != nil {
			log.Warnf("Bttv error %s %s", redemption.BroadcasterUserLogin, err)
			cpm.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add emote from: @%s error: %s", redemption.UserName, err.Error()))
		} else if emoteAdded != nil && emoteRemoved != nil {
			success = true
			cpm.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new emote: %s redeemed by @%s removed: %s", emoteAdded.Code, redemption.UserName, emoteRemoved.Code))
		} else if emoteAdded != nil {
			success = true
			cpm.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new emote: %s redeemed by @%s", emoteAdded.Code, redemption.UserName))
		} else {
			success = true
			cpm.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new emote: [unknown] redeemed by @%s", redemption.UserName))
		}
	} else {
		cpm.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add emote from @%s error: no bttv link found in message", redemption.UserName))
	}

	token, err := cpm.db.GetUserAccessToken(redemption.BroadcasterUserID)
	if err != nil {
		log.Errorf("Failed to get userAccess token to update redemption status for %s", redemption.BroadcasterUserID)
		return
	} else {
		err := cpm.helixClient.UpdateRedemptionStatus(redemption.BroadcasterUserID, token.AccessToken, redemption.Reward.ID, redemption.ID, success)
		if err != nil {
			log.Errorf("Failed to update redemption status %s", err.Error())
			return
		}
	}
}

func UnmarshallBttvAdditionalOptions(jsonString string) BttvAdditionalOptions {
	if jsonString == "{}" {
		return BttvAdditionalOptions{Slots: 1}
	}

	var additionalOptions BttvAdditionalOptions

	if err := json.Unmarshal([]byte(jsonString), &additionalOptions); err != nil {
		log.Error(err)
		return BttvAdditionalOptions{Slots: 1}
	}

	return additionalOptions
}
