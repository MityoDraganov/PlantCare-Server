package dtos

import "PlantCare/models"

type AddWebhookDto struct {
	CropPotID uint

	EndpointUrl      string  `json:"endpointUrl"`
	Description      *string `json:"description"`
	SubscribedEvents []models.Sensor  `json:"subscribedEvents"`
}
