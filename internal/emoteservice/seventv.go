package emoteservice

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/utils"
)

const DefaultSevenTvApiBaseUrl = "https://api.7tv.app/v2"
const DefaultSevenTvGqlBaseUrl = "https://api.7tv.app/v2/gql"

type SevenTvClient struct {
	apiToken   string
	apiBaseUrl string
	gqlBaseUrl string
}

func NewSevenTvClient(apiToken string) *SevenTvClient {
	return &SevenTvClient{
		apiToken:   apiToken,
		apiBaseUrl: DefaultSevenTvApiBaseUrl,
		gqlBaseUrl: DefaultSevenTvGqlBaseUrl,
	}
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
		Bearer(c.apiToken).
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
