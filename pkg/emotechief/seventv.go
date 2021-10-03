package emotechief

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/gempir/gempbot/pkg/channelpoint"
	"github.com/gempir/gempbot/pkg/dto"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/utils"
	nickHelix "github.com/nicklaw5/helix/v2"
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

func (ec *EmoteChief) SetSevenTvEmote(channelUserID, login, emoteId, channel string, slots int) (addedEmote *sevenTvEmote, removedEmote *sevenTvEmote, err error) {
	newEmote, err := getSevenTvEmote(emoteId)
	if err != nil {
		return
	}

	if utils.BitField.HasBits(int64(newEmote.Visibility), int64(EmoteVisibilityPrivate)) ||
		utils.BitField.HasBits(int64(newEmote.Visibility), int64(EmoteVisibilityUnlisted)) {
		err = fmt.Errorf("7tv emote %s has incorrect visibility", newEmote.Name)
		return
	}

	var userData SevenTvUserResponse
	err = ec.QuerySevenTvGQL(SEVEN_TV_USER_DATA_QUERY, map[string]interface{}{"id": login}, &userData)
	if err != nil {
		return
	}

	emotes := userData.Data.User.Emotes
	emotesLimit := userData.Data.User.EmoteSlots
	for _, emote := range emotes {
		if emote.Name == newEmote.Name {
			return nil, nil, fmt.Errorf("Emote code already added")
		}
	}
	log.Infof("Current 7tv emotes: %d/%d", len(userData.Data.User.Emotes), userData.Data.User.EmoteSlots)

	var removalTargetEmoteId string

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

	emoteAddType := dto.EMOTE_ADD_REMOVED_PREVIOUS

	if len(confirmedEmotesAdded) == slots {
		removalTargetEmoteId = confirmedEmotesAdded[len(confirmedEmotesAdded)-1].EmoteID
		log.Infof("Found removal target %s in %s", removalTargetEmoteId, channelUserID)
	} else if len(emotes) >= emotesLimit {
		emoteAddType = dto.EMOTE_ADD_REMOVED_RANDOM
		log.Infof("Didn't find previous emote history of %d emotes and limit reached, choosing random in %s", slots, channelUserID)
		removalTargetEmoteId = emotes[rand.Intn(len(emotes))].ID
	}

	// do we need to remove the emote?
	if removalTargetEmoteId != "" {
		var empty struct{}
		err := ec.QuerySevenTvGQL(SEVEN_TV_DELETE_EMOTE_QUERY, map[string]interface{}{"ch": userData.Data.User.ID, "re": "redemption", "em": removalTargetEmoteId}, &empty)
		if err != nil {
			return nil, nil, err
		}

		ec.db.CreateEmoteAdd(channelUserID, dto.REWARD_SEVENTV, removalTargetEmoteId, emoteAddType)
	}

	// add the emote
	var empty struct{}
	err = ec.QuerySevenTvGQL(SEVEN_TV_ADD_EMOTE_QUERY, map[string]interface{}{"ch": userData.Data.User.ID, "re": "redemption", "em": newEmote.ID}, &empty)
	if err != nil {
		return
	}

	ec.db.CreateEmoteAdd(channelUserID, dto.REWARD_SEVENTV, emoteId, dto.EMOTE_ADD_ADD)

	if removalTargetEmoteId != "" {
		removedEmote, err = getSevenTvEmote(removalTargetEmoteId)
		if err != nil {
			return
		}
	}

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

func (ec *EmoteChief) HandleSeventvRedemption(reward store.ChannelPointReward, redemption nickHelix.EventSubChannelPointsCustomRewardRedemptionEvent) {
	opts := channelpoint.UnmarshallSevenTvAdditionalOptions(reward.AdditionalOptions)
	success := false

	matches := sevenTvRegex.FindAllStringSubmatch(redemption.UserInput, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		emoteAdded, emoteRemoved, err := ec.SetSevenTvEmote(redemption.BroadcasterUserID, redemption.BroadcasterUserLogin, matches[0][1], redemption.BroadcasterUserLogin, opts.Slots)
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

	token, err := ec.db.GetUserAccessToken(redemption.BroadcasterUserID)
	if err != nil {
		log.Errorf("Failed to get userAccess token to update redemption status for %s", redemption.BroadcasterUserID)
		return
	} else {
		err := ec.helixClient.UpdateRedemptionStatus(redemption.BroadcasterUserID, token.AccessToken, redemption.Reward.ID, redemption.ID, success)
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
