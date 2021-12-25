package helixclient

import (
	"time"

	"github.com/gempir/gempbot/pkg/log"
	"github.com/nicklaw5/helix/v2"
)

func (c *Client) CreateEventSubSubscription(userID string, webHookUrl string, subType string) (*helix.EventSubSubscriptionsResponse, error) {
	// Twitch doesn't need a user token here, always an app token eventhough the user has to authenticate beforehand.
	// Internally they check if the app token has authenticated users
	response, err := c.Client.CreateEventSubSubscription(
		&helix.EventSubSubscription{
			Condition: helix.EventSubCondition{BroadcasterUserID: userID},
			Transport: helix.EventSubTransport{Method: "webhook", Callback: webHookUrl, Secret: c.eventSubSecret},
			Type:      subType,
			Version:   "1",
		},
	)

	return response, err
}

func (c *Client) CreateRewardEventSubSubscription(userID, webHookUrl, subType, rewardID string) (*helix.EventSubSubscriptionsResponse, error) {
	// Twitch doesn't need a user token here, always an app token eventhough the user has to authenticate beforehand.
	// Internally they check if the app token has authenticated users
	response, err := c.Client.CreateEventSubSubscription(
		&helix.EventSubSubscription{
			Condition: helix.EventSubCondition{BroadcasterUserID: userID, RewardID: rewardID},
			Transport: helix.EventSubTransport{Method: "webhook", Callback: webHookUrl, Secret: c.eventSubSecret},
			Type:      subType,
			Version:   "1",
		},
	)

	return response, err
}

func (c *Client) GetAllSubscriptions(eventType string) []helix.EventSubSubscription {
	subs := []helix.EventSubSubscription{}

	resp, err := c.Client.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{Type: eventType})
	if err != nil {
		log.Error(err)
	}
	subs = append(subs, resp.Data.EventSubSubscriptions...)

	for resp.Data.Pagination.Cursor != "" {
		log.Infof("Getting next subscriptions after %s", resp.Data.Pagination.Cursor)
		time.Sleep(time.Second * 5)
		resp, err = c.Client.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{Type: eventType, After: resp.Data.Pagination.Cursor})
		if err != nil {
			log.Error(err)
		}
		subs = append(subs, resp.Data.EventSubSubscriptions...)
	}

	return subs
}
