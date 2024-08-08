package wsutils

import (
	"PlantCare/websocket/wsTypes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)



func sendResponse(connection *wsTypes.Connection, response wsTypes.WsResponse) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return
	}

	select {
	case connection.Send <- responseBytes:
		// Message successfully queued for sending
	default:
		// Handle the case where sending is blocked or the channel is full
		fmt.Println("Error: Unable to send response, channel may be blocked.")
	}
}

// SendOkResponse creates and sends a response with Ok = true.
func SendValidResponse(connection *wsTypes.Connection, data interface{}) {
	response := wsTypes.WsResponse{
		Ok:   true,
		Status: 200,
		Data: toJSON(data),
	}

	sendResponse(connection, response)
}

// SendErrorResponse creates and sends a response with Ok = false and an optional status message.
func SendErrorResponse(connection *wsTypes.Connection, status int) {
	response := wsTypes.WsResponse{
		Ok:     false,
		Status: status,
		Data:   nil, // No data in error response
	}

	sendResponse(connection, response)
}

func SendMessages(connection *wsTypes.Connection) {
	for msg := range connection.Send {
		err := connection.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Error while sending message:", err)
			break
		}
	}
}

func toJSON(data interface{}) json.RawMessage {
	if data == nil {
		return nil
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return nil
	}
	return dataBytes
}