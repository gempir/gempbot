package helix

import (
	"fmt"

	nickHelix "github.com/nicklaw5/helix"
)

func (c *Client) CreateChannelPointsRewardAdd(userID string, webHookUrl string) (*nickHelix.EventSubSubscriptionsResponse, error) {
	// Twitch doesn't need a user token here, always an app token eventhough the user has to authenticate beforehand.
	// Internally they check if the app token has authenticated users
	response, err := c.Client.CreateEventSubSubscription(
		&nickHelix.EventSubSubscription{
			Condition: nickHelix.EventSubCondition{BroadcasterUserID: userID},
			Transport: nickHelix.EventSubTransport{Method: "webhook", Callback: webHookUrl, Secret: c.eventSubSecret},
			Type:      "channel.channel_points_custom_reward_redemption.add",
			Version:   "1",
		},
	)

	c.helixApiResponseStatus.WithLabelValues(fmt.Sprint(response.StatusCode), "CreateChannelPointsRewardAdd").Inc()

	return response, err
}
