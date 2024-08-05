package websocket

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Connection represents a single WebSocket connection
type Connection struct {
	Conn *websocket.Conn
	Send chan []byte
}

// HandleConnection handles WebSocket connections
func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}

	connection := &Connection{
		Conn: conn,
		Send: make(chan []byte),
	}

	go handleMessages(connection)
	go SendMessages(connection)
}

// handleMessages handles incoming messages from the client
func handleMessages(connection *Connection) {
	defer func() {
		connection.Conn.Close()
	}()

	for {
		_, msg, err := connection.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}

		fmt.Printf("Received message: %s\n", msg)

		// Process message (e.g., pass to custom middleware)
		// ...

		// Example: Echo back the message
		connection.Send <- msg
	}
}

// sendMessages sends messages to the client
func SendMessages(connection *Connection) {
	for msg := range connection.Send {
		err := connection.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Error while sending message:", err)
			break
		}
	}
}
