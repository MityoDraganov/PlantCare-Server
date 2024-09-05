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

	Updates    []Update

	Condition *Condition
	ActivePeriod ActivePeriod
}
