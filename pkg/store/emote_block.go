package store

import (
	"time"

	"github.com/gempir/gempbot/pkg/dto"
	"gorm.io/gorm/clause"
)

type EmoteBlock struct {
	ChannelTwitchID string         `gorm:"primarykey"`
	Type            dto.RewardType `gorm:"primarykey"`
	EmoteID         string         `gorm:"primarykey"`
	CreatedAt       time.Time
}

func (db *Database) IsEmoteBlocked(channelTwitchID string, emoteID string, emoteType dto.RewardType) bool {
	var emoteBlocks []EmoteBlock
	db.Client.Where("channel_twitch_id = ? AND emote_id = ? AND type = ?", channelTwitchID, emoteID, emoteType).Find(&emoteBlocks)

	return len(emoteBlocks) > 0
}

func (db *Database) GetEmoteBlocks(channelTwitchID string, page int, pageSize int) []EmoteBlock {
	var emoteBlocks []EmoteBlock
	db.Client.Where("channel_twitch_id = ?", channelTwitchID).Offset((page * pageSize) - pageSize).Limit(pageSize).Order("created_at desc").Find(&emoteBlocks)

	return emoteBlocks
}

func (db *Database) BlockEmotes(channelTwitchID string, emoteIds []string, emoteType string) error {
	var emoteBlocks []EmoteBlock
	for _, emoteId := range emoteIds {
		emoteBlock := EmoteBlock{
			ChannelTwitchID: channelTwitchID,
			EmoteID:         emoteId,
			Type:            dto.RewardType(emoteType),
		}
		emoteBlocks = append(emoteBlocks, emoteBlock)
	}

	res := db.Client.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&emoteBlocks)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
