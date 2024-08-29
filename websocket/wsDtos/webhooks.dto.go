package wsDtos

import (
	"PlantCare/dtos"
	"PlantCare/models"
)

type WebhookPayload struct {
	Sensor dtos.SensorDto
	Measurement models.Measurement
}