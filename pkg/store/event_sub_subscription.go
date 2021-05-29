package store

import (
	"errors"

	"gorm.io/gorm"
)

type EventSubSubscription struct {
	gorm.Model
	TargetTwitchID string `gorm:"index"`
	SubscriptionID string `gorm:"index"`
}

func (db *Database) AddEventSubSubscription(targetTwitchID string, subscriptionID string) {
	sub := EventSubSubscription{TargetTwitchID: targetTwitchID, SubscriptionID: subscriptionID}

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
