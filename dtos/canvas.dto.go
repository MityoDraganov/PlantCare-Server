package dtos

type CanvasDto struct {
	CropPotID uint `json:"cropPotID"`
	PinnedCards []PinnedCardDto `json:"pinnedCards"`
}