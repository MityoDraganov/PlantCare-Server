package dtos

type AddWebhookDto struct {
	CropPotID uint

	EndpointUrl      string  `json:"endpointUrl"`
	Description      *string `json:"description"`
	SubscribedEventsSerialNums []string  `json:"subscribedEventsSerialNums"`
}

type WebhookResponse struct {
	ID uint `json:"id"`
	EndpointUrl      string  `json:"endpointUrl"`
	Description      *string `json:"description"`
	SubscribedEvents []SensorDto `json:"subscribedEvents"`
}

