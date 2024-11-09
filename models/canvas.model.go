package models

import "gorm.io/gorm"


type Canvas struct {
	gorm.Model
	CropPotID   uint         `json:"cropPotId"`
	PinnedCards []PinnedCard `json:"pinnedCards" gorm:"constraint:OnDelete:CASCADE"`
}