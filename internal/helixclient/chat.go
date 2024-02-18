package helixclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nicklaw5/helix/v2"
)

type SendChatMessageResponse struct {
	Data []struct {
		MessageID string `json:"message_id"`
		IsSent    bool   `json:"is_sent"`
	} `json:"data"`
}

func (c *HelixClient) SendChatMessage(channelID string, message string) (*SendChatMessageResponse, error) {
	params := helix.SendChatMessageParams{
		BroadcasterID: channelID,
		Message:       message,
		SenderID:      c.botUserID,
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/helix/chat/messages", TWITCH_API), bytes.NewBuffer(jsonParams))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AppAccessToken.AccessToken))
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("[%d] %s", resp.StatusCode, string(body))
	}

	var msgs SendChatMessageResponse
	err = json.Unmarshal(body, &msgs)
	if err != nil {
		return nil, err
	}

	return &msgs, nil
}
