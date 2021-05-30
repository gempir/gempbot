package store

import (
	"fmt"

	"github.com/gempir/bitraft/pkg/config"
	"github.com/gempir/bitraft/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Client *gorm.DB
}

func NewDatabase(cfg *config.Config) *Database {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=bitraft port=5432 sslmode=disable TimeZone=Europe/Berlin", cfg.PostgresUsername, cfg.PostgresPassword)
	pdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: log.NewGormLogger()})
	if err != nil {
		panic("failed to connect postgres database")
	}

	// Migrate the schema
	err = pdb.AutoMigrate(&Editor{}, &ChannelPointReward{}, &EventSubSubscription{}, &UserAccessToken{}, &EmoteAdd{})
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.SqliteDatabase), &gorm.Config{Logger: log.NewGormLogger()})
	if err != nil {
		panic("failed to connect sqlite database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Editor{}, &ChannelPointReward{}, &EventSubSubscription{}, &UserAccessToken{}, &EmoteAdd{})
	if err != nil {
		panic(err)
	}

	var editors []Editor
	db.Find(&editors)
	pdb.Create(editors)

	var rewards []ChannelPointReward
	db.Find(&rewards)
	pdb.Create(rewards)

	var subs []EventSubSubscription
	db.Find(&subs)
	pdb.Create(subs)

	var tokens []UserAccessToken
	db.Find(&tokens)
	pdb.Create(tokens)

	var adds []EmoteAdd
	db.Find(&adds)
	pdb.Create(adds)

	return &Database{
		Client: db,
	}
}
