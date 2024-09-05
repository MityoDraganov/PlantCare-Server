package dtos

type AddWebhookDto struct {
	CropPotID                  uint
	EndpointUrl                string   `json:"endpointUrl"`
	Description                *string  `json:"description"`
	SubscribedEvents []SensorResponseDto `json:"subscribedEvents"`
}

type UpdateWebhookDto struct {
	EndpointUrl                *string   `json:"endpointUrl"`
	Description                *string   `json:"description"`
	SubscribedEvents *[]SensorResponseDto `json:"subscribedEvents"`
}

type WebhookResponse struct {
	ID               uint        `json:"id"`
	EndpointUrl      string      `json:"endpointUrl"`
	Description      *string     `json:"description"`
	SubscribedEvents []SensorResponseDto `json:"subscribedEvents"`
}
