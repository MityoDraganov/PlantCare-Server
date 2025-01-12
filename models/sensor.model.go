package models

import (
	"gorm.io/gorm"
)

type Type string

const (
	SoilTempType Type = "soil_temp"
	SoilMoistureType Type = "soil_moisture"
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

	Measurements []Measurement
	IsOfficial   bool
}
