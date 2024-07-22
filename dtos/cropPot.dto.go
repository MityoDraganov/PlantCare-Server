package dtos

import "time"

type CreateCropPot struct {
	Alias            string `json:"alias" validate:"required,min=8"`
	WateringInterval int    `json:"wateringInterval" validate:"required"` // in minutes
}

type CropPotResponse struct {
	ID               uint      `json:"id"`
	Alias            string    `json:"alias"`
	WateringInterval int       `json:"watering_interval"`
	LastWateredAt    time.Time `json:"last_watered_at"`
	IsArchived       bool      `json:"is_archived"`
}
