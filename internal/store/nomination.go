package store

import (
	"context"
	"time"
)

type Nomination struct {
	EmoteID         string `gorm:"primarykey"`
	ChannelTwitchID string `gorm:"primarykey"`
	ElectionID      uint   `gorm:"primarykey"`
	Votes           int
	EmoteCode       string
	NominatedBy     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (db *Database) CreateOrIncrementNomination(ctx context.Context, nomination Nomination) error {
	var prevNom Nomination
	db.Client.WithContext(ctx).Where("emote_id = ? AND channel_twitch_id = ? AND election_id = ?", nomination.EmoteID, nomination.ChannelTwitchID, nomination.ElectionID).First(&prevNom)
	if prevNom.Votes != 0 {
		nomination = prevNom
		nomination.Votes = prevNom.Votes + 1
	} else {
		nomination.Votes = 1
	}

	res := db.Client.WithContext(ctx).Save(&nomination)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
