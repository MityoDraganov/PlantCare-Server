package models

import "gorm.io/gorm"

type Sensor struct {
	gorm.Model
	CropPotID uint
	SerialNumber string
	Alias string
	Description *string

	Measurements []Measurement
	IsOfficial bool
}
