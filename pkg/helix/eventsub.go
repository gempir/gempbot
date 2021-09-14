package helix

import (
	nickHelix "github.com/nicklaw5/helix"
)

func (c *Client) CreateEventSubSubscription(userID string, webHookUrl string, subType string) (*nickHelix.EventSubSubscriptionsResponse, error) {
	// Twitch doesn't need a user token here, always an app token eventhough the user has to authenticate beforehand.
	// Internally they check if the app token has authenticated users
	response, err := c.Client.CreateEventSubSubscription(
		&nickHelix.EventSubSubscription{
			Condition: nickHelix.EventSubCondition{BroadcasterUserID: userID},
			Transport: nickHelix.EventSubTransport{Method: "webhook", Callback: webHookUrl, Secret: c.eventSubSecret},
			Type:      subType,
			Version:   "1",
		},
	)

	return response, err
}
