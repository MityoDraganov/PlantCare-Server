package models

import "gorm.io/gorm"

type Driver struct {
	gorm.Model

	DownloadUrl string `json:"downloadUrl"` // The URL to download the driver (software)

	Sensors []Sensor `gorm:"foreignKey:DriverID"`

	MarketplaceBannerUrl *string
	Alias                string `json:"alias"`
	ClerkUserID     string
}
