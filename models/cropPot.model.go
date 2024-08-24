package models

import (
	"time"

	"gorm.io/gorm"
)

type CropPot struct {
	gorm.Model
	IsPinned    bool
	IsArchived  bool
	Token       string `gorm:"size:255;uniqueIndex;not null"`
	User        User   `gorm:"foreignKey:ClerkUserID;references:ClerkID"`
	ClerkUserID *string
	Alias       string `json:"alias" gorm:"size:255"`

	LastWateredAt *time.Time

	Sensors  []Sensor
	Webhooks []Webhook `gorm:"foreignKey:CropPotID"`

	ControlSettingsID *uint `json:"controlSettingsId"`
	ControlSettings   *ControlSettings
}
