package store

import (
	"context"
	"errors"
	"time"
)

type UserAccessToken struct {
	OwnerTwitchID string `gorm:"primaryKey"`
	AccessToken   string
	RefreshToken  string
	Scopes        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (db *Database) SaveUserAccessToken(ctx context.Context, ownerId string, accessToken string, refreshToken string, scopes string) error {
	token := UserAccessToken{OwnerTwitchID: ownerId, AccessToken: accessToken, RefreshToken: refreshToken, Scopes: scopes}

	update := db.Client.WithContext(ctx).Model(&UserAccessToken{}).Where("owner_twitch_id = ?", ownerId).Updates(&token)
	if update.Error == nil {
		return nil
	}

	update = db.Client.Create(&token)

	return update.Error
}

func (db *Database) GetUserAccessToken(userID string) (UserAccessToken, error) {
	var token UserAccessToken
	result := db.Client.Where("owner_twitch_id = ?", userID).First(&token)
	if result.RowsAffected == 0 {
		return token, errors.New("not found")
	}

	return token, nil
}

func (db *Database) GetAllUserAccessToken() []UserAccessToken {
	var tokens []UserAccessToken
	db.Client.Find(&tokens)

	return tokens
}
