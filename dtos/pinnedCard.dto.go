package dtos

import "PlantCare/types"

type PinnedCardDto struct {
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	SensorID uint   `json:"sensorId"`
	StartLocation int    `json:"startLocation"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Type     types.CardType `json:"type"`
}
