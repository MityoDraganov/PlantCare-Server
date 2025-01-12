package utils

import (
	"fmt"
	"time"
)

func CoalesceString(s *string) *string {
	if s == nil {
		empty := ""
		return &empty
	}
	return s
}

func StringPtr(s string) *string {
	return &s
}

func DurationToTimeString(d time.Duration) string {
	hours := int(d / time.Hour)
	minutes := int((d % time.Hour) / time.Minute)
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
