package cronjobs

import (
	"PlantCare/controllers"
	"PlantCare/services"
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// RequestAllSensorData sends a command to all crop pots to retrieve their sensor data
func RequestAllSensorData() {
    log.Println("Requesting all sensor data from crop pots...")

    // Get all active crop pot connections
    connections := connectionManager.ConnManager.GetConnectionsByRole(wsTypes.PotRole)
    
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

			if err != nil {
				fmt.Println("Error getting lastMeasurementTime")
				return;
			}

            if time.Since(lastMeasurementTime) >= cropPot.MeasuremntInterval {
                // Send the readAllSensorData command if the measurement interval has passed
                sendReadAllSensorDataCommand(connection, cropPotID)
            } else {
                log.Printf("Skipping sensor data request for crop pot %s. Interval not met.", cropPotID)
            }
        }
    }
}

// sendReadAllSensorDataCommand sends the readAllSensorData command over WebSocket
func sendReadAllSensorDataCommand(connection *wsTypes.Connection, cropPotID string) {
    command := wsDtos.SensorCommand{
        Command: "readAllSensorData",
    }

    // Marshal the command into JSON
    commandBytes, err := json.Marshal(command)
    if err != nil {
        log.Printf("Failed to marshal sensor command for crop pot %s: %v", cropPotID, err)
        return
    }

    // Send the command via WebSocket
    err = wsutils.SendMessage(connection, "", wsTypes.HandleSensorDataRequest, commandBytes)
    if err != nil {
        log.Printf("Failed to send sensor data request to crop pot %s: %v", cropPotID, err)
        return
    }

    log.Printf("Sensor data request sent to crop pot %s", cropPotID)
}
