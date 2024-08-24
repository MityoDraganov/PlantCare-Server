package models

import (
	"gorm.io/gorm"
)

type Measurement struct {
	gorm.Model
	SensorID uint    `json:"sensorId"`
	Value    float32 `json:"value"`
}
