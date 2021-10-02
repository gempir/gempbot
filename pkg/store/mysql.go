package store

import (
	"time"

	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Client *gorm.DB
}

func NewDatabase(cfg *config.Config) *Database {
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
		panic("failed to connect database")
	}
	log.Infof("connected on %s:3306 to %s", cfg.DbHost, cfg.DbName)

	return &Database{
		Client: pdb,
	}
}

func (db *Database) Migrate() {
	log.Info("Migrating schema")
	err := db.Client.AutoMigrate(&ChannelPointReward{}, &EventSubSubscription{}, &UserAccessToken{}, &AppAccessToken{}, &EmoteAdd{}, &PredictionLog{}, &PredictionLogOutcome{}, &BotConfig{}, &Permission{}, EventSubMessage{})
	if err != nil {
		panic(err)
	}
	log.Info("Finished migrating schema")
}
