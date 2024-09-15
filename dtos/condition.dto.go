package dtos

type ConditionDto struct {
	DependentSensor *SensorResponseDto `json:"dependentSensor"`
	On              float32            `json:"on"`
	Off             float32            `json:"off"`
}

type ConditionRequestDto struct {
	DependentSensor *SensorRequestDto `json:"dependentSensor"`
	On              float32           `json:"on"`
	Off             float32           `json:"off"`
}
