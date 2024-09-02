package models

import "gorm.io/gorm"

type Update struct {
	gorm.Model
	ControlID uint
}