package models

import (
	"gorm.io/gorm"
)

type MeasurementGroup struct {
	gorm.Model
	Measurements    []Measurement      `gorm:"foreignKey:MeasurementGroupID"`
	CropPotID       uint               `json:"cropPotId"`
	ModelOutput 	*ModelOutput `json:"modelOutput"`
}
