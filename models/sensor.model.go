package models

import (
	"time"

	"gorm.io/gorm"
)

type Type string

const (
	TempType Type = "temperature"
)

type Sensor struct {
	gorm.Model
	CropPotID    uint
	SerialNumber string
	Driver       Driver
	Type         Type

	Alias              string        `json:"alias"`
	Description        *string       `json:"description"`
	MeasuremntInterval time.Duration `json:"measuremntInterval"`

	Measurements []Measurement
	IsOfficial   bool
}
