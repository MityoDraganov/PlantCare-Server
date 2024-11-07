package dtos

import (
	"PlantCare/models"
)

type SensorDto struct {
	ID                  uint    `json:"id"`
	SerialNumber        string  `json:"serialNumber"`
	Alias               *string  `json:"alias"`
	Description         *string `json:"description"`
	MeasurementInterval string  `json:"measurementInterval"`
	Type                models.Type
	IsAttached  bool `json:"isAttached"`


	DriverUrl           string  `json:"driverUrl"`
	Measurements    []models.Measurement `json:"measurements"`
}

type AttachSensor struct {
	SerialNumber string  `json:"serialNumber"`
	Alias        *string `json:"alias"`
	Description  *string
}

type SensorUserRequestDto struct {
	ID uint `json:"id"`
	Alias               string  `json:"alias"`
	Description         *string `json:"description"`
	MeasurementInterval string  `json:"measurementInterval"`
	DriverUrl           string  `json:"driverUrl"`
}

type SensorMeasurementsSummary struct {
	SensorType   models.Type
	Measurements []models.Measurement `json:"measurements"`
}
