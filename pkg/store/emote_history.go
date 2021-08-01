package store

import (
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

func (db *Database) GetEmoteAdded(channelTwitchID string, limit int) []EmoteAdd {
	var emotes []EmoteAdd

	db.Client.Where("channel_twitch_id = ?", channelTwitchID).Limit(limit).Order("updated_at desc").Find(&emotes)

	return emotes
}
