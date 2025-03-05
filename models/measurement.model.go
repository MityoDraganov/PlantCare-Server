package models

import (
	"PlantCare/websocket/wsTypes"
	"time"

	"gorm.io/gorm"
)

type Measurement struct {
	gorm.Model
	ID                 uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	MeasurementGroupID uint
	CreatedAt time.Time
	SensorID  uint         `json:"sensorId"`
	Value     float32      `json:"value"`
	Role      wsTypes.Role `json:"role"`
}
