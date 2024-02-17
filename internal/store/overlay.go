package store

import "gorm.io/gorm/clause"

type Overlay struct {
	ID            string `gorm:"primaryKey"`
	OwnerTwitchID string `gorm:"index"`
	RoomID        string `gorm:"index"`
}

func (db *Database) GetOverlays(userID string) []Overlay {
	var overlays []Overlay

	db.Client.Where("owner_twitch_id = ?", userID).Find(&overlays)

	return overlays
}

func (db *Database) GetOverlay(ID string, userID string) Overlay {
	var overlay Overlay

	db.Client.Where("id = ? AND owner_twitch_id = ?", ID, userID).First(&overlay)

	return overlay
}

func (db *Database) GetOverlayByRoomId(roomID string) Overlay {
	var overlay Overlay

	db.Client.Where("room_id = ?", roomID).First(&overlay)

	return overlay
}

func (db *Database) DeleteOverlay(ID string) {
	db.Client.Delete(&Overlay{}, "id = ?", ID)
}

func (db *Database) SaveOverlay(overlay Overlay) error {
	update := db.Client.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&overlay)

	return update.Error
}
