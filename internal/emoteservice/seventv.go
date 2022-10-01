package emoteservice

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/gempir/gempbot/internal/utils"
)

const DefaultSevenTvApiBaseUrl = "https://api.7tv.app/v2"
const DefaultSevenTvV3ApiBaseUrl = "https://7tv.io/v3"
const DefaultSevenTvGqlBaseUrl = "https://api.7tv.app/v2/gql"

type SevenTvClient struct {
	store      store.Store
	apiBaseUrl string
	gqlBaseUrl string
}

func NewSevenTvClient(store store.Store) *SevenTvClient {
	return &SevenTvClient{
		store:      store,
		apiBaseUrl: DefaultSevenTvApiBaseUrl,
		gqlBaseUrl: DefaultSevenTvGqlBaseUrl,
	}
}

type UserV3 struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	ProfilePictureURL string `json:"profile_picture_url"`
	DisplayName       string `json:"display_name"`
	Style             struct {
		Color int         `json:"color"`
		Paint interface{} `json:"paint"`
	} `json:"style"`
	Biography string `json:"biography"`
	Editors   []struct {
		ID          string `json:"id"`
		Permissions int    `json:"permissions"`
		Visible     bool   `json:"visible"`
		AddedAt     int64  `json:"added_at"`
	} `json:"editors"`
	Roles       []string `json:"roles"`
	Connections []struct {
		ID            string `json:"id"`
		Platform      string `json:"platform"`
		Username      string `json:"username"`
		DisplayName   string `json:"display_name"`
		LinkedAt      int64  `json:"linked_at"`
		EmoteCapacity int    `json:"emote_capacity"`
	} `json:"connections"`
}

func (c *SevenTvClient) GetUserV3(userID string) (UserV3, error) {
	// first figure out the bttvUserId for the channel, might cache this later on
	var userResp UserV3
	err := requests.
		URL(DefaultSevenTvV3ApiBaseUrl).
		Pathf("/users/%s", userID).
		ToJSON(&userResp).
		Fetch(context.Background())
	if err != nil {
		return UserV3{}, err
	}

	return userResp, nil
}

func (c *SevenTvClient) GetTwitchConnection(userID string) (string, error) {
	user, err := c.GetUserV3(userID)
	if err != nil {
		return "", err
	}

	for _, connection := range user.Connections {
		if connection.Platform == "TWITCH" {
			return connection.Username, nil
		}
	}

	return "", fmt.Errorf("no twitch connection found for user %s", userID)
}

func (c *SevenTvClient) GetEmote(emoteID string) (Emote, error) {
	var emoteData sevenTvEmote

	err := requests.URL(c.apiBaseUrl + "/emotes/" + emoteID).
		ToJSON(&emoteData).
		Fetch(context.Background())

	if utils.BitField.HasBits(int64(emoteData.Visibility), int64(EmoteVisibilityPrivate)) ||
		utils.BitField.HasBits(int64(emoteData.Visibility), int64(EmoteVisibilityUnlisted)) {

		return Emote{}, fmt.Errorf("emote %s has incorrect visibility", emoteData.Name)
	}

	return Emote{Code: emoteData.Name, ID: emoteData.ID}, err
}

func (c *SevenTvClient) RemoveEmote(channelUserID, emoteID string) error {
	user, err := c.GetUser(channelUserID)
	if err != nil {
		return err
	}

	var empty struct{}
	err = c.QuerySevenTvGQL(
		`mutation RemoveChannelEmote($ch: String!, $em: String!, $re: String!) {removeChannelEmote(channel_id: $ch, emote_id: $em, reason: $re) {emote_ids}}`,
		map[string]interface{}{
			"ch": user.ID,
			"re": "blocked emote",
			"em": emoteID,
		}, &empty,
	)

	return err
}

func (c *SevenTvClient) AddEmote(channelUserID string, emoteID string) error {
	user, err := c.GetUser(channelUserID)
	if err != nil {
		return err
	}

	var empty struct{}
	err = c.QuerySevenTvGQL(
		`mutation AddChannelEmote($ch: String!, $em: String!, $re: String!) {addChannelEmote(channel_id: $ch, emote_id: $em, reason: $re) {emote_ids}}`,
		map[string]interface{}{
			"ch": user.ID,
			"re": "bot.gempir.com redemption",
			"em": emoteID,
		}, &empty,
	)

	return err
}

func (c *SevenTvClient) GetUser(channelID string) (User, error) {
	var userData SevenTvUserResponse
	err := c.QuerySevenTvGQL(`
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
	`, map[string]interface{}{"id": channelID}, &userData)
	if err != nil {
		return User{}, err
	}

	var emotes []Emote
	for _, emote := range userData.Data.User.Emotes {
		emotes = append(emotes, Emote{ID: emote.ID, Code: emote.Name})
	}

	return User{ID: userData.Data.User.ID, Emotes: emotes, EmoteSlots: userData.Data.User.EmoteSlots}, nil
}

func (c *SevenTvClient) QuerySevenTvGQL(query string, variables map[string]interface{}, response interface{}) error {
	gqlQuery := gqlQuery{Query: query, Variables: variables}

	err := requests.
		URL(c.gqlBaseUrl).
		BodyJSON(gqlQuery).
		Bearer(c.store.GetSevenTvToken(context.Background())).
		ToJSON(&response).
		Fetch(context.Background())
	if err != nil {
		log.Infof("7tv query '%s' with '%v' resp: '%v'", query, variables, response)
		return err
	}

	log.Infof("7tv query '%s' with '%v' resp: '%v'", query, variables, response)

	return nil
}

const SEVEN_TV_ADD_EMOTE_QUERY = `mutation AddChannelEmote($ch: String!, $em: String!, $re: String!) {addChannelEmote(channel_id: $ch, emote_id: $em, reason: $re) {emote_ids}}`
