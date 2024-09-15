package dtos

type GroupRequestDto struct {
	CropPotID uint `json:"cropPotId"`
}

type GroupResponsetDto struct {
	CropPots []CropPotResponse
}
