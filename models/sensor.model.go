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

	DriverID    *uint    // Foreign key to reference the software driver
	Driver      *Driver `gorm:"foreignKey:DriverID;references:ID"` // Many sensors to one driver
	Type        Type
	IsAttached  bool

	Alias              *string       `json:"alias"`
	Description        *string       `json:"description"`
	MeasuremntInterval time.Duration `json:"measuremntInterval"`

	Measurements []Measurement
	IsOfficial   bool
}
