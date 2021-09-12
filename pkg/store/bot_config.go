package store

import (
	"context"
	"errors"
)

type BotConfig struct {
	OwnerTwitchID string `gorm:"primaryKey"`
	JoinBot       bool   `gorm:"index"`
}

func (db *Database) SaveBotConfig(ctx context.Context, botCfg BotConfig) error {
	updateMap, err := StructToMap(botCfg)
	if err != nil {
		return err
	}

	update := db.Client.WithContext(ctx).Model(&botCfg).Where("owner_twitch_id = ?", botCfg.OwnerTwitchID).Updates(&updateMap)
	if update.Error != nil {
		return update.Error
	}

	if update.RowsAffected > 0 {
		return nil
	}

	update = db.Client.Create(&botCfg)

	return update.Error
}

func (db *Database) GetAllJoinBotConfigs() []BotConfig {
	var botConfigs []BotConfig

	db.Client.Where("join_bot = true").Find(&botConfigs)

	return botConfigs
}

func (db *Database) GetBotConfig(userID string) (BotConfig, error) {
	var botConfig BotConfig
	result := db.Client.Where("owner_twitch_id = ?", userID).First(&botConfig)
	if result.RowsAffected == 0 {
		return botConfig, errors.New("not found")
	}

	return botConfig, nil
}
