package models

import (
	"gorm.io/gorm"
)

type SensorData struct {
    gorm.Model
	Temperature float32 `json:"temperature"`
	Moisture    float32 `json:"moisture"`
	WaterLevel  float32 `json:"waterLevel"`
	SunExposure float32 `json:"sunExposure"`

	CropPotID uint `json:"cropPotId"`
    CropPot   CropPot `gorm:"foreignKey:CropPotID;references:ID"`
}