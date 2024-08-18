package models

import "gorm.io/gorm"

type Sensor struct {
	gorm.Model
	SerialNumber string
	Alias string
	Description *string
}
