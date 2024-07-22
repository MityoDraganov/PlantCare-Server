package models

import (
	"time"

	"gorm.io/gorm"
)

type CropPot struct {
	gorm.Model
	Alias            string     `json:"alias" gorm:"not null"`
	WateringInterval int        `json:"wateringInterval" gorm:"not null"` // in minutes
	LastWateredAt    *time.Time `json:"lastWateredAt"`
	IsArchived       bool       `json:"isArchived"`
	UserID           uint       `json:"userId"`
	User             User       `json:"user"`
}
