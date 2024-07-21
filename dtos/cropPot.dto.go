package dtos

type CreateCropPot struct {
	Alias            string `json:"alias" validate:"required,min=8"`
	WateringInterval int    `json:"wateringInterval" validate:"required"` // in minutes
}
