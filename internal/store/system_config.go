package store

import (
	"context"

	"gorm.io/gorm/clause"
)

type SystemConfig struct {
	ConfigKey   string `gorm:"primaryKey"`
	ConfigValue string
}

func (db *Database) Save(ctx context.Context, configKey string, configValue string) error {
	var systemConfig SystemConfig = SystemConfig{
		ConfigKey:   configKey,
		ConfigValue: configValue,
	}

	update := db.Client.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&systemConfig)

	return update.Error
}

func (db *Database) GetConfigValue(ctx context.Context, configKey string) string {
	var cfgs []SystemConfig

	db.Client.WithContext(ctx).Where("configKey = ?", configKey).Find(&cfgs)

	if len(cfgs) < 1 {
		return ""
	}

	return cfgs[0].ConfigValue
}

func (db *Database) GetSevenTvToken(ctx context.Context) string {
	var cfgs []SystemConfig

	db.Client.WithContext(ctx).Where("config_key = ?", "SEVEN_TV_TOKEN").Find(&cfgs)

	if len(cfgs) < 1 {
		return ""
	}

	return cfgs[0].ConfigValue
}

func (db *Database) GetBttvToken(ctx context.Context) string {
	var cfgs []SystemConfig

	db.Client.WithContext(ctx).Where("config_key = ?", "BTTV_TOKEN").Find(&cfgs)

	if len(cfgs) < 1 {
		return ""
	}

	return cfgs[0].ConfigValue
}
