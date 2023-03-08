package store

import (
	"time"

	"github.com/teris-io/shortid"
)

type MediaPlayer struct {
	ChannelTwitchId string `gorm:"primaryKey"`
	Time            float32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type MediaQueue struct {
	ID              string `gorm:"primaryKey"`
	ChannelTwitchId string `gorm:"index"`
	Url             string
	Approved        bool
	Author          string
	Approver        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (db *Database) AddToQueue(queueItem MediaQueue) error {
	var err error
	queueItem.ID, err = shortid.Generate()
	if err != nil {
		return err
	}

	res := db.Client.Create(&queueItem)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *Database) GetQueue(channelID string) []MediaQueue {
	var queue []MediaQueue

	db.Client.Where("channel_twitch_id = ?", channelID).Find(&queue)

	return queue
}
