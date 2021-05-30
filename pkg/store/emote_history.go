package store

import (
	"errors"

	"gorm.io/gorm"
)

type EmoteAdd struct {
	gorm.Model
	ID              uint   `gorm:"primarykey,autoIncrement"`
	ChannelTwitchID string `gorm:"index"`
	EmoteID         string
}

func (db *Database) CreateEmoteAdd(channelTwitchID string, emoteID string) {
	add := EmoteAdd{ChannelTwitchID: channelTwitchID, EmoteID: emoteID}
	db.Client.Create(&add)
}

func (db *Database) GetOldestEmoteAdded(channelTwitchID string) (EmoteAdd, error) {
	var emote EmoteAdd
	result := db.Client.Where("channel_twitch_id = ?", channelTwitchID).Order("updated_at asc").First(&emote)
	if result.RowsAffected == 0 {
		return emote, errors.New("not found")
	}

	return emote, nil
}

func (db *Database) RemoveOldestEmoteAdd(channelTwitchID string) {
	db.Client.Where("channel_twitch_id = ?", channelTwitchID).Order("updated_at asc").Limit(1).Unscoped().Delete(&EmoteAdd{})
}
