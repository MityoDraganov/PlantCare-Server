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
	SerialNumber string `gorm:"type:varchar(100);uniqueIndex;not null"`

	Driver       Driver
	Type         Type
	IsAttached   bool

	Alias              *string       `json:"alias"`
	Description        *string       `json:"description"`
	MeasuremntInterval time.Duration `json:"measuremntInterval"`

	Measurements []Measurement
	IsOfficial   bool
}
