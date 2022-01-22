package store

import (
	"errors"
	"time"
)

type EventSubMessage struct {
	ID        string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index"`
}

func (db *Database) CreateEventSubMessage(message EventSubMessage) {
	db.Client.Create(&message)
}

func (db *Database) GetEventSubMessage(id string) (EventSubMessage, error) {
	var msg EventSubMessage
	result := db.Client.Where("id = ?", id).First(&msg)
	if result.RowsAffected == 0 {
		return msg, errors.New("not found")
	}

	return msg, nil
}
