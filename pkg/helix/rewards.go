package helix

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type GetRewardsResponse struct {
	Data []struct {
		BroadcasterName     string      `json:"broadcaster_name"`
		BroadcasterLogin    string      `json:"broadcaster_login"`
		BroadcasterID       string      `json:"broadcaster_id"`
		ID                  string      `json:"id"`
		Image               interface{} `json:"image"`
		BackgroundColor     string      `json:"background_color"`
		IsEnabled           bool        `json:"is_enabled"`
		Cost                int         `json:"cost"`
		Title               string      `json:"title"`
		Prompt              string      `json:"prompt"`
		IsUserInputRequired bool        `json:"is_user_input_required"`
		MaxPerStreamSetting struct {
			IsEnabled    bool `json:"is_enabled"`
			MaxPerStream int  `json:"max_per_stream"`
		} `json:"max_per_stream_setting"`
		MaxPerUserPerStreamSetting struct {
			IsEnabled           bool `json:"is_enabled"`
			MaxPerUserPerStream int  `json:"max_per_user_per_stream"`
		} `json:"max_per_user_per_stream_setting"`
		GlobalCooldownSetting struct {
			IsEnabled             bool `json:"is_enabled"`
			GlobalCooldownSeconds int  `json:"global_cooldown_seconds"`
		} `json:"global_cooldown_setting"`
		IsPaused     bool `json:"is_paused"`
		IsInStock    bool `json:"is_in_stock"`
		DefaultImage struct {
			URL1X string `json:"url_1x"`
			URL2X string `json:"url_2x"`
			URL4X string `json:"url_4x"`
		} `json:"default_image"`
		ShouldRedemptionsSkipRequestQueue bool        `json:"should_redemptions_skip_request_queue"`
		RedemptionsRedeemedCurrentStream  interface{} `json:"redemptions_redeemed_current_stream"`
		CooldownExpiresAt                 interface{} `json:"cooldown_expires_at"`
	} `json:"data"`
}

func (c *Client) GetRewards(userID, userAccessToken string) (*GetRewardsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/channel_points/custom_rewards?broadcaster_id="+userID, nil)
	req.Header.Set("authorization", "Bearer "+userAccessToken)
	req.Header.Set("client-id", c.clientID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var getRewardsResponse GetRewardsResponse
	err = json.NewDecoder(resp.Body).Decode(&getRewardsResponse)
	if err != nil {
		return nil, err
	}

	return &getRewardsResponse, nil
}
