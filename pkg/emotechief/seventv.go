package emotechief

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gempir/gempbot/pkg/channelpoint"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	nickHelix "github.com/nicklaw5/helix/v2"
)

var sevenTvRegex = regexp.MustCompile(`https?:\/\/7tv.app\/emotes\/(\w*)`)

const sevenTvApiBaseUrl = "https://api.7tv.app/v2"

type GqlQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type SevenTvUserResponse struct {
	Data struct {
		User struct {
			ID           string        `json:"id"`
			EmoteAliases []interface{} `json:"emote_aliases"`
			Emotes       []struct {
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

func (e *EmoteChief) SetSevenTvEmote(channelUserID, login, emoteId, channel string, slots int) (addedEmote *sevenTvEmote, removedEmote *sevenTvEmote, err error) {
	emote, err := getSevenTvEmote(emoteId)
	if err != nil {
		return
	}

	gqlQuery := GqlQuery{
		Query: `
		query GetUser($id: String!) {
			user(id: $id) {
			  ...FullUser
			}
		  }
		  
		fragment FullUser on User {
			id
			emote_aliases
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
		`,
		Variables: map[string]interface{}{"id": login},
	}

	data, err := json.Marshal(gqlQuery)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.7tv.app/v2/gql", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("authorization", "Bearer "+e.cfg.SevenTvToken)
	if err != nil {
		return
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return
	}

	var userData SevenTvUserResponse
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return
	}

	log.Infof("%d/%d", len(userData.Data.User.Emotes), userData.Data.User.EmoteSlots)

	return emote, nil, nil
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
