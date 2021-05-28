package store

import (
	"errors"

	"gorm.io/gorm"
)

type UserAccessToken struct {
	gorm.Model
	OwnerTwitchID string `gorm:"index"`
	AccessToken   string
	RefreshToken  string
	Scopes        string
}

func (db *Database) SaveUserAccessToken(ownerId string, accessToken string, refreshToken string, scopes string) error {
	token := UserAccessToken{OwnerTwitchID: ownerId, AccessToken: accessToken, RefreshToken: refreshToken, Scopes: scopes}

	update := db.Client.Model(&UserAccessToken{}).Where("owner_twitch_id = ?", ownerId).Updates(&token)
	if update.Error != nil {
		return update.Error
	}

	if update.RowsAffected > 0 {
		updateMap := map[string]interface{}{"deleted_at": nil}
		update := db.Client.Model(&token).Where("owner_twitch_id = ?", ownerId).Updates(&updateMap)
		if update.Error != nil {
			return update.Error
		}

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
