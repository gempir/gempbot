package store

import (
	"errors"

	"gorm.io/gorm"
)

type EventSubSubscription struct {
	gorm.Model
	OwnerTwitchID  string `gorm:"index"`
	SubscriptionID string `gorm:"index"`
}

func (db *Database) AddEventSubSubscription(ownerId string, subscriptionID string) {
	sub := EventSubSubscription{OwnerTwitchID: ownerId, SubscriptionID: subscriptionID}

	db.Client.Create(&sub)
}

func (db *Database) GetEventSubSubscription(ownerId string, subscriptionID string) (EventSubSubscription, error) {
	var sub EventSubSubscription
	result := db.Client.Where("owner_twitch_id = ? AND subscription_id = ?", ownerId, subscriptionID).First(&sub)
	if result.RowsAffected == 0 {
		return sub, errors.New("not found")
	}

	return sub, nil
}

func (db *Database) RemoveEventSubSubscription(ownerId string, subscriptionID string) {
	db.Client.Delete(Editor{}, "owner_twitch_id = ? AND subscription_id = ?", ownerId, subscriptionID)
}
