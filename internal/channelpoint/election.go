package channelpoint

import "github.com/gempir/gempbot/internal/dto"

type NominateReward struct {
	TwitchRewardConfig
	NominateAdditionalOptions
}

type NominateAdditionalOptions struct {
}

func (r *NominateReward) GetType() dto.RewardType {
	return dto.REWARD_SEVENTV
}

func (r *NominateReward) GetAdditionalOptions() interface{} {
	return r.NominateAdditionalOptions
}

func (r *NominateReward) GetConfig() TwitchRewardConfig {
	return r.TwitchRewardConfig
}

func (r *NominateReward) SetConfig(config TwitchRewardConfig) {
	r.TwitchRewardConfig = config
}
