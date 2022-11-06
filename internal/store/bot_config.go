package store

import (
	"context"
	"errors"

	"gorm.io/gorm/clause"
)

type BotConfig struct {
	OwnerTwitchID string `gorm:"primaryKey"`
	JoinBot       bool   `gorm:"index"`
	MediaCommands bool
}

func (db *Database) SaveBotConfig(ctx context.Context, botCfg BotConfig) error {
	update := db.Client.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&botCfg)

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

func (db *Database) GetAllMediaCommandsBotConfig() []BotConfig {
	var botConfigs []BotConfig

	db.Client.Where("media_commands = true").Find(&botConfigs)

	return botConfigs
}
