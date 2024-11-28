package wsDtos

type SensorDTO struct {
	SerialNumber string `json:"serialNumber"`
	Alias        string `json:"alias"`
	Description  *string
}

type UserMeasurementResponse struct {
	SensorID uint    `json:"sensorId"`
	Value    float32 `json:"value"`
}