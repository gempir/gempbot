package store

import "strings"

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

	db.Client.Delete(Editor{}, "editor_twitch_id IN (?) AND owner_user_id = ?", strings.Join(userIds, ","), ownerId)
}
