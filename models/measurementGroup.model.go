package models

import "gorm.io/gorm"

type MeasurementGroup struct {
	gorm.Model
	Measurements      []Measurement `gorm:"foreignKey:MeasurementGroupID"`
	CropPotID         uint 		`json:"cropPotId"`
	HealthStatus      *float32		`json:"healthStatus"`
}
