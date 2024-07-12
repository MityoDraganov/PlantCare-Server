package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"not null;unique;check:username_length_gte_8,len(username) >= 8"`
	Email        string `json:"email" gorm:"not null;unique;check:email_length_gte_3,len(email) >= 3"`
	PasswordHash string `json:"passwordHash" gorm:"not null;check:password_hash_length_gte_8,len(password_hash) >= 8"`
}
