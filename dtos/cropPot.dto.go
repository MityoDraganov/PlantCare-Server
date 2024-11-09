package dtos

import "PlantCare/models"

type CropPotResponse struct {
	ID         uint              `json:"id"`
	Alias      string            `json:"alias"`
	IsArchived bool              `json:"isArchived"`
	IsPinned   bool              `json:"isPinned"`
	Controls   []ControlDto      `json:"controls"`
	Sensors    []SensorDto       `json:"sensors"`
	Webhooks   []WebhookResponse `json:"webhooks"`
	Status     models.Status     `json:"status"`
	MeasurementInterval string  `json:"measurementInterval"`
	Canvas CanvasDto `json:"canvas"`
}

type CropPotRequest struct {
	Alias string `json:"alias" validate:"required,min=8"`
	IsPinned bool `json:"isPinned"`
	MeasurementInterval string  `json:"measurementInterval"`
}

type CreateCropPot struct {
	Alias           string `json:"alias" validate:"required,min=8"`
	ControlSettings *[]ControlDto
}
