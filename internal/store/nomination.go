package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/gempir/gempbot/internal/log"
	"gorm.io/gorm/clause"
)

type Nomination struct {
	EmoteID         string `gorm:"primarykey"`
	ChannelTwitchID string `gorm:"primarykey"`
	EmoteCode       string
	NominatedBy     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Votes           []NominationVote `gorm:"foreignKey:EmoteID,ChannelTwitchID;references:EmoteID,ChannelTwitchID"`
}

type NominationVote struct {
	EmoteID         string `gorm:"primarykey"`
	ChannelTwitchID string `gorm:"primarykey"`
	VoteBy          string `gorm:"primarykey"`
}

func (db *Database) ClearNominations(ctx context.Context, channelTwitchID string) error {
	res := db.Client.WithContext(ctx).Where("channel_twitch_id = ?", channelTwitchID).Delete(&Nomination{})
	if res.Error != nil {
		log.Error(res.Error)
	}

	res = db.Client.WithContext(ctx).Where("channel_twitch_id = ?", channelTwitchID).Delete(&NominationVote{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *Database) ClearNominationEmote(ctx context.Context, channelTwitchID string, emoteID string) error {
	res := s.Client.WithContext(ctx).Where("channel_twitch_id = ? AND emote_id = ?", channelTwitchID, emoteID).Delete(&Nomination{})
	if res.Error != nil {
		log.Error(res.Error)
	}

	res = s.Client.WithContext(ctx).Where("channel_twitch_id = ? AND emote_id = ?", channelTwitchID, emoteID).Delete(&NominationVote{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *Database) CreateNominationVote(ctx context.Context, vote NominationVote) error {
	res := db.Client.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&vote)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *Database) RemoveNominationVote(ctx context.Context, vote NominationVote) error {
	res := db.Client.WithContext(ctx).Where("emote_id = ? AND channel_twitch_id = ? AND vote_by = ?", vote.EmoteID, vote.ChannelTwitchID, vote.VoteBy).Delete(&NominationVote{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *Database) GetNominations(ctx context.Context, channelTwitchID string) ([]Nomination, error) {
	var nominations []Nomination
	res := db.Client.WithContext(ctx).Preload("Votes").Where("channel_twitch_id = ?", channelTwitchID).Find(&nominations)
	if res.Error != nil {
		return nominations, res.Error
	}

	sort.Slice(nominations, func(i, j int) bool {
		return len(nominations[i].Votes) > len(nominations[j].Votes)
	})

	return nominations, nil
}

func (db *Database) GetTopVotedNominated(ctx context.Context, channelTwitchID string, count int) ([]Nomination, error) {
	var votes []NominationVote
	db.Client.WithContext(ctx).Raw("SELECT emote_id, COUNT(*) FROM nomination_votes WHERE channel_twitch_id = ? GROUP BY (emote_id) ORDER BY COUNT(*) DESC LIMIT ?", channelTwitchID, count).Scan(&votes)

	if len(votes) == 0 {
		return []Nomination{}, fmt.Errorf("no votes found for channel %s", channelTwitchID)
	}

	var nominations []Nomination
	for _, vote := range votes {
		var nom Nomination
		res := db.Client.WithContext(ctx).Preload("Votes").Where("channel_twitch_id = ? AND emote_id = ?", channelTwitchID, vote.EmoteID).First(&nom)
		if res.Error != nil {
			continue
		}
		nominations = append(nominations, nom)
	}

	return nominations, nil
}

func (db *Database) CreateOrIncrementNomination(ctx context.Context, nomination Nomination) error {
	inputNomination := nomination
	var prevNom Nomination
	db.Client.WithContext(ctx).Preload("Votes").Where("emote_id = ? AND channel_twitch_id = ?", nomination.EmoteID, nomination.ChannelTwitchID).First(&prevNom)
	if len(prevNom.Votes) > 0 {
		log.Infof("incrementing nomination %s", nomination.NominatedBy)
		nomination = prevNom
		nomination.CreatedAt = time.Now()
		nomination.Votes = append(nomination.Votes, NominationVote{EmoteID: nomination.EmoteID, ChannelTwitchID: nomination.ChannelTwitchID, VoteBy: inputNomination.NominatedBy})
	} else {
		nomination.Votes = []NominationVote{{EmoteID: nomination.EmoteID, ChannelTwitchID: nomination.ChannelTwitchID, VoteBy: inputNomination.NominatedBy}}
	}

	res := db.Client.WithContext(ctx).Save(&nomination)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
