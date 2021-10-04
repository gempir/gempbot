package store

import (
	"errors"
	"time"
)

type EventSubSubscription struct {
	TargetTwitchID string `gorm:"primaryKey"`
	SubscriptionID string `gorm:"primaryKey"`
	Type           string `gorm:"index"`
	ForeignID      string `gorm:"index"`
	Version        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (db *Database) AddEventSubSubscription(targetTwitchID string, subscriptionID string, version string, subType string) {
	sub := EventSubSubscription{TargetTwitchID: targetTwitchID, SubscriptionID: subscriptionID, Version: version, Type: subType}

	db.Client.Create(&sub)
}

func (db *Database) GetEventSubSubscription(targetTwitchID string, subscriptionID string, subType string) (EventSubSubscription, error) {
	var sub EventSubSubscription
	result := db.Client.Where("target_twitch_id = ? AND subscription_id = ? AND type = ?", targetTwitchID, subscriptionID, subType).First(&sub)
	if result.RowsAffected == 0 {
		return sub, errors.New("not found")
	}

	return sub, nil
}

func (db *Database) RemoveEventSubSubscription(targetTwitchID string, subscriptionID string, subType string) {
	db.Client.Delete(&EventSubSubscription{}, "target_twitch_id = ? AND subscription_id = ? AND type = ?", targetTwitchID, subscriptionID, subType)
}
