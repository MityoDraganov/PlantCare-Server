package websocket

import (
	"PlantCare/websocket/eventHandlers"
	"PlantCare/websocket/wsTypes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"PlantCare/websocket/wsUtils"
)

// Read messages from the WebSocket connection
func HandleMessages(connection *wsTypes.Connection) {
	defer connection.Conn.Close()

	rateLimiter := wsutils.NewRateLimiter(10, time.Hour) // Initialize RateLimiter with a hardcoded limit

	for {
		_, msg, err := connection.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}

		// Process the received message
		ProcessMessage(msg, connection, rateLimiter)
	}
}

// Process the received message
func ProcessMessage(msg []byte, connection *wsTypes.Connection, rateLimiter *wsutils.RateLimiter) {
	var message wsTypes.Message
	err := json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println("Error while unmarshaling message:", err)
		return
	}

	fmt.Printf("Received message with event: %+v\n", message.Event)

	handler := &eventHandlers.Handler{}

	handlerValue := reflect.ValueOf(handler)
	method := handlerValue.MethodByName(string(message.Event))

	if method.IsValid() && method.Type().NumIn() == 2 {
		// Pass raw message data as is
		data := json.RawMessage(message.Data)

		if method.Type().In(0) == reflect.TypeOf(data) && method.Type().In(1) == reflect.TypeOf(connection) {
			args := []reflect.Value{reflect.ValueOf(data), reflect.ValueOf(connection)}

			// Wrap the method with rate limiting logic
			wrappedMethod := rateLimiter.RateLimitWrapper(func(d json.RawMessage, c *wsTypes.Connection) {
				method.Call(args)
			}, string(message.Event))

			// Call the wrapped method
			wrappedMethod(data, connection)
		} else {
			fmt.Println("Handler signature mismatch for event:", message.Event)
			wsutils.SendErrorResponse(connection, http.StatusBadRequest)
		}
	} else {
		fmt.Println("Unknown event:", message.Event)
		wsutils.SendErrorResponse(connection, http.StatusBadRequest)
	}
}

