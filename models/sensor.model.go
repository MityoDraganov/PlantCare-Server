package models

import (
	"time"

	"gorm.io/gorm"
)

type Sensor struct {
	gorm.Model
	CropPotID uint
	SerialNumber string
	Alias string
	Description *string
	MeasuremntInterval time.Duration

	Measurements []Measurement
	IsOfficial bool
}
