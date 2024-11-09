package wsDtos


type SensorMeasuremntDto struct {
	SensorSerialNumber string `json:"sensorSerialNumber"`
	Value float32 `json:"value"`
}
