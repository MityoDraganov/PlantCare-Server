package dtos

import (
	"PlantCare/models"
)

type SensorResponseDto struct {
	ID uint `json:"id"`
	SerialNumber string `json:"serialNumber"`
	Alias        string `json:"alias"`
	Description  *string `json:"description"`
	MeasurementInterval string `json:"measurementInterval"`

	Measurements []models.Measurement `json:"measurements"`
}

type SensorRequestDto struct {
	Alias        string `json:"alias"`
	Description  *string `json:"description"`
	MeasurementInterval string `json:"measurementInterval"`
}