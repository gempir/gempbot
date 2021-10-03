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
