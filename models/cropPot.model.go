package models

import (
	"time"

	"gorm.io/gorm"
)

type CropPot struct {
	gorm.Model
	Alias            string     `json:"alias" gorm:"not null;unique"`
	WateringInterval int        `json:"wateringInterval" gorm:"not null"` // in minutes
	LastWateredAt    *time.Time `json:"lastWateredAt"`   // timestamp of the last watering
}
