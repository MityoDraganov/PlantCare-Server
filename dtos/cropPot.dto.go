package dtos

type CropPotResponse struct {
	ID            uint              `json:"id"`
	Alias         string            `json:"alias"`
	IsArchived    bool              `json:"isArchived"`
	Controls      []ControlDto     `json:"controls"`
	Sensors       []SensorResponseDto       `json:"sensors"`
	Webhooks      []WebhookResponse `json:"webhooks"`
}

type CropPotRequest struct {
	Alias           string `json:"alias" validate:"required,min=8"`
}

type CreateCropPot struct {
	Alias           string `json:"alias" validate:"required,min=8"`
	ControlSettings *[]ControlDto
}
