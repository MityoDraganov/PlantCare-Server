package dtos

type DriverDto struct {
	DownloadUrl          string  `json:"downloadUrl"`
	MarketplaceBannerUrl *string `json:"marketplaceBannerUrl"`
	Alias                string  `json:"alias"`
}
