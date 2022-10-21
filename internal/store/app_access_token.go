package store

import (
	"context"
	"errors"
	"time"

	"github.com/gempir/gempbot/internal/log"
)

type AppAccessToken struct {
	AccessToken  string `gorm:"primaryKey"`
	RefreshToken string
	Scopes       string
	ExpiresIn    int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (db *Database) SaveAppAccessToken(ctx context.Context, accessToken string, refreshToken string, scopes string, expiresIn int) error {
	token := AppAccessToken{AccessToken: accessToken, RefreshToken: refreshToken, Scopes: scopes, ExpiresIn: expiresIn}

	update := db.Client.Create(&token)
	if update.Error != nil {
		return update.Error
	}

	update = db.Client.Where("updated_at < DATETIME('now', '-30 day')").Delete(&AppAccessToken{})
	if update.RowsAffected > 0 {
		log.Infof("Deleted %d old app access tokens", update.RowsAffected)
	}

	return update.Error
}

func (db *Database) GetAppAccessToken() (AppAccessToken, error) {
	var token AppAccessToken
	result := db.Client.Order("updated_at desc").First(&token)
	if result.RowsAffected == 0 {
		return token, errors.New("not found")
	}

	return token, nil
}
