package store

import (
	"github.com/gempir/bitraft/pkg/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Client *gorm.DB
}

func NewDatabase(sqliteDatabase string) *Database {
	db, err := gorm.Open(sqlite.Open(sqliteDatabase), &gorm.Config{Logger: log.NewGormLogger()})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Editor{}, &ChannelPointReward{}, &EventSubSubscription{}, &UserAccessToken{}, &EmoteAdd{})
	if err != nil {
		panic(err)
	}

	return &Database{
		Client: db,
	}
}
