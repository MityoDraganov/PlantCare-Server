package wsDtos

type SensorCommand struct {
	Command string `json:"command"`
}

type FirmwareCommand struct {
	Command string `json:"command"`
	DownloadUrl string `json:"downloadUrl"`
}