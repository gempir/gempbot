package channelpoint

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

type ChannelPointManager struct {
	cfg         *config.Config
	helixClient helixclient.Client
	db          *store.Database
}

func NewChannelPointManager(cfg *config.Config, helixClient helixclient.Client, db *store.Database) *ChannelPointManager {
	return &ChannelPointManager{
		cfg:         cfg,
		helixClient: helixClient,
		db:          db,
	}
}

func (cpm *ChannelPointManager) DeleteChannelPointReward(userID, rewardID string) error {
	err := cpm.helixClient.DeleteReward(userID, rewardID)
	if err != nil {
		return err
	}

	cpm.db.DeleteChannelPointRewardById(userID, rewardID)
	return nil
}

func (cpm *ChannelPointManager) CreateOrUpdateChannelPointReward(userID string, request TwitchRewardConfig, rewardID string) (TwitchRewardConfig, error) {
	req := helixclient.CreateCustomRewardRequest{
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

	resp, err := cpm.helixClient.CreateOrUpdateReward(userID, req, rewardID)
	if err != nil {
		return TwitchRewardConfig{}, err
	}

	return TwitchRewardConfig{
		Title:                             resp.Title,
		ApproveOnly:                       request.ApproveOnly,
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

type Reward interface {
	GetType() dto.RewardType
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
	ApproveOnly                       bool
	Enabled                           bool `json:"enabled"`
	ID                                string
}

type BttvReward struct {
	TwitchRewardConfig
	BttvAdditionalOptions
}

type BttvAdditionalOptions struct {
	Slots int
}

func (r *BttvReward) GetType() dto.RewardType {
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

func (r *SevenTvReward) GetType() dto.RewardType {
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

func MarshallReward(reward Reward) string {
	js, err := json.Marshal(reward)
	if err != nil {
		log.Infof("failed to marshal BttvReward %s", err)
		return ""
	}

	return string(js)
}

func CreateStoreRewardFromReward(userID string, reward Reward) store.ChannelPointReward {
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
		ApproveOnly:                       reward.GetConfig().ApproveOnly,
		AdditionalOptions:                 string(addOpts),
	}
}

type rewardRequestBody struct {
	ID                                string
	OwnerTwitchID                     string
	Type                              dto.RewardType
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
	ApproveOnly                       bool
	Enabled                           bool
}

type bttvRewardRequestBody struct {
	AdditionalOptionsParsed BttvAdditionalOptions
}

type sevenTvRewardRequestBody struct {
	AdditionalOptionsParsed SevenTvAdditionalOptions
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
		ApproveOnly:                       body.ApproveOnly,
		Enabled:                           body.Enabled,
		ID:                                body.ID,
	}
}

func CreateRewardFromBody(body io.ReadCloser) (Reward, error) {
	bodyBytes, err := io.ReadAll(body)
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
	}

	return nil, errors.New("unknown reward")
}
