package store

import (
	"time"

	"github.com/nicklaw5/helix/v2"
)

type EventSubSubscription struct {
	TargetTwitchID string `gorm:"primary_key"`
	SubscriptionID string `gorm:"primary_key;index"`
	Type           string
	ForeignID      string `gorm:"index"`
	Version        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (db *Database) AddEventSubSubscription(targetTwitchID, subscriptionID, version, subType, foreignID string) {
	sub := EventSubSubscription{TargetTwitchID: targetTwitchID, SubscriptionID: subscriptionID, Version: version, Type: subType, ForeignID: foreignID}

	db.Client.Create(&sub)
}

func (db *Database) GetAllSubscriptions() []EventSubSubscription {
	var subs []EventSubSubscription
	db.Client.Order("updated_at desc").Find(&subs)
	return subs
}

func (db *Database) GetAllPredictionSubscriptions(userID string) []EventSubSubscription {
	var subs []EventSubSubscription
	db.Client.Where("target_twitch_id = ? AND type IN (?, ?, ?)", userID, helix.EventSubTypeChannelPredictionBegin, helix.EventSubTypeChannelPredictionLock, helix.EventSubTypeChannelPredictionEnd).Find(&subs)
	return subs
}

func (db *Database) HasEventSubSubscription(subscriptionID string) bool {
	var subs []EventSubSubscription
	result := db.Client.Where("subscription_id = ?", subscriptionID).Find(&subs)
	if result.Error != nil {
		return false
	}

	return len(subs) > 0
}

func (db *Database) RemoveEventSubSubscription(subscriptionID string) {
	db.Client.Delete(&EventSubSubscription{}, "subscription_id = ?", subscriptionID)
}
