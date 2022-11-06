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
	Id              uint64 `gorm:"primaryKey,autoIncrement:true"`
	Url             string
	Approved        bool
	Author          string
	Approver        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (db *Database) AddToQueue(queueItem MediaQueue) error {
	res := db.Client.Create(&queueItem)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
