package dtos

type CropPotResponse struct {
	ID            uint              `json:"id"`
	Alias         string            `json:"alias"`
	IsArchived    bool              `json:"isArchived"`
	Controls      []ControlDto     `json:"controls"`
	Sensors       []SensorDto       `json:"sensors"`
	Webhooks      []WebhookResponse `json:"webhooks"`
}

type CreateCropPot struct {
	Alias           string `json:"alias" validate:"required,min=8"`
	ControlSettings *[]ControlDto
}
