package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivePeriod struct {
	gorm.Model
	ControlID uint
	Start time.Duration
	End   time.Duration
}
