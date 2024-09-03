package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivePeriod struct {
	gorm.Model
	ControlID uint
	Start     time.Duration
	End       time.Duration
	Days      uint8 `gorm:"type:tinyint"`
}


func (a *ActivePeriod) SetDays(days ...int) {
	var bitmask uint8
	for _, day := range days {
		bitmask |= uint8(day)
	}
	a.Days = bitmask
}

func (a *ActivePeriod) HasDay(day int) bool {
	return a.Days&uint8(day) != 0
}
