package store

import "time"

type MediaPlayer struct {
	ChannelTwitchId string `gorm:"primaryKey"`
	CurrentTime     float32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
