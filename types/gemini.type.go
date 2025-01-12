package types

type ModelOutput struct {
	PercentageHealthy uint8  `json:"percentageHealthy"`
	PlantName         string `json:"plantName"`
	CertentyPercantage uint8  `json:"certentyPercantage"`
	IsPlantRecognised bool   `json:"isPlantRecognised"`
}