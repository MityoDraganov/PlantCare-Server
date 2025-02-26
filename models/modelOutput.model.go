package models

import "gorm.io/gorm"

type ModelOutput struct {
	gorm.Model
	MeasurementGroupID uint
	PercentageHealthy uint8  `json:"percentageHealthy"`
	PlantName         string `json:"plantName"`
	CertentyPercantage uint8  `json:"certentyPercantage"`
	IsPlantRecognised bool   `json:"isPlantRecognised"`
}