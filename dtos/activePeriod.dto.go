package dtos

import "time"

type ActivePeriod struct {
	ControlID uint `json:"controlId"`
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
}