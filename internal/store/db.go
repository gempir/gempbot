package store

import (
	"context"
	"time"

	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
	"github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
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
}

type Database struct {
	Client     *gorm.DB
	PsqlClient *gorm.DB
}

func NewDatabase(cfg *config.Config) *Database {
	if cfg.DbHost == "" {
		panic("No database host specified")
	}

	mysqlConfig := mysql.Config{
		User:                 cfg.DbUsername,
		Passwd:               cfg.DbPassword,
		Addr:                 cfg.DbHost + ":3306",
		Net:                  "tcp",
		DBName:               cfg.DbName,
		Loc:                  time.Local,
		ParseTime:            true,
		AllowNativePasswords: true,
		TLSConfig:            "true",
	}

	pdb, err := gorm.Open(gormMysql.Open(mysqlConfig.FormatDSN()), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect mysql database " + err.Error())
	}
	log.Infof("connected on %s:3306 to %s", cfg.DbHost, cfg.DbName)

	psql, err := gorm.Open(gormPostgres.Open(cfg.DSN), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect psql database " + err.Error())
	}

	return &Database{
		Client:     pdb,
		PsqlClient: psql,
	}
}

func (db *Database) Migrate() {
	log.Info("Migrating schema")
	err := db.Client.AutoMigrate(SystemConfig{}, ChannelPointReward{}, EventSubSubscription{}, UserAccessToken{}, AppAccessToken{}, EmoteAdd{}, BotConfig{}, Permission{}, EventSubMessage{}, EmoteBlock{})
	if err != nil {
		panic("Failed to migrate, " + err.Error())
	}
	log.Info("Finished migrating schema")
}

func (db *Database) MigratePsql() {
	log.Info("Migrating schema")
	err := db.PsqlClient.AutoMigrate(SystemConfig{}, ChannelPointReward{}, EventSubSubscription{}, UserAccessToken{}, AppAccessToken{}, EmoteAdd{}, BotConfig{}, Permission{}, EventSubMessage{}, EmoteBlock{})
	if err != nil {
		panic("Failed to migrate, " + err.Error())
	}
	log.Info("Finished migrating schema")
}
