package wsDtos

import "time"

type NotificationDto struct {
	Title *string `json:"title"`
	Data   interface{} `json:"data"`
	IsRead bool   `json:"isRead"`
	Timestamp time.Time `json:"timestamp"`
}
