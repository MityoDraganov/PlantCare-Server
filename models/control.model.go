package models

import (
	"gorm.io/gorm"
)

type Control struct {
	gorm.Model
	CropPotID    uint
	SerialNumber string
	Alias        string
	Description  *string
	IsOfficial bool
	IsAttached bool
	DriverUrl string

	
	DriverID    *uint    // Foreign key to reference the software driver
	Driver      *Driver `gorm:"foreignKey:DriverID;references:ID"` // Many sensors to one driver
	DependantSensorSerial *string
	MinValue          *int
	MaxValue          *int
}
