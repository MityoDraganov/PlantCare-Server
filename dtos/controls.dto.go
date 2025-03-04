package dtos

type ControlDto struct {
	ID           uint    `json:"id"`
	SerialNumber string  `json:"serialNumber"`
	Alias        string  `json:"alias"`
	Description  *string `json:"description"`

	IsOfficial bool `json:"isOfficial"`

	MinValue          *int   `json:"minValue"`
	MaxValue          *int   `json:"maxValue"`
	DependantSensorId *uint `json:"dependantSensor"`
	DriverUrl         string `json:"driverUrl"`
	//ActivePeriod ActivePeriod `json:"activePeriod"`
}

type AttachControlDto struct {
	SerialNumber string
}

type ControlUserRequestDto struct {
	ID                uint    `json:"id"`
	Alias             string  `json:"alias"`
	Description       *string `json:"description"`
	DriverUrl         string  `json:"driverUrl"`
	DependantSensorSerial string   `json:"dependantSensorSerial"`
	MinValue          *int    `json:"minValue"`
	MaxValue          *int    `json:"maxValue"`
	// Add other fields if necessary:
	IsOfficial        bool `json:"isOfficial"`
	IsEditing         bool `json:"isEditing"`
}
type ControlRequestDto struct {
	ID          uint    `json:"id"`
	Alias       string  `json:"alias"`
	Description *string `json:"description"`

	Condition ConditionRequestDto `json:"condition"`

	//ActivePeriod ActivePeriod `json:"activePeriod"`
}
