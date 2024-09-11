package models

import "gorm.io/gorm"

type Condition struct {
	gorm.Model
	ControlID uint

	DependentSensorID *uint
	DependentSensor   *Sensor    `gorm:"foreignKey:DependentSensorID"`
	On float32
	Off float32
}