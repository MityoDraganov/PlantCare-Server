package models

import (
	"time"

	"gorm.io/gorm"
)

type CropPot struct {
    gorm.Model
    Token            string     `json:"token" gorm:"size:255;uniqueIndex;not null"`
    Alias            string     `json:"alias" gorm:"size:255"`
    WateringInterval int        `json:"wateringInterval"` // in minutes
    LastWateredAt    *time.Time `json:"lastWateredAt"`
    IsArchived       bool       `json:"isArchived"`
    ClerkUserID      *string    `json:"clerkUserId"`
    User            User       `gorm:"foreignKey:ClerkUserID;references:ClerkID"`
}