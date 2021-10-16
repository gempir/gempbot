package store

import (
	"github.com/gempir/gempbot/pkg/dto"
	"gorm.io/gorm/clause"
)

type Permission struct {
	ChannelTwitchId string `gorm:"primaryKey"`
	TwitchID        string `gorm:"primaryKey"`
	Editor          bool   `gorm:"default:false"`
	Prediction      bool   `gorm:"default:false"`
}

func (db *Database) GetChannelUserPermissions(userID string, channelID string) Permission {
	var perm Permission

	db.Client.Where("twitch_id = ? AND channel_twitch_id = ?", userID, channelID).First(&perm)

	return perm
}

func (db *Database) GetChannelPermissions(channelID string) []Permission {
	var permissions []Permission

	db.Client.Where("channel_twitch_id = ?", channelID).Find(&permissions)

	return permissions
}

func (db *Database) GetUserPermissions(userID string) []Permission {
	var permissions []Permission

	if userID == dto.GEMPIR_USER_ID {
		db.Client.Distinct("channel_twitch_id").Find(&permissions)
	} else {
		db.Client.Where("twitch_id = ?", userID).Find(&permissions)
	}

	return permissions
}

func (db *Database) DeletePermission(channelID string, userID string) {
	db.Client.Delete(&Permission{}, "channel_twitch_id = ? AND twitch_id = ?", channelID, userID)
}

func (db *Database) SavePermission(permission Permission) error {
	update := db.Client.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&permission)

	return update.Error
}
