package wsDtos

type FirmwareUpdate struct {
    Event    string `json:"Event"`
    Data     []byte `json:"Data"`  // This will hold the firmware data
}