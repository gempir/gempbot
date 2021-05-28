package store

import (
	"strings"

	"gorm.io/gorm"
)

type Editor struct {
	gorm.Model
	OwnerTwitchID  string `gorm:"index"`
	EditorTwitchID string `gorm:"index"`
}

func (db *Database) AddEditors(ownerId string, userIds []string) {
	if len(userIds) == 0 {
		return
	}

	var editors []Editor
	for _, id := range userIds {
		editors = append(editors, Editor{OwnerTwitchID: ownerId, EditorTwitchID: id})
	}

	db.Client.Create(&editors)
}

func (db *Database) RemoveEditors(ownerId string, userIds []string) {
	if len(userIds) == 0 {
		return
	}

	db.Client.Delete(Editor{}, "editor_twitch_id IN (?) AND owner_twitch_id = ?", strings.Join(userIds, ","), ownerId)
}

func (db *Database) GetEditors(userID string) []Editor {
	var editors []Editor
	db.Client.Where("owner_twitch_id = ? OR editor_twitch_id = ?", userID, userID).Find(&editors)

	return editors
}
