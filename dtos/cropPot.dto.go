package dtos

import (
	"time"
)

type ControlSettingsResponse struct {
	WateringInterval int
}

// CropPotResponse represents the response DTO for a CropPot
type CropPotResponse struct {
    ID              uint                      `json:"id"`
    Alias           string                    `json:"alias"`
    LastWateredAt   *time.Time                `json:"lastWateredAt"`
    IsArchived      bool                      `json:"isArchived"`
    ControlSettings *ControlSettingsResponse  `json:"controlSettings"`
    SensorData      []SensorDataResponse      `json:"sensorData"`
}

// CreateCropPot represents the data transfer object for creating/updating CropPot
type CreateCropPot struct {
	Alias             string             `json:"alias" validate:"required,min=8"`
	ControlSettings   *ControlSettingsDto `json:"controlSettings,omitempty"`
}