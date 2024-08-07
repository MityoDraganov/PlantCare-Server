package wsutils

import (
	"PlantCare/websocket/eventHandlers"
	"PlantCare/websocket/wsTypes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gorilla/websocket"
)

// Read messages from the WebSocket connection
func HandleMessages(connection *wstypes.Connection) {
	defer connection.Conn.Close()

	for {
		_, msg, err := connection.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}

		// Process the received message
		ProcessMessage(msg, connection)
	}
}

// Process the received message
func ProcessMessage(msg []byte, connection *wstypes.Connection) {
	var message wstypes.Message
	err := json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println("Error while unmarshaling message:", err)
		return
	}

	fmt.Printf("Received message with event: %+v\n", message.Event)

	handler := &eventHandlers.Handler{
	}

	handlerValue := reflect.ValueOf(handler)
	method := handlerValue.MethodByName(message.Event)

	if method.IsValid() && method.Type().NumIn() == 2 {
		// Pass raw message data as is
		data := json.RawMessage(message.Data)

		if method.Type().In(0) == reflect.TypeOf(data) && method.Type().In(1) == reflect.TypeOf(connection) {
			args := []reflect.Value{reflect.ValueOf(data), reflect.ValueOf(connection)}
			method.Call(args)
		} else {
			fmt.Println("Handler signature mismatch for event:", message.Event)
		}
	} else {
		fmt.Println("Unknown event:", message.Event)
	}

	// Echo back the message as an example
	response, _ := json.Marshal(message)
	connection.Send <- response
}



func SendMessages(connection *wstypes.Connection) {
	for msg := range connection.Send {
		err := connection.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Error while sending message:", err)
			break
		}
	}
}