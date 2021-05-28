package store

import (
	"errors"

	"gorm.io/gorm"
)

type ChannelPointReward struct {
	gorm.Model
	OwnerTwitchID                     string `gorm:"index"`
	Type                              string `gorm:"index"`
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
	RewardID                          string `gorm:"index"`
}

func (db *Database) GetChannelPointRewards(userID string) []ChannelPointReward {
	var rewards []ChannelPointReward

	db.Client.Where("owner_twitch_id = ?", userID).Find(&rewards)

	return rewards
}

func (db *Database) GetEnabledChannelPointRewardByID(rewardID string) (ChannelPointReward, error) {
	var reward ChannelPointReward
	result := db.Client.Where("reward_id = ? AND enabled = ?", rewardID, true).First(&reward)
	if result.RowsAffected == 0 {
		return reward, errors.New("not found")
	}

	return reward, nil
}

func (db *Database) GetChannelPointReward(userID string, rewardType string) (ChannelPointReward, error) {
	var reward ChannelPointReward
	result := db.Client.Where("owner_twitch_id = ? AND type = ?", userID, rewardType).First(&reward)
	if result.RowsAffected == 0 {
		return reward, errors.New("not found")
	}

	return reward, nil
}

func (db *Database) DeleteChannelPointReward(userID string, rewardType string) {
	db.Client.Where("owner_twitch_id = ? AND type = ?", userID, rewardType).Delete(&ChannelPointReward{})
}

func (db *Database) GetDistinctRewardsPerUser() []ChannelPointReward {
	var rewards []ChannelPointReward
	db.Client.Distinct("owner_twitch_id").Find(&rewards)

	return rewards
}

func (db *Database) SaveReward(reward ChannelPointReward) error {
	update := db.Client.Model(&reward).Where("owner_twitch_id = ? AND type = ?", reward.OwnerTwitchID, reward.Type).Updates(&reward)
	if update.Error != nil {
		return update.Error
	}

	if update.RowsAffected > 0 {
		updateMap := map[string]interface{}{"deleted_at": nil}
		update := db.Client.Model(&reward).Where("owner_twitch_id = ? AND type = ?", reward.OwnerTwitchID, reward.Type).Updates(&updateMap)
		if update.Error != nil {
			return update.Error
		}

		return nil
	}

	update = db.Client.Create(&reward)

	return update.Error
}
