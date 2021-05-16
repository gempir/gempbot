package helix

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type CreateCustomRewardRequest struct {
	Title                             string `json:"title"`
	Prompt                            string `json:"prompt,omitempty"`
	Cost                              int    `json:"cost"`
	IsEnabled                         bool   `json:"is_enabled"`
	BackgroundColor                   string `json:"background_color,omitempty"`
	IsUserInputRequired               bool   `json:"is_user_input_required,omitempty"`
	IsMaxPerStreamEnabled             bool   `json:"is_max_per_stream_enabled"`
	MaxPerStream                      int    `json:"max_per_stream,omitempty"`
	IsMaxPerUserPerStreamEnabled      bool   `json:"is_max_per_user_per_stream_enabled"`
	MaxPerUserPerStream               int    `json:"max_per_user_per_stream,omitempty"`
	IsGlobalCooldownEnabled           bool   `json:"is_global_cooldown_enabled"`
	GlobalCoolDownSeconds             int    `json:"global_cooldown_seconds,omitempty"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CreateCustomRewardResponse struct {
	Data []CreateCustomRewardResponseDataItem `json:"data"`
}

type CreateCustomRewardResponseDataItem struct {
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
}

func (c *Client) CreateOrUpdateReward(userID, userAccessToken string, reward CreateCustomRewardRequest, rewardID string) (CreateCustomRewardResponseDataItem, error) {
	log.Infof("Updating Reward for user %s reward title: %s", userID, reward.Title)
	marshalled, err := json.Marshal(reward)
	if err != nil {
		return CreateCustomRewardResponseDataItem{}, err
	}

	url := "https://api.twitch.tv/helix/channel_points/custom_rewards?broadcaster_id=" + userID
	method := http.MethodPost
	if rewardID != "" {
		url = url + "&id=" + rewardID
		method = http.MethodPatch
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(marshalled))
	req.Header.Set("authorization", "Bearer "+userAccessToken)
	req.Header.Set("client-id", c.clientID)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Error(err)
		return CreateCustomRewardResponseDataItem{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return CreateCustomRewardResponseDataItem{}, err
	}

	log.Infof("[%d][%s] %s", resp.StatusCode, method, url)

	if resp.StatusCode >= 400 {
		var response ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return CreateCustomRewardResponseDataItem{}, fmt.Errorf("Failed to unmarshal reward error response: %s", err.Error())
		}

		return CreateCustomRewardResponseDataItem{}, fmt.Errorf("Failed to create reward: %s", response.Message)
	}

	var response CreateCustomRewardResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return CreateCustomRewardResponseDataItem{}, err
	}

	if len(response.Data) != 1 {
		return CreateCustomRewardResponseDataItem{}, fmt.Errorf("%d amount of rewards returned after creation", len(response.Data))
	}

	return response.Data[0], nil
}
