package models

import "gorm.io/gorm"

type Driver struct {
	gorm.Model
	SensorID uint

	DownloadUrl string
}