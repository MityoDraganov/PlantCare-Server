package dtos

import "PlantCare/models"

type AddWebhookDto struct {
	CropPotID uint

	EndpointUrl      string  `json:"endpointUrl"`
	Description      *string `json:"description"`
	SubscribedEventsSerialNums []string  `json:"subscribedEventsSerialNums"`
}

type WebhookResponse struct {
	Sensor SensorWebhookResponse

	Measurement models.Measurement
	IsOfficial bool
}
