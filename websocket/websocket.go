package websocket

import (
	"PlantCare/websocket/eventHandlers"
	"PlantCare/websocket/wsTypes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"PlantCare/websocket/wsUtils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}



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
			wsutils.SendErrorResponse(connection, http.StatusBadRequest)
		}
	} else {
		fmt.Println("Unknown event:", message.Event)
		wsutils.SendErrorResponse(connection, http.StatusBadRequest)
	}
}