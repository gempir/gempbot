package emotechief

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/gempir/gempbot/pkg/channelpoint"
	"github.com/gempir/gempbot/pkg/dto"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/utils"
	helix "github.com/nicklaw5/helix/v2"
)

var sevenTvRegex = regexp.MustCompile(`https?:\/\/7tv.app\/emotes\/(\w*)`)

const sevenTvApiBaseUrl = "https://api.7tv.app/v2"

const (
	EmoteVisibilityPrivate int32 = 1 << iota
	EmoteVisibilityGlobal
	EmoteVisibilityUnlisted
	EmoteVisibilityOverrideBTTV
	EmoteVisibilityOverrideFFZ
	EmoteVisibilityOverrideTwitchGlobal
	EmoteVisibilityOverrideTwitchSubscriber
	EmoteVisibilityZeroWidth
	EmoteVisibilityPermanentlyUnlisted

	EmoteVisibilityAll int32 = (1 << iota) - 1
)

type SevenTvUserResponse struct {
	Data struct {
		User struct {
			ID     string `json:"id"`
			Emotes []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Status     int    `json:"status"`
				Visibility int    `json:"visibility"`
				Width      []int  `json:"width"`
				Height     []int  `json:"height"`
			} `json:"emotes"`
			EmoteSlots int `json:"emote_slots"`
		} `json:"user"`
	} `json:"data"`
}

func (ec *EmoteChief) VerifySetSevenTvEmote(channelUserID, emoteId, channel, redeemedByUsername string, slots int) (newEmote *sevenTvEmote, emoteAddType dto.EmoteChangeType, userData *SevenTvUserResponse, removalTargetEmoteId string, err error) {
	if ec.db.IsEmoteBlocked(channelUserID, emoteId, dto.REWARD_SEVENTV) {
		return nil, dto.EMOTE_ADD_ADD, nil, "", errors.New("Emote is blocked")
	}

	newEmote, err = getSevenTvEmote(emoteId)
	if err != nil {
		return
	}

	if utils.BitField.HasBits(int64(newEmote.Visibility), int64(EmoteVisibilityPrivate)) ||
		utils.BitField.HasBits(int64(newEmote.Visibility), int64(EmoteVisibilityUnlisted)) {
		err = fmt.Errorf("7tv emote %s has incorrect visibility", newEmote.Name)
		return
	}

	err = ec.QuerySevenTvGQL(SEVEN_TV_USER_DATA_QUERY, map[string]interface{}{"id": channel}, &userData)
	if err != nil {
		return
	}

	emotes := userData.Data.User.Emotes
	emotesLimit := userData.Data.User.EmoteSlots
	for _, emote := range emotes {
		if emote.Name == newEmote.Name {
			return nil, dto.EMOTE_ADD_ADD, nil, "", fmt.Errorf("Emote code \"%s\" already added", newEmote.Name)
		}
	}
	log.Infof("Current 7tv emotes: %d/%d", len(userData.Data.User.Emotes), userData.Data.User.EmoteSlots)

	emotesAdded := ec.db.GetEmoteAdded(channelUserID, dto.REWARD_SEVENTV, slots)
	log.Infof("Total Previous emotes %d in %s", len(emotesAdded), channelUserID)

	confirmedEmotesAdded := []store.EmoteAdd{}
	for _, emote := range emotesAdded {
		for _, sharedEmote := range emotes {
			if emote.EmoteID == sharedEmote.ID {
				confirmedEmotesAdded = append(confirmedEmotesAdded, emote)
			}
		}
	}

	emoteAddType = dto.EMOTE_ADD_REMOVED_PREVIOUS

	if len(confirmedEmotesAdded) == slots {
		removalTargetEmoteId = confirmedEmotesAdded[len(confirmedEmotesAdded)-1].EmoteID
		log.Infof("Found removal target %s in %s", removalTargetEmoteId, channelUserID)
	} else if len(emotes) >= emotesLimit {
		log.Infof("7tv Userdata for %s %v", channelUserID, userData)
		if len(emotes) == 0 {
			return nil, dto.EMOTE_ADD_ADD, nil, "", errors.New("emotes limit reached and can't find amount of emotes added to choose random")
		}

		emoteAddType = dto.EMOTE_ADD_REMOVED_RANDOM
		log.Infof("Didn't find previous emote history of %d emotes and limit reached, choosing random in %s", slots, channelUserID)
		removalTargetEmoteId = emotes[rand.Intn(len(emotes))].ID
	}

	return
}

func (ec *EmoteChief) SetSevenTvEmote(channelUserID, emoteId, channel, redeemedByUsername string, slots int) (addedEmote *sevenTvEmote, removedEmote *sevenTvEmote, err error) {
	newEmote, emoteAddType, userData, removalTargetEmoteId, err := ec.VerifySetSevenTvEmote(channelUserID, emoteId, channel, redeemedByUsername, slots)
	if err != nil {
		return nil, nil, err
	}

	// do we need to remove the emote?
	if removalTargetEmoteId != "" {
		var empty struct{}
		err := ec.QuerySevenTvGQL(
			SEVEN_TV_DELETE_EMOTE_QUERY,
			map[string]interface{}{
				"ch": userData.Data.User.ID,
				"re": fmt.Sprintf("removed for redemption by %s, new emote: %s", redeemedByUsername, newEmote.Name),
				"em": removalTargetEmoteId,
			}, &empty,
		)
		if err != nil {
			return nil, nil, err
		}

		ec.db.CreateEmoteAdd(channelUserID, dto.REWARD_SEVENTV, removalTargetEmoteId, emoteAddType)
	}

	removedEmoteName := removalTargetEmoteId
	if removalTargetEmoteId != "" {
		var sevenTvErr error
		removedEmote, sevenTvErr = getSevenTvEmote(removalTargetEmoteId)
		if sevenTvErr == nil {
			removedEmoteName = removedEmote.Name
		}
		if sevenTvErr != nil {
			log.Errorf("Failed to fetch removed Emote %s", sevenTvErr.Error())
		}
	}

	// add the emote
	var empty struct{}
	err = ec.QuerySevenTvGQL(
		SEVEN_TV_ADD_EMOTE_QUERY,
		map[string]interface{}{
			"ch": userData.Data.User.ID,
			"re": fmt.Sprintf("redemption by %s, replaced: %s", redeemedByUsername, removedEmoteName),
			"em": newEmote.ID,
		}, &empty,
	)
	if err != nil {
		return
	}

	ec.db.CreateEmoteAdd(channelUserID, dto.REWARD_SEVENTV, emoteId, dto.EMOTE_ADD_ADD)

	return newEmote, removedEmote, nil
}

const SEVEN_TV_ADD_EMOTE_QUERY = `mutation AddChannelEmote($ch: String!, $em: String!, $re: String!) {addChannelEmote(channel_id: $ch, emote_id: $em, reason: $re) {emote_ids}}`
const SEVEN_TV_DELETE_EMOTE_QUERY = `mutation RemoveChannelEmote($ch: String!, $em: String!, $re: String!) {removeChannelEmote(channel_id: $ch, emote_id: $em, reason: $re) {emote_ids}}`
const SEVEN_TV_USER_DATA_QUERY = `
query GetUser($id: String!) {
	user(id: $id) {
	  ...FullUser
	}
  }
  
fragment FullUser on User {
	id
	emotes {
		id
		name
		status
		visibility
		width
		height
	}
	emote_slots
}
`

type GqlQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func (ec *EmoteChief) QuerySevenTvGQL(query string, variables map[string]interface{}, response interface{}) error {
	data, err := json.Marshal(GqlQuery{Query: query, Variables: variables})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.7tv.app/v2/gql", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Bearer "+ec.cfg.SevenTvToken)
	if err != nil {
		return err
	}

	resp, err := ec.httpClient.Do(req)
	if err != nil {
		return err
	}
	log.Infof("%d 7tv query '%s' with '%v'", resp.StatusCode, query, variables)

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("Error %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (ec *EmoteChief) VerifySeventvRedemption(reward store.ChannelPointReward, redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent) bool {
	opts := channelpoint.UnmarshallSevenTvAdditionalOptions(reward.AdditionalOptions)

	matches := sevenTvRegex.FindAllStringSubmatch(redemption.UserInput, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		_, _, _, _, err := ec.VerifySetSevenTvEmote(redemption.BroadcasterUserID, matches[0][1], redemption.UserLogin, redemption.BroadcasterUserLogin, opts.Slots)
		if err != nil {
			log.Warnf("7tv error %s %s", redemption.BroadcasterUserLogin, err)
			ec.chatClient.WaitForConnect()
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add 7tv emote from: @%s error: %s", redemption.UserName, err.Error()))
			return false
		}

		return true
	}

	ec.chatClient.WaitForConnect()
	ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add 7tv emote from @%s error: no 7tv link found in message", redemption.UserName))
	return false
}

func (ec *EmoteChief) HandleSeventvRedemption(reward store.ChannelPointReward, redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent, updateStatus bool) {
	opts := channelpoint.UnmarshallSevenTvAdditionalOptions(reward.AdditionalOptions)
	success := false

	matches := sevenTvRegex.FindAllStringSubmatch(redemption.UserInput, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		emoteAdded, emoteRemoved, err := ec.SetSevenTvEmote(redemption.BroadcasterUserID, matches[0][1], redemption.BroadcasterUserLogin, redemption.UserName, opts.Slots)
		if err != nil {
			log.Warnf("7tv error %s %s", redemption.BroadcasterUserLogin, err)
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add 7tv emote from: @%s error: %s", redemption.UserName, err.Error()))
		} else if emoteAdded != nil && emoteRemoved != nil {
			success = true
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new 7tv emote: %s redeemed by @%s removed: %s", emoteAdded.Name, redemption.UserName, emoteRemoved.Name))
		} else if emoteAdded != nil {
			success = true
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new 7tv emote: %s redeemed by @%s", emoteAdded.Name, redemption.UserName))
		} else {
			success = true
			ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("✅ Added new 7tv emote: [unknown] redeemed by @%s", redemption.UserName))
		}
	} else {
		ec.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Failed to add 7tv emote from @%s error: no 7tv link found in message", redemption.UserName))
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

type sevenTvEmote struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner struct {
		ID          string `json:"id"`
		TwitchID    string `json:"twitch_id"`
		Login       string `json:"login"`
		DisplayName string `json:"display_name"`
		Role        struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Position int    `json:"position"`
			Color    int    `json:"color"`
			Allowed  int    `json:"allowed"`
			Denied   int    `json:"denied"`
			Default  bool   `json:"default"`
		} `json:"role"`
	} `json:"owner"`
	Visibility       int           `json:"visibility"`
	VisibilitySimple []interface{} `json:"visibility_simple"`
	Mime             string        `json:"mime"`
	Status           int           `json:"status"`
	Tags             []interface{} `json:"tags"`
	Width            []int         `json:"width"`
	Height           []int         `json:"height"`
	Urls             [][]string    `json:"urls"`
}

func getSevenTvEmote(emoteID string) (*sevenTvEmote, error) {
	if emoteID == "" {
		return nil, nil
	}

	response, err := http.Get(sevenTvApiBaseUrl + "/emotes/" + emoteID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if response.StatusCode <= 100 || response.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad 7tv response: %d", response.StatusCode)
	}

	var emoteResponse sevenTvEmote
	err = json.NewDecoder(response.Body).Decode(&emoteResponse)
	if err != nil {
		return nil, err
	}

	return &emoteResponse, nil
}
