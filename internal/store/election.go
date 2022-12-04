package store

import (
	"context"
	"time"
)

type Election struct {
	ID                   uint   `gorm:"primarykey,autoIncrement"`
	ChannelTwitchID      string `gorm:"index"`
	Hours                int
	NominationCost       int
	ChannelPointRewardID string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	LastRunAt            *time.Time
}

func (db *Database) GetAllElections(ctx context.Context) ([]Election, error) {
	var elections []Election
	res := db.Client.WithContext(ctx).Find(&elections)
	if res.Error != nil {
		return elections, res.Error
	}

	return elections, nil
}

func (db *Database) CreateOrUpdateElection(ctx context.Context, election Election) error {
	if election.Hours < 1 {
		election.Hours = 1
	}

	res := db.Client.WithContext(ctx).Save(&election)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *Database) GetElection(ctx context.Context, channelTwitchID string) (Election, error) {
	var election Election
	res := db.Client.WithContext(ctx).Where("channel_twitch_id = ?", channelTwitchID).First(&election)
	if res.Error != nil {
		return election, res.Error
	}

	return election, nil
}

func (db *Database) DeleteElection(ctx context.Context, channelTwitchID string) error {
	res := db.Client.WithContext(ctx).Where("channel_twitch_id = ?", channelTwitchID).Delete(&Election{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
