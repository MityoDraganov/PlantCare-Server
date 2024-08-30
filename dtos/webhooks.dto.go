package dtos

type AddWebhookDto struct {
	CropPotID                  uint
	EndpointUrl                string   `json:"endpointUrl"`
	Description                *string  `json:"description"`
	SubscribedEvents []SensorDto `json:"subscribedEvents"`
}

type UpdateWebhookDto struct {
	EndpointUrl                *string   `json:"endpointUrl"`
	Description                *string   `json:"description"`
	SubscribedEvents *[]SensorDto `json:"subscribedEvents"`
}

type WebhookResponse struct {
	ID               uint        `json:"id"`
	EndpointUrl      string      `json:"endpointUrl"`
	Description      *string     `json:"description"`
	SubscribedEvents []SensorDto `json:"subscribedEvents"`
}
