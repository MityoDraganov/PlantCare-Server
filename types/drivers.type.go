package types


type Sensor struct {
	SerialNumber string `json:"serialNumber"`
	Type         string `json:"type"`
	Name         string `json:"name,omitempty"` // Optional field for sensors with a name
}

type Control struct {
	SerialNumber    string `json:"serialNumber"`
	Type            string `json:"type"`
	MinValue        int    `json:"minValue"`
	MaxValue        int    `json:"maxValue"`
	DependantSensor DependantSensor `json:"dependantSensor"`
}

type Config struct {
	Sensors []Sensor `json:"sensors"`
	Controls []Control `json:"controls"`
}

type DriverConfig struct {
	SerialNumber    string          `json:"serialNumber"`
	DriverURL       string          `json:"driverUrl"`
	DependantSensorSerial DependantSensor `json:"dependantSensorSerial"`
	MinValue        int             `json:"minValue"`
	MaxValue        int             `json:"maxValue"`
}

type DriverJsonConfig struct {
	SerialNumber    string          `json:"serialNumber"`
	DriverURL       string          `json:"driverUrl"`
	DependantSensor DependantSensor `json:"dependantSensor"`
	MinValue        int             `json:"minValue"`
	MaxValue        int             `json:"maxValue"`
	Classname 	 string          `json:"classname"`
}




type DependantSensor struct {
	SerialNumber string `json:"serialNumber"`
}