package dtos

import "time"

type ActivePeriod struct {
	ID uint `json:"id"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}