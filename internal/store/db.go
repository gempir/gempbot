package store

import (
	"context"

	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Store interface {
	IsEmoteBlocked(channelUserID string, emoteID string, rewardType dto.RewardType) bool
	GetEmoteAdded(channelUserID string, rewardType dto.RewardType, slots int) []EmoteAdd
	CreateEmoteAdd(channelUserId string, rewardType dto.RewardType, emoteID string, changeType dto.EmoteChangeType)
	GetUserAccessToken(userID string) (UserAccessToken, error)
	GetAppAccessToken() (AppAccessToken, error)
	SaveAppAccessToken(ctx context.Context, accessToken string, refreshToken string, scopes string, expiresIn int) error
	SaveUserAccessToken(ctx context.Context, ownerId string, accessToken string, refreshToken string, scopes string) error
	GetAllUserAccessToken() []UserAccessToken
	GetSevenTvToken(ctx context.Context) string
	GetBttvToken(ctx context.Context) string
	CreateOrUpdateElection(ctx context.Context, election Election) error
	GetElection(ctx context.Context, channelTwitchID string) (Election, error)
	DeleteElection(ctx context.Context, channelTwitchID string) error
	GetAllElections(ctx context.Context) ([]Election, error)
	SaveReward(reward ChannelPointReward) error
	CreateOrIncrementNomination(ctx context.Context, nomination Nomination) error
	GetTopVotedNominated(ctx context.Context, channelTwitchID string) (Nomination, error)
	GetNominations(ctx context.Context, channelTwitchID string) ([]Nomination, error)
	GetActiveElection(ctx context.Context, channelTwitchID string) (Election, error)
	ClearNominations(ctx context.Context, channelTwitchID string) error
	ClearNominationEmote(ctx context.Context, channelTwitchID string, emoteID string) error
	DeleteChannelPointRewardById(userID string, rewardID string)
	GetChannelPointReward(userID string, rewardType dto.RewardType) (ChannelPointReward, error)
	CreateNominationVote(ctx context.Context, vote NominationVote) error
}

type Database struct {
	Client     *gorm.DB
	PsqlClient *gorm.DB
}

func NewDatabase(cfg *config.Config) *Database {
	psql, err := gorm.Open(gormPostgres.Open(cfg.DSN), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect psql database " + err.Error())
	}

	log.Infof("connected on postgres")

	return &Database{
		Client: psql,
	}
}

func (db *Database) Migrate() {
	log.Info("Migrating schema")
	err := db.Client.AutoMigrate(
		SystemConfig{},
		ChannelPointReward{},
		EventSubSubscription{},
		UserAccessToken{},
		AppAccessToken{},
		EmoteAdd{},
		BotConfig{},
		Permission{},
		EventSubMessage{},
		EmoteBlock{},
		MediaPlayer{},
		MediaQueue{},
		Election{},
		Nomination{},
		NominationVote{},
	)
	if err != nil {
		panic("Failed to migrate, " + err.Error())
	}
	log.Info("Finished migrating schema")
}
