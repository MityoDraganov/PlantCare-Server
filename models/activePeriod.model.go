package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivePeriod struct {
	gorm.Model
	ControlID uint
	Start time.Time `gorm:"type:time"`
	End   time.Time `gorm:"type:time"`
}
