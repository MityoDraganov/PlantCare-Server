package models

import (
	"gorm.io/gorm"
)

type ControlSettings struct {
    gorm.Model
    WateringInterval int        `json:"wateringInterval"` // in minutes


    CropPotID uint `json:"cropPotId"`
}
