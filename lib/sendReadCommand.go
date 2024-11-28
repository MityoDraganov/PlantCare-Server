package lib

import (
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"log"
)

// sendReadAllSensorDataCommand sends the readAllSensorData command over WebSocket
func SendReadAllSensorDataCommand(connection *wsTypes.Connection, cropPotID string) {
    command := wsDtos.SensorCommand{
        Command: "readAllSensorData",
    }


    // Send the command via WebSocket
    err := wsutils.SendMessage(connection, "", wsTypes.HandleSensorDataRequest, command)
    if err != nil {
        log.Printf("Failed to send sensor data request to crop pot %s: %v", cropPotID, err)
        connectionManager.ConnManager.RemoveConnection(cropPotID);
        return
    }

    log.Printf("Sensor data request sent to crop pot %s", cropPotID)
}
