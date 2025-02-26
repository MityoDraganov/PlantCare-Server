package types


type Sensor struct {
	SerialNumber string `json:"serialNumber"`
	Type         string `json:"type"`
	Name         string `json:"name,omitempty"` // Optional field for sensors with a name
}

type Control struct {
	SerialNumber    string `json:"serialNumber"`
	Type            string `json:"type"`
	DependantSensor struct {
		SerialNumber string `json:"serialNumber"`
		MinValue     int    `json:"minValue"`
		MaxValue     int    `json:"maxValue"`
	} `json:"dependantSensor"`
}

type Config struct {
	Sensors []Sensor `json:"sensors"`
	Controls []Control `json:"controls"`
}