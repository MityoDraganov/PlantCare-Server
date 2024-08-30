package dtos

import (
	"time"
)

type ControlSettingsResponse struct {
    Id uint `json:"id"`
	WateringInterval int `json:"wateringInterval"`
}

// CropPotResponse represents the response DTO for a CropPot
type CropPotResponse struct {
    ID              uint                      `json:"id"`
    Alias           string                    `json:"alias"`
    LastWateredAt   *time.Time                `json:"lastWateredAt"`
    IsArchived      bool                      `json:"isArchived"`
    ControlSettings *ControlSettingsResponse  `json:"controlSettings"`
    Sensors      []SensorDto      `json:"sensors"`
    Webhooks []WebhookResponse `json:"webhooks"`
}

// CreateCropPot represents the data transfer object for creating/updating CropPot
type CreateCropPot struct {
	Alias             string             `json:"alias" validate:"required,min=8"`
	ControlSettings   *ControlSettingsDto `json:"controlSettings,omitempty"`
}