package store

import (
	"context"
	"time"
)

type Nomination struct {
	EmoteID         string `gorm:"primarykey"`
	ChannelTwitchID string `gorm:"primarykey"`
	Votes           int
	EmoteCode       string
	NominatedBy     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s *Database) ClearNominations(ctx context.Context, channelTwitchID string) error {
	res := s.Client.WithContext(ctx).Where("channel_twitch_id = ?", channelTwitchID).Delete(&Nomination{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *Database) GetNominations(ctx context.Context, channelTwitchID string, page int, pageSize int) ([]Nomination, error) {
	var nominations []Nomination
	res := db.Client.WithContext(ctx).Where("channel_twitch_id = ?", channelTwitchID).Order("votes desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&nominations)
	if res.Error != nil {
		return nominations, res.Error
	}

	return nominations, nil
}

func (db *Database) GetTopVotedNominated(ctx context.Context, channelTwitchID string) (Nomination, error) {
	var nomination Nomination
	res := db.Client.WithContext(ctx).Where("channel_twitch_id = ?", channelTwitchID).Order("votes desc").First(&nomination)
	if res.Error != nil {
		return nomination, res.Error
	}

	return nomination, nil
}

func (db *Database) CreateOrIncrementNomination(ctx context.Context, nomination Nomination) error {
	var prevNom Nomination
	db.Client.WithContext(ctx).Where("emote_id = ? AND channel_twitch_id = ?", nomination.EmoteID, nomination.ChannelTwitchID).First(&prevNom)
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
