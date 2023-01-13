package store

import (
	"context"
	"time"

	"github.com/gempir/gempbot/internal/dto"
)

type EmoteLog struct {
	CreatedAt       time.Time `gorm:"primarykey"`
	EmoteID         string    `gorm:"primarykey"`
	EmoteCode       string
	AddedBy         string
	ChannelTwitchID string
	Type            dto.RewardType
}

func (db *Database) AddEmoteLogEntry(ctx context.Context, emoteLog EmoteLog) {
	db.Client.Create(&emoteLog)
}

func (db *Database) GetEmoteLogEntries(ctx context.Context, channelTwitchID string, limit int, page int) []EmoteLog {
	if limit > 100 {
		limit = 100
	}

	var emoteLogs []EmoteLog
	db.Client.Where("channel_twitch_id = ?", channelTwitchID).Order("created_at desc").Offset((page * limit) - limit).Limit(limit).Find(&emoteLogs)
	return emoteLogs
}
