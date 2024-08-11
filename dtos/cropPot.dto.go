package dtos

import "time"

type ControlSettingsResponse struct {
	WateringInterval int
}

// CropPotResponse represents the response DTO for a CropPot
type CropPotResponse struct {
	ID               uint
	Alias            string
	LastWateredAt    *time.Time 
	IsArchived       bool
	ControlSettings  *ControlSettingsResponse
	
}

// CreateCropPot represents the data transfer object for creating/updating CropPot
type CreateCropPot struct {
	Alias             string             `json:"alias" validate:"required,min=8"`
	ControlSettings   *ControlSettingsDto `json:"controlSettings,omitempty"`
}