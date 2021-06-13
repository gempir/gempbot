package store

import "time"

const (
	PREDICTION_RESOLVED = "resolved"
	PREDICTION_CANCELED = "canceled"
)

type PredictionLog struct {
	ID               string `gorm:"primaryKey"`
	OwnerTwitchID    string `gorm:"index"`
	Title            string
	WinningOutcomeID string
	Status           string
	StartedAt        time.Time
	LockedAt         time.Time
	EndedAt          time.Time
}

type PredictionLogOutcome struct {
	ID            string `gorm:"primaryKey"`
	PredictionID  string `gorm:"index"`
	Title         string
	Color         string
	Users         int
	ChannelPoints int
}

func (db *Database) SavePrediction(log PredictionLog) error {
	update := db.Client.Model(&log).Where("id = ?", log.ID).Updates(&log)
	if update.Error != nil {
		return update.Error
	}

	if update.RowsAffected > 0 {
		return nil
	}

	update = db.Client.Create(&log)

	return update.Error
}

func (db *Database) SaveOutcome(log PredictionLogOutcome) error {
	update := db.Client.Model(&log).Where("id = ?", log.ID).Updates(&log)
	if update.Error != nil {
		return update.Error
	}

	if update.RowsAffected > 0 {
		return nil
	}

	update = db.Client.Create(&log)

	return update.Error
}
