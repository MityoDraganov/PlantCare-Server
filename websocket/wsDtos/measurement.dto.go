package wsDtos


type SensorMeasuremntDto struct {
	SensorSerialNumber string `json:"sensorSerialNumber"`
	Value float32 `json:"value"`
}

type MlDataDto struct {
	CropType string `json:"cropType"`
	Ph float32 `json:"ph"`
	Temperature float32 `json:"temperature"`
	SoilMoisture float32 `json:"humidity"`
}