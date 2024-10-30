package models

import "gorm.io/gorm"

type PinnedCard struct {
	gorm.Model
	Title    string `json:"title"`
	Icon     string	`json:"icon"`
	SensorID uint	`json:"sensorId"`
	Location []int `json:"location"`
}
