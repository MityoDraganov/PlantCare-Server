package dtos

type PinnedCard struct {
	Title    string `json:"title"`
	Icon     string	`json:"icon"`
	SensorID uint	`json:"sensorId"`
	Location []int `json:"location"`
}
