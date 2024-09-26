package wsDtos

type SensorDTO struct {
	SerialNumber string `json:"serialNumber"`
	Alias        string `json:"alias"`
	Description  *string
}
