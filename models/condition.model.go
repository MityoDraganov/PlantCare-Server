package models

import "gorm.io/gorm"

type Condition struct {
	gorm.Model
	ControlID uint

	DependentSensor *Sensor
	On float32
	Off float32
}