package models

import (
	"gorm.io/gorm"
)


type Passenger struct {
	gorm.Model
	PassengerInfoID uint      `json:"passenger_info_id"`
	PassengerInfo   User      `gorm:"foreignKey:PassengerInfoID"`
	StartPoint      Coordinates `json:"startPoint" validate:"required"`
	EndPoint        Coordinates `json:"endPoint" validate:"required"`
}