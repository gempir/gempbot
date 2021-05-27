package store

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Editor struct {
	gorm.Model
	OwnerTwitchID  string `gorm:"index"`
	EditorTwitchID string `gorm:"index"`
}

type Database struct {
	Client *gorm.DB
}

func NewDatabase(sqliteDatabase string) *Database {
	db, err := gorm.Open(sqlite.Open(sqliteDatabase), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Editor{}, &ChannelPointReward{})
	if err != nil {
		panic(err)
	}

	return &Database{
		Client: db,
	}
}
