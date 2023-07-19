package emoteservice

import (
	"context"
	"errors"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

const DefaultSevenTvApiBaseUrl = "https://7tv.io/v3"
const DefaultSevenTvGqlV3BaseUrl = "https://api.7tv.app/v3/gql"

type SevenTvClient struct {
	store      store.Store
	apiBaseUrl string
	gqlBaseUrl string
}

func NewSevenTvClient(store store.Store) *SevenTvClient {
	return &SevenTvClient{
		store:      store,
		apiBaseUrl: DefaultSevenTvApiBaseUrl,
		gqlBaseUrl: DefaultSevenTvGqlV3BaseUrl,
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
		URL(c.apiBaseUrl).
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
		URL(c.apiBaseUrl).
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
	var userResp UserResponse

	err := requests.URL(c.apiBaseUrl + "/users/twitch/" + channelID).
		ToJSON(&userResp).
		Fetch(context.Background())
	if err != nil {
		return User{}, err
	}

	var emotes []Emote
	for _, emote := range userResp.EmoteSet.Emotes {
		emotes = append(emotes, Emote{ID: emote.ID, Code: emote.Name})
	}

	return User{ID: userResp.User.ID, Emotes: emotes, EmoteSlots: userResp.EmoteCapacity}, nil
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
		log.Infof("7TV query '%s' with '%v' resp: '%v'", query, variables, response)
		return err
	}

	log.Infof("7tv query '%s' with '%v' resp: '%v'", query, variables, response)

	return nil
}

const SEVEN_TV_ADD_EMOTE_QUERY = `mutation AddChannelEmote($ch: String!, $em: String!, $re: String!) {addChannelEmote(channel_id: $ch, emote_id: $em, reason: $re) {emote_ids}}`

type UserResponse struct {
	ID            string      `json:"id"`
	Platform      string      `json:"platform"`
	Username      string      `json:"username"`
	DisplayName   string      `json:"display_name"`
	LinkedAt      int64       `json:"linked_at"`
	EmoteCapacity int         `json:"emote_capacity"`
	EmoteSetID    interface{} `json:"emote_set_id"`
	EmoteSet      struct {
		ID         string        `json:"id"`
		Name       string        `json:"name"`
		Flags      int           `json:"flags"`
		Tags       []interface{} `json:"tags"`
		Immutable  bool          `json:"immutable"`
		Privileged bool          `json:"privileged"`
		Emotes     []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Flags     int    `json:"flags"`
			Timestamp int64  `json:"timestamp"`
			ActorID   string `json:"actor_id"`
			Data      struct {
				ID        string   `json:"id"`
				Name      string   `json:"name"`
				Flags     int      `json:"flags"`
				Lifecycle int      `json:"lifecycle"`
				State     []string `json:"state"`
				Listed    bool     `json:"listed"`
				Animated  bool     `json:"animated"`
				Owner     struct {
					ID          string `json:"id"`
					Username    string `json:"username"`
					DisplayName string `json:"display_name"`
					AvatarURL   string `json:"avatar_url"`
					Style       struct {
					} `json:"style"`
					Roles []string `json:"roles"`
				} `json:"owner"`
				Host struct {
					URL   string `json:"url"`
					Files []struct {
						Name       string `json:"name"`
						StaticName string `json:"static_name"`
						Width      int    `json:"width"`
						Height     int    `json:"height"`
						FrameCount int    `json:"frame_count"`
						Size       int    `json:"size"`
						Format     string `json:"format"`
					} `json:"files"`
				} `json:"host"`
			} `json:"data"`
		} `json:"emotes"`
		EmoteCount int `json:"emote_count"`
		Capacity   int `json:"capacity"`
		Owner      struct {
			ID          string `json:"id"`
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			AvatarURL   string `json:"avatar_url"`
			Style       struct {
				Color int `json:"color"`
			} `json:"style"`
			Roles []string `json:"roles"`
		} `json:"owner"`
	} `json:"emote_set"`
	User struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		CreatedAt   int64  `json:"created_at"`
		AvatarURL   string `json:"avatar_url"`
		Biography   string `json:"biography"`
		Style       struct {
			Color int `json:"color"`
		} `json:"style"`
		EmoteSets []struct {
			ID       string        `json:"id"`
			Name     string        `json:"name"`
			Flags    int           `json:"flags"`
			Tags     []interface{} `json:"tags"`
			Capacity int           `json:"capacity"`
		} `json:"emote_sets"`
		Editors []struct {
			ID          string `json:"id"`
			Permissions int    `json:"permissions"`
			Visible     bool   `json:"visible"`
			AddedAt     int64  `json:"added_at"`
		} `json:"editors"`
		Roles       []string `json:"roles"`
		Connections []struct {
			ID            string      `json:"id"`
			Platform      string      `json:"platform"`
			Username      string      `json:"username"`
			DisplayName   string      `json:"display_name"`
			LinkedAt      int64       `json:"linked_at"`
			EmoteCapacity int         `json:"emote_capacity"`
			EmoteSetID    interface{} `json:"emote_set_id"`
			EmoteSet      struct {
				ID         string        `json:"id"`
				Name       string        `json:"name"`
				Flags      int           `json:"flags"`
				Tags       []interface{} `json:"tags"`
				Immutable  bool          `json:"immutable"`
				Privileged bool          `json:"privileged"`
				Capacity   int           `json:"capacity"`
				Owner      interface{}   `json:"owner"`
			} `json:"emote_set"`
		} `json:"connections"`
	} `json:"user"`
}
