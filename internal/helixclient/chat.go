package helixclient

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/nicklaw5/helix/v2"
)

type SendChatMessageResponse struct {
	Data []struct {
		MessageID string `json:"message_id"`
		IsSent    bool   `json:"is_sent"`
	} `json:"data"`
}

func (c *HelixClient) SendChatMessage(channelID string, message string) (*SendChatMessageResponse, error) {
	var resp SendChatMessageResponse

	params := helix.SendChatMessageParams{
		BroadcasterID: channelID,
		Message:       message,
		SenderID:      c.botUserID,
	}

	err := requests.
		URL(TWITCH_API).
		BodyJSON(params).
		Bearer(c.AppAccessToken.AccessToken).
		Header("Client-Id", c.clientID).
		ContentType("application/json").
		Path("/helix/chat/messages").
		ToJSON(&resp).
		Post().
		Fetch(context.Background())

	return &resp, err
}
