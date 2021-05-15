package helix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	IsEnabled                         bool   `json:"is_enabled,omitempty"`
	BackgroundColor                   string `json:"background_color,omitempty"`
	IsUserInputRequired               bool   `json:"is_user_input_required,omitempty"`
	IsMaxPerStreamEnabled             bool   `json:"is_max_per_stream_enabled,omitempty"`
	MaxPerStream                      int    `json:"max_per_stream,omitempty"`
	IsMaxPerUserPerStreamEnabled      bool   `json:"is_max_per_user_per_stream_enabled,omitempty"`
	MaxPerUserPerStream               int    `json:"max_per_user_per_stream,omitempty"`
	IsGlobalCooldownEnabled           bool   `json:"is_global_cooldown_enabled,omitempty"`
	GlobalCoolDownSeconds             int    `json:"global_cooldown_seconds,omitempty"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue,omitempty"`
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

func (c *Client) CreateReward(userID, userAccessToken string, reward CreateCustomRewardRequest) (CreateCustomRewardResponseDataItem, error) {
	marshalled, err := json.Marshal(reward)
	if err != nil {
		return CreateCustomRewardResponseDataItem{}, err
	}

	log.Info(bytes.NewBuffer(marshalled))

	req, err := http.NewRequest(http.MethodPost, "https://api.twitch.tv/helix/channel_points/custom_rewards?broadcaster_id="+userID, bytes.NewBuffer(marshalled))
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

	if resp.StatusCode >= 400 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return CreateCustomRewardResponseDataItem{}, fmt.Errorf("Failed to read bad reward create response %s", err.Error())
		}
		log.Errorf("Failed to create reward %s", string(bodyBytes))
		return CreateCustomRewardResponseDataItem{}, errors.New("Failed to create reward")
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
