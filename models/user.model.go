package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    ClerkID   string    `json:"clerkId" gorm:"size:255;not null;unique"`
    CropPots  []CropPot `gorm:"foreignKey:ClerkUserID;references:ClerkID"`  // Foreign key for CropPots
    IsAdmin_  bool
    Inbox     []Message  `gorm:"foreignKey:ClerkUserID;references:ClerkID"` // Specify the foreign key for Inbox
}
