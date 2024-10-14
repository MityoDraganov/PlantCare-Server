package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Title      *string
	Data       string
	IsRead      bool
	ClerkUserID string // Foreign key to reference the User's ClerkID
	User        *User  `gorm:"foreignKey:ClerkUserID;references:ClerkID"` // User relation defined via ClerkUserID
}
