package store

import (
	"context"

	"github.com/gempir/gempbot/pkg/dto"
	"gorm.io/gorm"
)

type EmoteAdd struct {
	gorm.Model
	ID              uint                `gorm:"primarykey,autoIncrement"`
	ChannelTwitchID string              `gorm:"index"`
	Type            dto.RewardType      `gorm:"index"`
	ChangeType      dto.EmoteChangeType `gorm:"index"`
	EmoteID         string
}

func (db *Database) CreateEmoteAdd(channelTwitchID string, addType dto.RewardType, emoteID string, emoteChangeType dto.EmoteChangeType) {
	add := EmoteAdd{ChannelTwitchID: channelTwitchID, Type: addType, EmoteID: emoteID, ChangeType: emoteChangeType}
	db.Client.Create(&add)
}

func (db *Database) GetEmoteAdded(channelTwitchID string, addType dto.RewardType, limit int) []EmoteAdd {
	var emotes []EmoteAdd

	db.Client.Where("channel_twitch_id = ? AND type = ? AND change_type = ?", channelTwitchID, addType, dto.EMOTE_ADD_ADD).Order("updated_at desc").Limit(limit).Find(&emotes)

	return emotes
}

func (db *Database) GetEmoteHistory(ctx context.Context, ownerTwitchID string, page int, pageSize int) []EmoteAdd {
	var emoteHistory []EmoteAdd
	db.Client.WithContext(ctx).Where("channel_twitch_id = ?", ownerTwitchID).Offset((page * pageSize) - pageSize).Limit(pageSize).Order("updated_at desc").Find(&emoteHistory)

	return emoteHistory
}
