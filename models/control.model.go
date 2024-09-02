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
	OnCondition float32
	OffCondition float32

	ActivePeriod *ActivePeriod `gorm:"foreignKey:ControlID"`
}
