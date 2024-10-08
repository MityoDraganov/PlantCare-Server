package models

import "gorm.io/gorm"

type Webhook struct {
	gorm.Model
	CropPotID uint

	EndpointUrl      string  `json:"endpointUrl"`
	Description      *string `json:"description"`
	SubscribedEvents []Sensor `gorm:"many2many:webhook_sensors;"`
}
