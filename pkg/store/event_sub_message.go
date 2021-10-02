package store

import "time"

type EventSubMessage struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
}
