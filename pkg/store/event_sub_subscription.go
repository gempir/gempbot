package store

import (
	"errors"
	"time"
)

type EventSubSubscription struct {
	TargetTwitchID string `gorm:"primaryKey"`
	SubscriptionID string `gorm:"primaryKey"`
	Version        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (db *Database) AddEventSubSubscription(targetTwitchID string, subscriptionID string, version string) {
	sub := EventSubSubscription{TargetTwitchID: targetTwitchID, SubscriptionID: subscriptionID, Version: version}

	db.Client.Create(&sub)
}

func (db *Database) GetEventSubSubscription(targetTwitchID string, subscriptionID string) (EventSubSubscription, error) {
	var sub EventSubSubscription
	result := db.Client.Where("target_twitch_id = ? AND subscription_id = ?", targetTwitchID, subscriptionID).First(&sub)
	if result.RowsAffected == 0 {
		return sub, errors.New("not found")
	}

	return sub, nil
}

func (db *Database) RemoveEventSubSubscription(targetTwitchID string, subscriptionID string) {
	db.Client.Delete(&EventSubSubscription{}, "target_twitch_id = ? AND subscription_id = ?", targetTwitchID, subscriptionID)
}
