package models

import (
	"PlantCare/websocket/wsTypes"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	StatusResponse wsTypes.StatusResponse
	Event          wsTypes.Event
	Title          *string
	Data           string
	IsRead         bool
	ClerkUserID    string // Foreign key to reference the User's ClerkID
	User           *User  `gorm:"foreignKey:ClerkUserID;references:ClerkID"` // User relation defined via ClerkUserID
}
