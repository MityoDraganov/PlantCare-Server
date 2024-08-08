package dtos

import "time"

type ControlSettingsResponse struct {
	WateringInterval int        `json:"wateringInterval"` // in minutes
	LastWateredAt    *time.Time `json:"lastWateredAt"`
}

// CropPotResponse represents the response DTO for a CropPot
type CropPotResponse struct {
	ID               uint                   `json:"id"`
	Alias            string                 `json:"alias"`
	WateringInterval int                    `json:"wateringInterval"`
	LastWateredAt    time.Time              `json:"lastWateredAt"`
	IsArchived       bool                   `json:"isArchived"`
	ControlSettings  *ControlSettingsResponse `json:"controlSettings,omitempty"` // Optional field
}

type ControlSettingsDTO struct {
	WateringInterval int        `json:"wateringInterval"` // in minutes
	LastWateredAt    *time.Time `json:"lastWateredAt"`
}

// CreateCropPot represents the data transfer object for creating/updating CropPot
type CreateCropPot struct {
	Alias             string             `json:"alias" validate:"required,min=8"`
	WateringInterval  int                `json:"wateringInterval" validate:"required"` // in minutes
	ControlSettings   *ControlSettingsDTO `json:"controlSettings,omitempty"` // Optional field
}