package models

import (
	"time"

	"gorm.io/gorm"
)

type ControlSettings struct {
    gorm.Model
    WateringInterval int        `json:"wateringInterval"` // in minutes
    LastWateredAt    *time.Time `json:"lastWateredAt"`

    CropPotID uint `json:"cropPotId"`
}
