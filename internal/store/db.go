package store

import (
	"context"
	"time"

	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	SaveReward(reward ChannelPointReward) error
	DeleteChannelPointRewardById(userID string, rewardID string)
	GetChannelPointReward(userID string, rewardType dto.RewardType) (ChannelPointReward, error)
}

type Database struct {
	Client     *gorm.DB
	PsqlClient *gorm.DB
}

func NewDatabase(cfg *config.Config) *Database {
	psql, err := gorm.Open(gormPostgres.Open(cfg.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect psql database " + err.Error())
	}

	// Configure connection pool
	sqlDB, err := psql.DB()
	if err != nil {
		panic("failed to get database instance " + err.Error())
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

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
	)
	if err != nil {
		panic("Failed to migrate, " + err.Error())
	}
	log.Info("Finished migrating schema")
}
