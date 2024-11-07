package models

import "gorm.io/gorm"

type PinnedCard struct {
	gorm.Model
	CanvasID      uint
	UserID        uint
	Title         string `json:"title"`
	Icon          string `json:"icon"`
	SensorID      uint   `json:"sensorId"`
	StartLocation int    `json:"startLocation"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
}
