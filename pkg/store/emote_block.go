package store

import (
	"time"

	"github.com/gempir/gempbot/pkg/dto"
)

type EmoteBlock struct {
	ChannelTwitchID string         `gorm:"primarykey"`
	Type            dto.RewardType `gorm:"primarykey"`
	EmoteID         string         `gorm:"primarykey"`
	CreatedAt       time.Time
}

func (db *Database) IsEmoteBlocked(channelTwitchID string, emoteID string, rewardType dto.RewardType) bool {
	var emoteBlocks []EmoteBlock
	db.Client.Where("channel_twitch_id = ? AND emote_id = ? AND type = ?", channelTwitchID, emoteID, rewardType).Find(&emoteBlocks)

	return len(emoteBlocks) > 0
}

func (db *Database) GetEmoteBlocks(channelTwitchID string, page int, pageSize int) []EmoteBlock {
	var emoteBlocks []EmoteBlock
	db.Client.Where("channel_twitch_id = ?", channelTwitchID).Offset((page * pageSize) - pageSize).Limit(pageSize).Order("created_at desc").Find(&emoteBlocks)

	return emoteBlocks
}
