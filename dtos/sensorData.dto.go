package dtos

import "time"

type SensorDataResponse struct {
    CreatedAt time.Time `json:"createdAt"`
    Temperature float32 `json:"temperature"`
    Moisture    float32 `json:"moisture"`
    WaterLevel  float32 `json:"waterLevel"`
    SunExposure float32 `json:"sunExposure"`
}

type CustomSensorDataResponse struct {
    FieldAlias string  `json:"fieldAlias"`
    DataValue  float64 `json:"sensorValue"`
}
