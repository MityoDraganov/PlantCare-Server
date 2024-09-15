package dtos

import "PlantCare/models"

type ControlDto struct {
	ID uint `json:"id"`
	SerialNumber string  `json:"serialNumber"`
	Alias        string  `json:"alias"`
	Description  *string `json:"description"`

	Updates    []models.Update `json:"updates"`
	IsOfficial bool            `json:"isOfficial"`

	Condition *ConditionDto `json:"condition"`

	ActivePeriod ActivePeriod `json:"activePeriod"`
}

type ControlRequestDto struct {
	ID uint `json:"id"`
	Alias        string  `json:"alias"`
	Description  *string `json:"description"`

	Condition ConditionRequestDto `json:"condition"`

	ActivePeriod ActivePeriod `json:"activePeriod"`
}

