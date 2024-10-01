package wsDtos

import "time"

type NotificationDto struct {
	Title *string `json:"title"`
	Text   string `json:"text"`
	IsRead bool   `json:"isRead"`
	Timestamp time.Time `json:"timestamp"`
}
