package store

type Permission struct {
	TwitchID        string `gorm:"primaryKey"`
	ChannelTwitchId string `gorm:"primaryKey"`
	Prediction      bool
}

func (db *Database) GetPermission(userID string, channelID string) Permission {
	var perm Permission

	db.Client.Where("twitch_id = ? AND channel_twitch_id = ?", userID, channelID).First(&perm)

	return perm
}

func (db *Database) SavePermission(perm Permission) error {
	updateMap, err := StructToMap(perm)
	if err != nil {
		return err
	}

	update := db.Client.Model(&perm).Where("twitch_id = ? AND channel_twitch_id = ?", perm.TwitchID, perm.ChannelTwitchId).Updates(&updateMap)
	if update.Error != nil {
		return update.Error
	}

	if update.RowsAffected > 0 {
		return nil
	}

	update = db.Client.Create(&perm)

	return update.Error
}
