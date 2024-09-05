package wsDtos

import (
	"PlantCare/dtos"
	"PlantCare/models"
)

type WebhookPayload struct {
	Sensor dtos.SensorResponseDto
	Measurement models.Measurement
}