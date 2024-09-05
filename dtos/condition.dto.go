package dtos

type ConditionDto struct {
	DependentSensor *SensorResponseDto
	On float32
	Off float32
}