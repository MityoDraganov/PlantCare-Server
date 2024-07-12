package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"not null;unique"`
	Email        string `json:"email" gorm:"not null;unique"`
	PasswordHash string `json:"passwordHash" gorm:"not null"`
}
