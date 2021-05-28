package store

import (
	gorm_logger "github.com/gempir/bitraft/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Client *gorm.DB
}

func NewDatabase(sqliteDatabase string) *Database {
	db, err := gorm.Open(sqlite.Open(sqliteDatabase), &gorm.Config{Logger: gorm_logger.New()})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Editor{}, &ChannelPointReward{}, &EventSubSubscription{})
	if err != nil {
		panic(err)
	}

	return &Database{
		Client: db,
	}
}
