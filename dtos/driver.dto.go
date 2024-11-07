package dtos

type DriverDto struct {
	Id uint `json:"id"`
	DownloadUrl          string          `json:"downloadUrl"`
	MarketplaceBannerUrl string         `json:"marketplaceBannerUrl"`
	Alias                string          `json:"alias"`
	User                 UserResponseDto `json:"user"`
	IsUploader bool `json:"isUploader"`
}
