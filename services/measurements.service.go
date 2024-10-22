package services

import (
	"PlantCare/models"
	"log"
	"time"
)

// GetLastMeasurementTimeFromSensors retrieves the most recent measurement time from an array of sensors
func GetLastMeasurementTimeFromSensors(sensors []models.Sensor) (time.Time, error) {
	var latestTime time.Time
	initialized := false

	// Loop through each sensor
	for _, sensor := range sensors {
		// Check if the sensor has any measurements
		if len(sensor.Measurements) == 0 {
			log.Printf("No measurements found for sensor ID %d", sensor.ID)
			continue
		}

		// Loop through the measurements of the sensor
		for _, measurement := range sensor.Measurements {
			if !initialized || measurement.CreatedAt.After(latestTime) {
				latestTime = measurement.CreatedAt
				initialized = true
			}
		}
	}

	if !initialized {
		log.Println("No measurements found across all sensors.")
		return time.Time{}, nil
	}

	return latestTime, nil
}
