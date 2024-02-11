package helixclient

import "github.com/nicklaw5/helix/v2"

func (c *HelixClient) SendChatMessage(params *helix.SendChatMessageParams) (*helix.ChatMessageResponse, error) {
	return c.Client.SendChatMessage(params)
}
