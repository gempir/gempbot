package helixclient

import (
	"time"

	"github.com/gempir/gempbot/internal/log"
	"github.com/nicklaw5/helix/v2"
)

func (c *HelixClient) CreateEventSubSubscription(userID string, webHookUrl string, subType string) (*helix.EventSubSubscriptionsResponse, error) {
	c.Client.SetAppAccessToken(c.AppAccessToken.AccessToken)
	c.Client.SetUserAccessToken("")
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

func (c *HelixClient) CreateRewardEventSubSubscription(userID, webHookUrl, subType, rewardID string) (*helix.EventSubSubscriptionsResponse, error) {
	c.Client.SetAppAccessToken(c.AppAccessToken.AccessToken)
	c.Client.SetUserAccessToken("")
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

func (c *HelixClient) GetAllSubscriptions(eventType string) []helix.EventSubSubscription {
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

func (c *HelixClient) RemoveEventSubSubscription(id string) (*helix.RemoveEventSubSubscriptionParamsResponse, error) {
	return c.Client.RemoveEventSubSubscription(id)
}

func (c *HelixClient) GetEventSubSubscriptions(params *helix.EventSubSubscriptionsParams) (*helix.EventSubSubscriptionsResponse, error) {
	return c.Client.GetEventSubSubscriptions(params)
}
