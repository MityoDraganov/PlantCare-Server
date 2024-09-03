package dtos

type ActivePeriod struct {
	ID    uint   `json:"id"`
	Start string `json:"start"`
	End   string `json:"end"`
	Days  []int `json:"days"`
}
