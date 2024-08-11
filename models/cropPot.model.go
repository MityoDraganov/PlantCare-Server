package models

import (
	"time"

	"gorm.io/gorm"
)

type CropPot struct {
    gorm.Model
    Token            string          `gorm:"size:255;uniqueIndex;not null"`
    Alias            string          `json:"alias" gorm:"size:255"`
    LastWateredAt    *time.Time
    IsArchived       bool
    ClerkUserID      *string
    User             User            `gorm:"foreignKey:ClerkUserID;references:ClerkID"`

    SensorData       []SensorData    `gorm:"foreignKey:CropPotID"`

    // Foreign key to ControlSettings
    ControlSettingsID *uint          `json:"controlSettingsId"`
    ControlSettings   *ControlSettings `gorm:"foreignKey:ControlSettingsID"`
}
