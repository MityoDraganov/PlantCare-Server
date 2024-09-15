package models

import (
	"gorm.io/gorm"
)

type CropPot struct {
	gorm.Model
	GroupID     *uint
	IsPinned    bool
	IsArchived  bool
	Token       string `gorm:"size:255;uniqueIndex;not null"`
	User        User   `gorm:"foreignKey:ClerkUserID;references:ClerkID"`
	ClerkUserID *string
	Alias       string `json:"alias" gorm:"size:255"`

	Sensors  []Sensor
	Webhooks []Webhook `gorm:"foreignKey:CropPotID"`
	Controls []Control `gorm:"foreignKey:CropPotID"`
}
