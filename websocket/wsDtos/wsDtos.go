package wsDtos


type SensorDataDto struct {
	Temperature float32 `json:"temperature"`
	Moisture    float32 `json:"moisture"`
	WaterLevel  float32 `json:"waterLevel"`
	SunExposure float32 `json:"sunExposure"`
}