package store

import "time"

type MediaPlayer struct {
	ChannelTwitchId string `gorm:"primaryKey"`
	Time            float32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type MediaQueue struct {
	ChannelTwitchId string `gorm:"primaryKey"`
	Url             string
	Approved        bool
	Author          string
	Approver        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
