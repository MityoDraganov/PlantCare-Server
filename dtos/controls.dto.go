package dtos

import "PlantCare/models"

type ControlDto struct {
	ID uint `json:"id"`
	SerialNumber string  `json:"serialNumber"`
	Alias        string  `json:"alias"`
	Description  *string `json:"description"`

	Updates    []models.Update `json:"updates"`
	IsOfficial bool            `json:"isOfficial"`

	OnCondition  float32 `json:"onCondition"`
	OffCondition float32 `json:"offCondition"`

	ActivePeriod ActivePeriod `json:"activePeriod"`
}

type ControlRequestDto struct {
	ID uint `json:"id"`
	SerialNumber string  `json:"serialNumber"`
	Alias        string  `json:"alias"`
	Description  *string `json:"description"`

	OnCondition  float32 `json:"onCondition"`
	OffCondition float32 `json:"offCondition"`

	ActivePeriod models.ActivePeriod `json:"activePeriod"`
}

