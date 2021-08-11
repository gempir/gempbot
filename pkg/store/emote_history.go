package store

import (
	"gorm.io/gorm"
)

type EmoteAdd struct {
	gorm.Model
	ID              uint   `gorm:"primarykey,autoIncrement"`
	ChannelTwitchID string `gorm:"index"`
	Type            string `gorm:"index"`
	EmoteID         string
}

func (db *Database) CreateEmoteAdd(channelTwitchID string, addType string, emoteID string) {
	add := EmoteAdd{ChannelTwitchID: channelTwitchID, Type: addType, EmoteID: emoteID}
	db.Client.Create(&add)
}

func (db *Database) GetEmoteAdded(channelTwitchID string, addType string, limit int) []EmoteAdd {
	var emotes []EmoteAdd

	db.Client.Where("channel_twitch_id = ? AND type = ?", channelTwitchID, addType).Order("updated_at desc").Limit(limit).Find(&emotes)

	return emotes
}
