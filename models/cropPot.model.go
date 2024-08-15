package models

import (
	"time"

	"gorm.io/gorm"
)
type CropPot struct {
    gorm.Model
    IsPinned bool
    Token              string             `gorm:"size:255;uniqueIndex;not null"`
    User               User               `gorm:"foreignKey:ClerkUserID;references:ClerkID"`
    Alias              string             `json:"alias" gorm:"size:255"`
    LastWateredAt      *time.Time
    IsArchived         bool
    ClerkUserID        *string
    SensorDatas        []SensorData       `gorm:"foreignKey:CropPotID"` // Removed the pointer to the slice
    CustomSensorFields []CustomSensorField `gorm:"foreignKey:CropPotID"` // Changed to plural and removed pointer
    ControlSettingsID  *uint              `json:"controlSettingsId"`
    ControlSettings    *ControlSettings
}
