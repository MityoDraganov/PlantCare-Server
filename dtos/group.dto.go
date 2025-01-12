package dtos

import "PlantCare/models"



type GroupRequestDto struct {
	CropPotID uint `json:"cropPotId"`
}

type GroupResponsetDto struct {
	CropPots []CropPotResponse
}



type MeasurementGroupDto struct {
	MeasurementGroupID uint `json:"measurementGroupId"`
	Measurements      []models.Measurement `gorm:"foreignKey:MeasurementGroupID" json:"measurements"`
	CropPotID         uint 		`json:"cropPotId"`
	HealthStatus      *float32		`json:"healthStatus"`
}

