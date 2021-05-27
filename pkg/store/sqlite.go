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

func NewDatabase(sqliteDatabase string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(sqliteDatabase), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Editor{})
	if err != nil {
		panic(err)
	}

	return db
}



