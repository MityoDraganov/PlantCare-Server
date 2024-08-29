package dtos

import (
	"PlantCare/models"
)

type SensorDto struct {
	ID uint `json:"id"`
	SerialNumber string `json:"serialNumber"`
	Alias        string `json:"alias"`
	Description  *string `json:"description"`

	Measurements []models.Measurement `json:"measurements"`
	IsOfficial   bool `json:"isOfficial"`
}