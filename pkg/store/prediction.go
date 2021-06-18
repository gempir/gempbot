package store

import (
	"errors"
	"time"

	"github.com/gempir/bitraft/pkg/dto"
)

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

func (o *PredictionLogOutcome) GetColorEmoji() string {
	if o.Color == dto.Outcome_First {
		return "🟦"
	}

	return "🟪"
}

func (db *Database) GetActivePrediction(ownerTwitchID string) (PredictionLog, error) {
	var reward PredictionLog
	result := db.Client.Where("owner_twitch_id = ? AND winning_outcome_id = ''", ownerTwitchID).Order("started_at desc").First(&reward)
	if result.RowsAffected == 0 {
		return reward, errors.New("not found")
	}

	return reward, nil
}

func (db *Database) GetOutcomes(predictionID string) []PredictionLogOutcome {
	var outcomes []PredictionLogOutcome
	db.Client.Where("prediction_id = ?", predictionID).Find(&outcomes)

	return outcomes
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
