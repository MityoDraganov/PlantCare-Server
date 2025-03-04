package utils

import (
	"fmt"
	"time"
)

func UintPtrToString(u *uint) *string {
	if u == nil {
		return nil
	}
	s := fmt.Sprint(*u)
	return &s
}

func CoalesceString(s *string) *string {
	if s == nil {
		empty := ""
		return &empty
	}
	return s
}

func CoalesceInt(i *int) *int {
	if i == nil {
		zero := 0
		return &zero
	}
	return i
}

func StringPtr(s string) *string {
	return &s
}

func DurationToTimeString(d time.Duration) string {
	hours := int(d / time.Hour)
	minutes := int((d % time.Hour) / time.Minute)
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
