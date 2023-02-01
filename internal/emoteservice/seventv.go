package emoteservice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

const DefaultSevenTvApiBaseUrl = "https://api.7tv.app/v2"
const DefaultSevenTvV3ApiBaseUrl = "https://7tv.io/v3"
const DefaultSevenTvGqlBaseUrl = "https://api.7tv.app/v2/gql"
const DefaultSevenTvGqlV3BaseUrl = "https://api.7tv.app/v3/gql"

type SevenTvClient struct {
	store        store.Store
	apiBaseUrl   string
	gqlBaseUrl   string
	gqlV3BaseUrl string
}

func NewSevenTvClient(store store.Store) *SevenTvClient {
	return &SevenTvClient{
		store:        store,
		apiBaseUrl:   DefaultSevenTvApiBaseUrl,
		gqlBaseUrl:   DefaultSevenTvGqlBaseUrl,
		gqlV3BaseUrl: DefaultSevenTvGqlV3BaseUrl,
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

func (c *SevenTvClient) GetTwitchConnection(twitchUserID string) (ConnectionResponse, error) {
	var resp ConnectionResponse
	err := requests.
		URL(DefaultSevenTvV3ApiBaseUrl).
		Pathf("/v3/users/twitch/%s", twitchUserID).
		ToJSON(&resp).
		Fetch(context.Background())
	if err != nil {
		return ConnectionResponse{}, err
	}

	return resp, nil
}

func (c *SevenTvClient) GetEmote(emoteID string) (Emote, error) {
	var emoteData sevenTvEmote

	err := requests.URL(c.apiBaseUrl + "/emotes/" + emoteID).
		ToJSON(&emoteData).
		Fetch(context.Background())

	if !emoteData.Listed {
		return Emote{}, fmt.Errorf("emote %s is not listed", emoteData.Name)
	}

	return Emote{Code: emoteData.Name, ID: emoteData.ID}, err
}

type ChangeEmoteResponse struct {
	Errors []struct {
		Message    string   `json:"message"`
		Path       []string `json:"path"`
		Extensions struct {
			Code   int `json:"code"`
			Fields struct {
			} `json:"fields"`
			Message string `json:"message"`
		} `json:"extensions"`
	} `json:"errors"`
	Data struct {
		EmoteSet interface{} `json:"emoteSet"`
	} `json:"data"`
}

func (c *SevenTvClient) RemoveEmote(channelUserID, emoteID string) error {
	connection, err := c.GetTwitchConnection(channelUserID)
	if err != nil {
		return errors.New("Could not find 7TV twitch connection for user " + err.Error())
	}

	var resp ChangeEmoteResponse
	err = c.QuerySevenTvGQL(
		`mutation addEmote($emoteSet: ObjectID!, $emoteId: ObjectID!) {
			emoteSet(id: $emoteSet) {
				emotes(id: $emoteId, action: REMOVE) {
					id
					name
				}
			}
		}`,
		map[string]interface{}{
			"emoteId":  emoteID,
			"emoteSet": connection.EmoteSet.ID,
		}, &resp,
		true,
	)

	if len(resp.Errors) > 0 {
		log.Errorf("7tv GQL error: %v", resp)
		errorMessages := make([]string, 0)
		for _, err := range resp.Errors {
			errorMessages = append(errorMessages, err.Message)
		}

		return errors.New(strings.Join(errorMessages, ", "))
	}

	return err
}

func (c *SevenTvClient) AddEmote(channelUserID, emoteID string) error {
	connection, err := c.GetTwitchConnection(channelUserID)
	if err != nil {
		return errors.New("Could not find 7TV twitch connection for user " + err.Error())
	}

	var resp ChangeEmoteResponse
	err = c.QuerySevenTvGQL(
		`mutation addEmote($emoteSet: ObjectID!, $emoteId: ObjectID!) {
			emoteSet(id: $emoteSet) {
				emotes(id: $emoteId, action: ADD) {
					id
					name
				}
			}
		}`,
		map[string]interface{}{
			"emoteId":  emoteID,
			"emoteSet": connection.EmoteSet.ID,
		}, &resp,
		true,
	)

	if len(resp.Errors) > 0 {
		log.Errorf("7tv GQL error: %v", resp)
		errorMessages := make([]string, 0)
		for _, err := range resp.Errors {
			errorMessages = append(errorMessages, err.Message)
		}

		return errors.New(strings.Join(errorMessages, ", "))
	}

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
	`, map[string]interface{}{"id": channelID}, &userData, false)
	if err != nil {
		return User{}, err
	}

	var emotes []Emote
	for _, emote := range userData.Data.User.Emotes {
		emotes = append(emotes, Emote{ID: emote.ID, Code: emote.Name})
	}

	return User{ID: userData.Data.User.ID, Emotes: emotes, EmoteSlots: userData.Data.User.EmoteSlots}, nil
}

func (c *SevenTvClient) QuerySevenTvGQL(query string, variables map[string]interface{}, response interface{}, v3 bool) error {
	gqlQuery := gqlQuery{Query: query, Variables: variables}

	gqlBaseUrl := c.gqlBaseUrl
	if v3 {
		gqlBaseUrl = c.gqlV3BaseUrl
	}

	err := requests.
		URL(gqlBaseUrl).
		BodyJSON(gqlQuery).
		Bearer(c.store.GetSevenTvToken(context.Background())).
		ToJSON(&response).
		Fetch(context.Background())
	if err != nil {
		log.Infof("7TV query '%s' with '%v' resp: '%v'", query, variables, response)
		return err
	}

	log.Infof("7tv query '%s' with '%v' resp: '%v'", query, variables, response)

	return nil
}

const SEVEN_TV_ADD_EMOTE_QUERY = `mutation AddChannelEmote($ch: String!, $em: String!, $re: String!) {addChannelEmote(channel_id: $ch, emote_id: $em, reason: $re) {emote_ids}}`
