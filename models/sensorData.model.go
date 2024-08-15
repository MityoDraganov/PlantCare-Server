package models

import (
	"gorm.io/gorm"
)

type SensorData struct {
    gorm.Model
    CropPotID   uint
    Temperature float32 `json:"temperature"`
    Moisture    float32 `json:"moisture"`
    WaterLevel  float32 `json:"waterLevel"`
    SunExposure float32 `json:"sunExposure"`
}

type CustomSensorField struct {
    gorm.Model
    CropPotID        uint              `json:"cropPotId"`
    FieldAlias       string            `json:"fieldAlias"`
    CustomSensorData []CustomSensorData `gorm:"foreignKey:CustomSensorFieldID"`
}

type CustomSensorData struct {
    gorm.Model
    CustomSensorFieldID uint    `json:"customSensorFieldId"`
    DataValue           float64 `json:"sensorValue"`
}
