package models

import "gorm.io/gorm"

type Driver struct {
	gorm.Model

	DownloadUrl string `json:"downloadUrl"`

	Sensors []Sensor `gorm:"foreignKey:DriverID"`
	Controls []Control `gorm:"foreignKey:DriverID"`

	MarketplaceBannerUrl *string
	Alias                string `json:"alias"`
	UploadedByUserID     string
	IsMarketplaceFeatured bool
}
