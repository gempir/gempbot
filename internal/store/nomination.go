package store

import (
	"context"
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
	Votes           []NominationVote     `gorm:"foreignKey:EmoteID,ChannelTwitchID;references:EmoteID,ChannelTwitchID"`
	Downvotes       []NominationDownvote `gorm:"foreignKey:EmoteID,ChannelTwitchID;references:EmoteID,ChannelTwitchID"`
}

type NominationVote struct {
	EmoteID         string `gorm:"primarykey"`
	ChannelTwitchID string `gorm:"primarykey"`
	VoteBy          string `gorm:"primarykey"`
}

type NominationDownvote struct {
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

	res = db.Client.WithContext(ctx).Where("channel_twitch_id = ?", channelTwitchID).Delete(&NominationDownvote{})
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

	res = s.Client.WithContext(ctx).Where("channel_twitch_id = ? AND emote_id = ?", channelTwitchID, emoteID).Delete(&NominationDownvote{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *Database) GetNominations(ctx context.Context, channelTwitchID string) ([]Nomination, error) {
	var nominations []Nomination
	res := db.Client.WithContext(ctx).Preload("Votes").Preload("Downvotes").Where("channel_twitch_id = ?", channelTwitchID).Find(&nominations)
	if res.Error != nil {
		return nominations, res.Error
	}

	sort.Slice(nominations, func(i, j int) bool {
		return (len(nominations[i].Votes) - len(nominations[i].Downvotes)) > (len(nominations[j].Votes) - len(nominations[j].Downvotes))
	})

	return nominations, nil
}

func (s *Database) RemoveNomination(ctx context.Context, channelTwitchID string, emoteID string) error {
	res := s.Client.WithContext(ctx).Where("channel_twitch_id = ? AND emote_id = ?", channelTwitchID, emoteID).Delete(&Nomination{})
	if res.Error != nil {
		return res.Error
	}

	res = s.Client.WithContext(ctx).Where("channel_twitch_id = ? AND emote_id = ?", channelTwitchID, emoteID).Delete(&NominationVote{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *Database) CountNominations(ctx context.Context, channelTwitchID string, userID string) (int, error) {
	var count int64
	res := s.Client.WithContext(ctx).Model(&Nomination{}).Where("channel_twitch_id = ? AND nominated_by = ?", channelTwitchID, userID).Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}

	return int(count), nil
}

func (s *Database) IsAlreadyNominated(ctx context.Context, channelTwitchID string, emoteID string) (bool, error) {
	var count int64
	res := s.Client.WithContext(ctx).Model(&Nomination{}).Where("channel_twitch_id = ? and emote_id = ?", channelTwitchID, emoteID).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return int(count) > 0, nil
}

func (db *Database) GetNomination(ctx context.Context, channelTwitchID string, emoteID string) (Nomination, error) {
	var nomination Nomination
	res := db.Client.WithContext(ctx).Preload("Votes").Preload("Downvotes").Where("channel_twitch_id = ? AND emote_id = ?", channelTwitchID, emoteID).First(&nomination)
	if res.Error != nil {
		return nomination, res.Error
	}

	return nomination, nil
}

func (db *Database) CreateOrIncrementNomination(ctx context.Context, nomination Nomination) error {
	inputNomination := nomination
	var prevNom Nomination
	db.Client.WithContext(ctx).Preload("Votes").Preload("Downvotes").Where("emote_id = ? AND channel_twitch_id = ?", nomination.EmoteID, nomination.ChannelTwitchID).First(&prevNom)
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

func (db *Database) CreateNominationDownvote(ctx context.Context, downvote NominationDownvote) error {
	res := db.Client.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&downvote)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *Database) RemoveNominationDownvote(ctx context.Context, downvote NominationDownvote) error {
	res := db.Client.WithContext(ctx).Where("emote_id = ? AND channel_twitch_id = ? AND vote_by = ?", downvote.EmoteID, downvote.ChannelTwitchID, downvote.VoteBy).Delete(&NominationDownvote{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
