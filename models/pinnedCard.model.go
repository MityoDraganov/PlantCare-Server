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
	width         int    `json:"width"`
	height        int    `json:"height"`
}
