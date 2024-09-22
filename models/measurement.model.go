package models

import (
	"time"

	"gorm.io/gorm"
)

type Measurement struct {
	gorm.Model
	CreatedAt time.Time
	SensorID uint    `json:"sensorId"`
	Value    float32 `json:"value"`
}
