package cronjobs

import (
	"PlantCare/controllers"
	"PlantCare/lib"
	"PlantCare/services"
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsTypes"
	"fmt"
	"log"
	"time"
)

// RequestAllSensorData sends a command to all crop pots to retrieve their sensor data
func RequestAllSensorData() {
    log.Println("Requesting all sensor data from crop pots...")

    // Get all active crop pot connections
    connections := connectionManager.ConnManager.GetConnectionsByRole(wsTypes.PotRole)
    fmt.Println(connections)
    for _, connection := range connections {
        // Retrieve the crop pot ID from the WebSocket connection context
        cropPotID := connection.Context.Value(wsTypes.CropPotIDKey).(string)

        // Fetch crop pot details from the database
        cropPot, err := controllers.FindCropPotById(cropPotID)
        if err != nil {
            log.Printf("Error fetching crop pot details for ID %s: %v", cropPotID, err)
            continue
        }

        // Get the measurement interval from the crop pot details
        if cropPot.MeasuremntInterval > 0 {
            lastMeasurementTime, err := services.GetLastMeasurementTimeFromSensors(cropPot.Sensors)
			fmt.Println(lastMeasurementTime)
			if err != nil {
				fmt.Println("Error getting lastMeasurementTime")
				return;
			}

            if time.Since(lastMeasurementTime) >= cropPot.MeasuremntInterval {
				fmt.Println("sendReadAllSensorDataCommand")
                lib.SendReadAllSensorDataCommand(connection, cropPotID)
            }
        }
    }
}

