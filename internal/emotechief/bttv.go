package emotechief

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/nicklaw5/helix/v2"
)

func (e *EmoteChief) VerifySetBttvEmote(channelUserID, emoteId, channel string, slots int) (addedEmote *bttvEmoteResponse, emoteAddType dto.EmoteChangeType, bttvUserId string, removalTargetEmoteId string, err error) {
	if e.db.IsEmoteBlocked(channelUserID, emoteId, dto.REWARD_BTTV) {
		return nil, dto.EMOTE_ADD_ADD, "", "", errors.New("Emote is blocked")
	}

	bttvToken := e.db.GetBttvToken(context.Background())

	addedEmote, err = getBttvEmote(emoteId)
	if err != nil {
		return
	}

	if !addedEmote.Sharing {
		err = errors.New("Emote is not shared")
		return
	}

	// first figure out the bttvUserId for the channel, might cache this later on
	var userResp bttvUserResponse
	err = requests.
		URL(BTTV_API).
		Pathf("/3/cached/users/twitch/%s", channelUserID).
		ToJSON(&userResp).
		Fetch(context.Background())
	if err != nil {
		return
	}

	bttvUserId = userResp.ID

	// figure out the limit for the channel, might also chache this later on with expiry
	var dashboards dashboardsResponse
	err = requests.
		URL(BTTV_API).
		Pathf("/3/account/dashboards").
		Bearer(bttvToken).
		ToJSON(&dashboards).
		Fetch(context.Background())
	if err != nil {
		return
	}

	var dbCfg dashboardCfg
	for _, db := range dashboards {
		if db.ID == bttvUserId {
			dbCfg = db
		}
	}
	if dbCfg.ID == "" {
		err = errors.New("No permission to moderate, add gempbot as BetterTTV editor")
		return
	}
	sharedEmotesLimit := dbCfg.Limits.Sharedemotes

	// figure currently added emotes
	var dashboard bttvDashboardResponse
	err = requests.
		URL(BTTV_API).
		Pathf("/3/users/%s", bttvUserId).
		Param("limited", "false").
		Param("personal", "false").
		ToJSON(&dashboard).
		Fetch(context.Background())
	if err != nil {
		return
	}

	for _, emote := range dashboard.Sharedemotes {
		if emote.ID == emoteId {
			err = errors.New("Emote already added")
			return
		}
		if emote.Code == addedEmote.Code {
			err = fmt.Errorf("Emote code \"%s\" already added", addedEmote.Code)
			return
		}
	}

	for _, emote := range dashboard.Channelemotes {
		if emote.ID == emoteId {
			err = fmt.Errorf("Emote \"%s\" already a channel emote", emote.Code)
			return
		}
		if emote.Code == addedEmote.Code {
			err = fmt.Errorf("Emote code \"%s\" already a channel emote", addedEmote.Code)
			return
		}
	}
	log.Infof("Current shared emotes: %d/%d", len(dashboard.Sharedemotes), sharedEmotesLimit)

	emotesAdded := e.db.GetEmoteAdded(channelUserID, dto.REWARD_BTTV, slots)
	log.Infof("Total Previous emotes %d in %s", len(emotesAdded), channelUserID)

	if len(emotesAdded) > 0 {
		oldestEmote := emotesAdded[len(emotesAdded)-1]
		if !oldestEmote.Blocked {
			for _, sharedEmote := range dashboard.Sharedemotes {
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
	if removalTargetEmoteId == "" && len(dashboard.Sharedemotes) >= sharedEmotesLimit {
		if len(dashboard.Sharedemotes) == 0 {
			return nil, dto.EMOTE_ADD_ADD, "", "", errors.New("emotes limit reached and can't find amount of emotes added to choose random")
		}

		emoteAddType = dto.EMOTE_ADD_REMOVED_RANDOM
		log.Infof("Didn't find previous emote history of %d emotes and limit reached, choosing random in %s", slots, channelUserID)
		removalTargetEmoteId = dashboard.Sharedemotes[rand.Intn(len(dashboard.Sharedemotes))].ID
	}

	return
}

func (e *EmoteChief) RemoveBttvEmote(channelUserID, emoteID string) (*bttvEmoteResponse, error) {
	bttvToken := e.db.GetBttvToken(context.Background())

	var userResp bttvUserResponse
	err := requests.
		URL(BTTV_API).
		Pathf("/3/cached/users/twitch/%s", channelUserID).
		ToJSON(&userResp).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	bttvUserId := userResp.ID

	err = requests.
		URL(BTTV_API).
		Pathf("/3/emotes/%s/shared/%s", emoteID, bttvUserId).
		Bearer(bttvToken).
		Method(http.MethodDelete).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	e.db.CreateEmoteAdd(channelUserID, dto.REWARD_BTTV, emoteID, dto.EMOTE_ADD_REMOVED_BLOCKED)
	log.Infof("Blocked channelId: %s emoteId: %s", channelUserID, emoteID)

	return getBttvEmote(emoteID)
}

func (e *EmoteChief) SetBttvEmote(channelUserID, emoteId, channel string, slots int) (addedEmote *bttvEmoteResponse, removedEmote *bttvEmoteResponse, err error) {
	addedEmote, emoteAddType, bttvUserId, removalTargetEmoteId, err := e.VerifySetBttvEmote(channelUserID, emoteId, channel, slots)
	if err != nil {
		return nil, nil, err
	}

	bttvToken := e.db.GetBttvToken(context.Background())

	// do we need to remove the emote?
	if removalTargetEmoteId != "" {
		err = requests.
			URL(BTTV_API).
			Pathf("/3/emotes/%s/shared/%s", removalTargetEmoteId, bttvUserId).
			Bearer(bttvToken).
			Method(http.MethodDelete).
			Fetch(context.Background())
		if err != nil {
			return
		}

		e.db.CreateEmoteAdd(channelUserID, dto.REWARD_BTTV, removalTargetEmoteId, emoteAddType)
		log.Infof("Deleted channelId: %s emoteId: %s", channelUserID, removalTargetEmoteId)

		removedEmote, _ = getBttvEmote(removalTargetEmoteId)
	}

	// Add new emote
	err = requests.
		URL(BTTV_API).
		Pathf("/3/emotes/%s/shared/%s", emoteId, bttvUserId).
		Bearer(bttvToken).
		Method(http.MethodPut).
		Fetch(context.Background())
	if err != nil {
		return
	}

	log.Infof("Added channelId: %s emoteId: %s", channelUserID, emoteId)
	e.db.CreateEmoteAdd(channelUserID, dto.REWARD_BTTV, emoteId, dto.EMOTE_ADD_ADD)

	return
}

func getBttvEmote(emoteID string) (*bttvEmoteResponse, error) {
	if emoteID == "" {
		return nil, nil
	}

	var emoteResp bttvEmoteResponse
	err := requests.
		URL(BTTV_API).
		Pathf("/3/emotes/%s", emoteID).
		ToJSON(&emoteResp).
		Fetch(context.Background())

	return &emoteResp, err
}

var bttvRegex = regexp.MustCompile(`https?:\/\/betterttv.com\/emotes\/(\w*)`)

func (ec *EmoteChief) VerifyBttvRedemption(reward store.ChannelPointReward, redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent) bool {
	opts := channelpoint.UnmarshallBttvAdditionalOptions(reward.AdditionalOptions)

	emoteID, err := GetBttvEmoteId(redemption.UserInput)
	if err == nil {
		_, _, _, _, err := ec.VerifySetBttvEmote(redemption.BroadcasterUserID, emoteID, redemption.BroadcasterUserLogin, opts.Slots)
		if err != nil {
			log.Warnf("Bttv error %s %s", redemption.BroadcasterUserLogin, err)
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add bttv emote from @%s error: %s", redemption.UserName, err.Error()))
			return false
		}

		return true
	}

	ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add bttv emote from @%s error: %s", redemption.UserName, err.Error()))
	return false
}

func GetBttvEmoteId(message string) (string, error) {
	matches := bttvRegex.FindAllStringSubmatch(message, -1)

	if len(matches) == 1 && len(matches[0]) == 2 {
		return matches[0][1], nil
	}

	return "", errors.New("no bttv emote link found")
}

func (ec *EmoteChief) HandleBttvRedemption(reward store.ChannelPointReward, redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent, updateStatus bool) {
	opts := channelpoint.UnmarshallBttvAdditionalOptions(reward.AdditionalOptions)
	success := false

	emoteID, err := GetBttvEmoteId(redemption.UserInput)
	if err == nil {
		emoteAdded, emoteRemoved, err := ec.SetBttvEmote(redemption.BroadcasterUserID, emoteID, redemption.BroadcasterUserLogin, opts.Slots)
		if err != nil {
			log.Warnf("Bttv error %s %s", redemption.BroadcasterUserLogin, err)
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add bttv emote from @%s error: %s", redemption.UserName, err.Error()))
		} else if emoteAdded != nil && emoteRemoved != nil {
			success = true
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new bttv emote %s redeemed by @%s removed: %s", emoteAdded.Code, redemption.UserName, emoteRemoved.Code))
		} else if emoteAdded != nil {
			success = true
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new bttv emote %s redeemed by @%s", emoteAdded.Code, redemption.UserName))
		} else {
			success = true
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new bttv emote [unknown] redeemed by @%s", redemption.UserName))
		}
	} else {
		ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add bttv emote from @%s error: %s", redemption.UserName, err.Error()))
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

const BTTV_API = "https://api.betterttv.net"

type bttvUserResponse struct {
	ID string `json:"id"`
}

type bttvDashboardResponse struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Displayname   string   `json:"displayName"`
	Providerid    string   `json:"providerId"`
	Bots          []string `json:"bots"`
	Channelemotes []struct {
		ID             string    `json:"id"`
		Code           string    `json:"code"`
		Imagetype      string    `json:"imageType"`
		Userid         string    `json:"userId"`
		Createdat      time.Time `json:"createdAt"`
		Updatedat      time.Time `json:"updatedAt"`
		Global         bool      `json:"global"`
		Live           bool      `json:"live"`
		Sharing        bool      `json:"sharing"`
		Approvalstatus string    `json:"approvalStatus"`
	} `json:"liveEmotes"`
	Sharedemotes []struct {
		ID             string    `json:"id"`
		Code           string    `json:"code"`
		Imagetype      string    `json:"imageType"`
		Createdat      time.Time `json:"createdAt"`
		Updatedat      time.Time `json:"updatedAt"`
		Global         bool      `json:"global"`
		Live           bool      `json:"live"`
		Sharing        bool      `json:"sharing"`
		Approvalstatus string    `json:"approvalStatus"`
		User           struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Displayname string `json:"displayName"`
			Providerid  string `json:"providerId"`
		} `json:"user"`
	} `json:"sharedEmotes"`
}

type bttvEmoteResponse struct {
	ID             string    `json:"id"`
	Code           string    `json:"code"`
	Imagetype      string    `json:"imageType"`
	Createdat      time.Time `json:"createdAt"`
	Updatedat      time.Time `json:"updatedAt"`
	Global         bool      `json:"global"`
	Live           bool      `json:"live"`
	Sharing        bool      `json:"sharing"`
	Approvalstatus string    `json:"approvalStatus"`
	User           struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Displayname string `json:"displayName"`
		Providerid  string `json:"providerId"`
	} `json:"user"`
}

type dashboardsResponse []dashboardCfg

type dashboardCfg struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Displayname string `json:"displayName"`
	Providerid  string `json:"providerId"`
	Avatar      string `json:"avatar"`
	Limits      struct {
		Channelemotes  int `json:"channelEmotes"`
		Sharedemotes   int `json:"sharedEmotes"`
		Personalemotes int `json:"personalEmotes"`
	} `json:"limits"`
}
