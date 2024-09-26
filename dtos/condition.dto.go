package dtos

type ConditionDto struct {
	DependentSensor *SensorDto `json:"dependentSensor"`
	On              float32            `json:"on"`
	Off             float32            `json:"off"`
}

type ConditionRequestDto struct {
	DependentSensor *SensorDto `json:"dependentSensor"`
	On              float32           `json:"on"`
	Off             float32           `json:"off"`
}
