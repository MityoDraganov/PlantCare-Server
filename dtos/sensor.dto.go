package dtos

import (
	"PlantCare/models"
)

type SensorDto struct {
	ID uint `json:"id"`
	SerialNumber string `json:"serialNumber"`
	Alias        string `json:"alias"`
	Description  *string `json:"description"`
	MeasuremntInterval string `json:"measuremntInterval"`

	Measurements []models.Measurement `json:"measurements"`
}