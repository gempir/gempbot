package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gempir/bitraft/pkg/dto"
	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/log"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/labstack/echo/v4"
)

type Reward interface {
	GetType() string
	GetConfig() TwitchRewardConfig
	SetConfig(config TwitchRewardConfig)
	GetAdditionalOptions() interface{}
}

type TwitchRewardConfig struct {
	Title                             string `json:"title"`
	Prompt                            string `json:"prompt"`
	Cost                              int    `json:"cost"`
	BackgroundColor                   string `json:"backgroundColor"`
	IsMaxPerStreamEnabled             bool   `json:"isMaxPerStreamEnabled"`
	MaxPerStream                      int    `json:"maxPerStream"`
	IsUserInputRequired               bool   `json:"isUserInputRequired"`
	IsMaxPerUserPerStreamEnabled      bool   `json:"isMaxPerUserPerStreamEnabled"`
	MaxPerUserPerStream               int    `json:"maxPerUserPerStream"`
	IsGlobalCooldownEnabled           bool   `json:"isGlobalCooldownEnabled"`
	GlobalCooldownSeconds             int    `json:"globalCooldownSeconds"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"shouldRedemptionsSkipRequestQueue"`
	Enabled                           bool   `json:"enabled"`
	ID                                string
}

type BttvReward struct {
	TwitchRewardConfig
	BttvAdditionalOptions
}

type BttvAdditionalOptions struct {
	Slots int
}

func (r *BttvReward) GetType() string {
	return dto.REWARD_BTTV
}

func (r *BttvReward) GetAdditionalOptions() interface{} {
	return r.BttvAdditionalOptions
}

func (r *BttvReward) GetConfig() TwitchRewardConfig {
	return r.TwitchRewardConfig
}

func (r *BttvReward) SetConfig(config TwitchRewardConfig) {
	r.TwitchRewardConfig = config
}

type SevenTvReward struct {
	TwitchRewardConfig
	SevenTvAdditionalOptions
}

type SevenTvAdditionalOptions struct {
	Slots int
}

func (r *SevenTvReward) GetType() string {
	return dto.REWARD_SEVENTV
}

func (r *SevenTvReward) GetAdditionalOptions() interface{} {
	return r.SevenTvAdditionalOptions
}

func (r *SevenTvReward) GetConfig() TwitchRewardConfig {
	return r.TwitchRewardConfig
}

func (r *SevenTvReward) SetConfig(config TwitchRewardConfig) {
	r.TwitchRewardConfig = config
}

type TimeoutReward struct {
	TwitchRewardConfig
	TimeoutAdditionalOptions
}

type TimeoutAdditionalOptions struct {
	Length int
}

func (r *TimeoutReward) GetType() string {
	return dto.REWARD_TIMEOUT
}

func (r *TimeoutReward) GetConfig() TwitchRewardConfig {
	return r.TwitchRewardConfig
}

func (r *TimeoutReward) SetConfig(config TwitchRewardConfig) {
	r.TwitchRewardConfig = config
}

func (r *TimeoutReward) GetAdditionalOptions() interface{} {
	return r.TimeoutAdditionalOptions
}

func MarshallReward(reward Reward) string {
	js, err := json.Marshal(reward)
	if err != nil {
		log.Infof("failed to marshal BttvReward %s", err)
		return ""
	}

	return string(js)
}

func (s *Server) handleRewardDeletion(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.Param("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.Param("userID"))
		if err != nil {
			return err
		}
	}

	reward, err := s.db.GetChannelPointReward(c.Param("userID"), c.Param("type"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	token, err := s.db.GetUserAccessToken(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no accessToken to edit reward")
	}

	s.db.DeleteChannelPointReward(c.Param("userID"), c.Param("type"))

	err = s.helixClient.DeleteReward(c.Param("userID"), token.AccessToken, reward.RewardID)
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (s *Server) handleRewardSingleRead(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.Param("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.Param("userID"))
		if err != nil {
			return err
		}
	}

	reward, err := s.db.GetChannelPointReward(c.Param("userID"), c.Param("type"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, reward)
}

func (s *Server) handleRewardCreateOrUpdate(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.Param("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.Param("userID"))
		if err != nil {
			return err
		}
	}

	newReward, err := createRewardFromBody(c.Request().Body)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failure reading body")
	}

	rewardID := ""
	reward, err := s.db.GetChannelPointReward(c.Param("userID"), newReward.GetType())
	if err == nil {
		rewardID = reward.RewardID
	}

	config, err := s.createOrUpdateChannelPointReward(c.Param("userID"), newReward.GetConfig(), rewardID)
	if err != nil {
		log.Errorf("Failed saving reward to twitch: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failure saving reward to twitch")
	}

	newReward.SetConfig(config)

	err = s.db.SaveReward(createStoreRewardFromReward(c.Param("userID"), newReward))
	if err != nil {
		log.Errorf("Failed saving reward: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failure saving reward")
	}

	_ = s.subscribeChannelPoints(c.Param("userID"))

	return nil
}

func createStoreRewardFromReward(userID string, reward Reward) store.ChannelPointReward {
	addOpts, err := json.Marshal(reward.GetAdditionalOptions())
	if err != nil {
		log.Error(err)
	}

	return store.ChannelPointReward{
		OwnerTwitchID:                     userID,
		Type:                              reward.GetType(),
		RewardID:                          reward.GetConfig().ID,
		Title:                             reward.GetConfig().Title,
		Prompt:                            reward.GetConfig().Prompt,
		Cost:                              reward.GetConfig().Cost,
		BackgroundColor:                   reward.GetConfig().BackgroundColor,
		IsMaxPerStreamEnabled:             reward.GetConfig().IsMaxPerStreamEnabled,
		MaxPerStream:                      reward.GetConfig().MaxPerStream,
		IsUserInputRequired:               reward.GetConfig().IsUserInputRequired,
		IsMaxPerUserPerStreamEnabled:      reward.GetConfig().IsMaxPerUserPerStreamEnabled,
		MaxPerUserPerStream:               reward.GetConfig().MaxPerUserPerStream,
		IsGlobalCooldownEnabled:           reward.GetConfig().IsGlobalCooldownEnabled,
		GlobalCooldownSeconds:             reward.GetConfig().GlobalCooldownSeconds,
		ShouldRedemptionsSkipRequestQueue: reward.GetConfig().ShouldRedemptionsSkipRequestQueue,
		Enabled:                           reward.GetConfig().Enabled,
		AdditionalOptions:                 string(addOpts),
	}
}

type rewardRequestBody struct {
	ID                                string
	OwnerTwitchID                     string
	Type                              string
	RewardID                          string
	CreatedAt                         time.Time
	UpdatedAt                         time.Time
	Title                             string
	Prompt                            string
	Cost                              int
	BackgroundColor                   string
	IsMaxPerStreamEnabled             bool
	MaxPerStream                      int
	IsUserInputRequired               bool
	IsMaxPerUserPerStreamEnabled      bool
	MaxPerUserPerStream               int
	IsGlobalCooldownEnabled           bool
	GlobalCooldownSeconds             int
	ShouldRedemptionsSkipRequestQueue bool
	Enabled                           bool
}

type bttvRewardRequestBody struct {
	AdditionalOptionsParsed BttvAdditionalOptions
}

type sevenTvRewardRequestBody struct {
	AdditionalOptionsParsed SevenTvAdditionalOptions
}

func createRewardFromBody(body io.ReadCloser) (Reward, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var data rewardRequestBody
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}

	rewardConfig := createTwitchRewardConfigFromRequestBody(data)

	switch data.Type {
	case dto.REWARD_BTTV:
		var addOpts bttvRewardRequestBody
		if err := json.Unmarshal(bodyBytes, &addOpts); err != nil {
			return nil, err
		}

		if addOpts.AdditionalOptionsParsed.Slots < 1 {
			addOpts.AdditionalOptionsParsed.Slots = 1
		}

		return &BttvReward{
			TwitchRewardConfig:    rewardConfig,
			BttvAdditionalOptions: addOpts.AdditionalOptionsParsed,
		}, nil
	case dto.REWARD_SEVENTV:
		var addOpts sevenTvRewardRequestBody
		if err := json.Unmarshal(bodyBytes, &addOpts); err != nil {
			return nil, err
		}

		if addOpts.AdditionalOptionsParsed.Slots < 1 {
			addOpts.AdditionalOptionsParsed.Slots = 1
		}

		return &SevenTvReward{
			TwitchRewardConfig:       rewardConfig,
			SevenTvAdditionalOptions: addOpts.AdditionalOptionsParsed,
		}, nil
	case dto.REWARD_TIMEOUT:
		return &TimeoutReward{
			TwitchRewardConfig: createTwitchRewardConfigFromRequestBody(data),
		}, nil
	}

	return nil, errors.New("unknown reward")
}

func createTwitchRewardConfigFromRequestBody(body rewardRequestBody) TwitchRewardConfig {
	return TwitchRewardConfig{
		Title:                             body.Title,
		Prompt:                            body.Prompt,
		Cost:                              body.Cost,
		BackgroundColor:                   body.BackgroundColor,
		IsMaxPerStreamEnabled:             body.IsMaxPerStreamEnabled,
		MaxPerStream:                      body.MaxPerStream,
		IsUserInputRequired:               true,
		IsMaxPerUserPerStreamEnabled:      body.IsMaxPerUserPerStreamEnabled,
		MaxPerUserPerStream:               body.MaxPerUserPerStream,
		IsGlobalCooldownEnabled:           body.IsGlobalCooldownEnabled,
		GlobalCooldownSeconds:             body.GlobalCooldownSeconds,
		ShouldRedemptionsSkipRequestQueue: false,
		Enabled:                           body.Enabled,
		ID:                                body.ID,
	}
}

func (s *Server) createOrUpdateChannelPointReward(userID string, request TwitchRewardConfig, rewardID string) (TwitchRewardConfig, error) {
	token, err := s.db.GetUserAccessToken(userID)
	if err != nil {
		return TwitchRewardConfig{}, err
	}

	req := helix.CreateCustomRewardRequest{
		Title:                             request.Title,
		Prompt:                            request.Prompt,
		Cost:                              request.Cost,
		IsEnabled:                         request.Enabled,
		BackgroundColor:                   request.BackgroundColor,
		IsUserInputRequired:               true,
		ShouldRedemptionsSkipRequestQueue: false,
		IsMaxPerStreamEnabled:             false,
		IsMaxPerUserPerStreamEnabled:      false,
		IsGlobalCooldownEnabled:           false,
	}

	if request.MaxPerStream != 0 {
		req.IsMaxPerStreamEnabled = true
		req.MaxPerStream = request.MaxPerStream
	}

	if request.MaxPerUserPerStream != 0 {
		req.IsMaxPerUserPerStreamEnabled = true
		req.MaxPerUserPerStream = request.MaxPerUserPerStream
	}

	if request.GlobalCooldownSeconds != 0 {
		req.IsGlobalCooldownEnabled = true
		req.GlobalCoolDownSeconds = request.GlobalCooldownSeconds
	}

	resp, err := s.helixClient.CreateOrUpdateReward(userID, token.AccessToken, req, rewardID)
	if err != nil {
		return TwitchRewardConfig{}, err
	}

	return TwitchRewardConfig{
		Title:                             resp.Title,
		Prompt:                            resp.Prompt,
		Cost:                              resp.Cost,
		BackgroundColor:                   resp.BackgroundColor,
		IsMaxPerStreamEnabled:             resp.MaxPerStreamSetting.IsEnabled,
		MaxPerStream:                      resp.MaxPerStreamSetting.MaxPerStream,
		IsUserInputRequired:               resp.IsUserInputRequired,
		IsMaxPerUserPerStreamEnabled:      resp.MaxPerUserPerStreamSetting.IsEnabled,
		MaxPerUserPerStream:               resp.MaxPerUserPerStreamSetting.MaxPerUserPerStream,
		IsGlobalCooldownEnabled:           resp.GlobalCooldownSetting.IsEnabled,
		GlobalCooldownSeconds:             resp.GlobalCooldownSetting.GlobalCooldownSeconds,
		ShouldRedemptionsSkipRequestQueue: resp.ShouldRedemptionsSkipRequestQueue,
		Enabled:                           resp.IsEnabled,
		ID:                                resp.ID,
	}, nil
}

func UnmarshallTimeoutAdditionalOptions(data string) (TimeoutAdditionalOptions, error) {
	var addOpts TimeoutAdditionalOptions
	err := json.Unmarshal([]byte(data), &addOpts)
	if err != nil {
		return TimeoutAdditionalOptions{}, err
	}

	return addOpts, nil
}

func UnmarshallBttvAdditionalOptions(jsonString string) BttvAdditionalOptions {
	if jsonString == "{}" {
		return BttvAdditionalOptions{Slots: 1}
	}

	var additionalOptions BttvAdditionalOptions

	if err := json.Unmarshal([]byte(jsonString), &additionalOptions); err != nil {
		log.Error(err)
		return BttvAdditionalOptions{Slots: 1}
	}

	return additionalOptions
}

func UnmarshallSevenTvAdditionalOptions(jsonString string) SevenTvAdditionalOptions {
	if jsonString == "{}" {
		return SevenTvAdditionalOptions{Slots: 1}
	}

	var additionalOptions SevenTvAdditionalOptions

	if err := json.Unmarshal([]byte(jsonString), &additionalOptions); err != nil {
		log.Error(err)
		return SevenTvAdditionalOptions{Slots: 1}
	}

	return additionalOptions
}
