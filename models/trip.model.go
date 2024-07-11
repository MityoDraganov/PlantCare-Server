package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Coordinates [2]float64


// Scan implements the sql.Scanner interface for Coordinates
func (c *Coordinates) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert value to byte slice")
	}

	err := json.Unmarshal(bytes, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal coordinates: %v", err)
	}

	return nil
}

func (c Coordinates) Value() (driver.Value, error) {
	return json.Marshal(c)
}

type Trip struct {
	gorm.Model
	DriverID           uint         `json:"driver_id"`
	Driver             User         `gorm:"foreignKey:DriverID"`
	CurrentPassengers  []Passenger  `gorm:"many2many:trip_passengers;"`
	StartPoint         Coordinates  `json:"startPoint" validate:"required"`
	EndPoint           Coordinates   `json:"endPoint" validate:"required"`
	MaxAvailableSeats  int          `json:"maxAvailableSeats" validate:"required"`
	Description        string       `json:"description"`
}