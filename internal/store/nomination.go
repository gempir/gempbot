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

func (s *Database) ClearNominations(ctx context.Context, channelTwitchID string, electionID uint) error {
	res := s.Client.WithContext(ctx).Where("channel_twitch_id = ? AND election_id = ?", channelTwitchID, electionID).Delete(&Nomination{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *Database) GetNominations(ctx context.Context, channelTwitchID string, electionID uint, page int, pageSize int) ([]Nomination, error) {
	var nominations []Nomination
	res := db.Client.WithContext(ctx).Where("channel_twitch_id = ? AND election_id = ?", channelTwitchID, electionID).Order("votes desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&nominations)
	if res.Error != nil {
		return nominations, res.Error
	}

	return nominations, nil
}

func (db *Database) GetTopVotedNominated(ctx context.Context, channelTwitchID string, electionID uint) (Nomination, error) {
	var nomination Nomination
	res := db.Client.WithContext(ctx).Where("channel_twitch_id = ? AND election_id = ?", channelTwitchID, electionID).Order("votes desc").First(&nomination)
	if res.Error != nil {
		return nomination, res.Error
	}

	return nomination, nil
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
