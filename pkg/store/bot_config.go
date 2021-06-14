package store

type BotConfig struct {
	OwnerTwitchID string `gorm:"primaryKey"`
	Login         string
}

func (db *Database) SaveBotConfig(botCfg BotConfig) error {
	update := db.Client.Model(&botCfg).Where("owner_twitch_id = ?", botCfg.OwnerTwitchID).Updates(&botCfg)
	if update.Error != nil {
		return update.Error
	}

	if update.RowsAffected > 0 {
		return nil
	}

	update = db.Client.Create(&botCfg)

	return update.Error
}

func (db *Database) GetAllBotConfigs() []BotConfig {
	var botConfigs []BotConfig

	db.Client.Find(&botConfigs)

	return botConfigs
}
