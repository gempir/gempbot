package store

import (
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
	RewardID                          string
}

func (db *Database) GetChannelPointRewards(userID string) []ChannelPointReward {
	var rewards []ChannelPointReward

	db.Client.Where("owner_twitch_id = ?", userID).Find(&rewards)

	return rewards
}

func (db *Database) GetChannelPointReward(userID string, rewardType string) ChannelPointReward {
	var reward ChannelPointReward
	db.Client.Where("owner_twitch_id = ? AND type = ?", userID, rewardType).First(&reward)

	return reward
}

func (db *Database) SaveReward(reward ChannelPointReward) error {
	update := db.Client.Model(&reward).Where("owner_twitch_id = ? AND type = ?", reward.OwnerTwitchID, reward.Type).Updates(&reward)
	if update.Error != nil {
		return update.Error
	}

	if update.RowsAffected > 0 {
		return nil
	}

	update = db.Client.Create(&reward)

	return update.Error
}
