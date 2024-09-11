package wsutils

import (
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

func SendFirmwareUpdate(conn *wsTypes.Connection, firmwareData []byte) {
	update := wsDtos.FirmwareUpdate{
		Event: "FirmwareUpdate",
		Data:  firmwareData,
	}

	// Serialize the update to JSON
	message, err := json.Marshal(update)
	if err != nil {
		fmt.Println("Error while marshaling firmware update:", err)
		return
	}

	// Send the message over WebSocket
	err = conn.Conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		fmt.Println("Error while sending firmware update:", err)
	}
}
