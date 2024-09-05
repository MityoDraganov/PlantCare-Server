package models

import (
	"time"

	"gorm.io/gorm"
)

type Sensor struct {
	gorm.Model
	CropPotID          uint
	SerialNumber       string
	Alias              string        `json:"alias"`
	Description        *string       `json:"description"`
	MeasuremntInterval time.Duration `json:"measuremntInterval"`

	Measurements []Measurement
	IsOfficial   bool
}
