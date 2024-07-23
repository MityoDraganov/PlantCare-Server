package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    ClerkID   string    `json:"clerkId" gorm:"size:255;not null;unique"`
    CropPots  []CropPot `gorm:"foreignKey:ClerkUserID;references:ClerkID"`
    IsAdmin_  bool
}