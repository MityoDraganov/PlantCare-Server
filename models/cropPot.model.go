package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	StatusOnline   Status = "online"
	StatusUpdating Status = "updating"
	StatusOffline  Status = "offline"
)

type CropPot struct {
	gorm.Model
	GroupID     *uint
	Token       string `gorm:"size:255;uniqueIndex;not null"`
	User        User   `gorm:"foreignKey:ClerkUserID;references:ClerkID"`
	ClerkUserID *string
	Status      Status

	Alias      string `json:"alias" gorm:"size:255"`
	IsArchived bool
	IsPinned   bool `json:"isPinned"`

	Sensors  []Sensor
	Webhooks []Webhook `gorm:"foreignKey:CropPotID"`
	Controls []Control `gorm:"foreignKey:CropPotID"`

	MeasuremntInterval time.Duration `json:"measuremntInterval"`

	Canvas Canvas `gorm:"foreignKey:CropPotID"`
}
