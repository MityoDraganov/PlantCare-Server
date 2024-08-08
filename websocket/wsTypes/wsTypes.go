package wsTypes

import (
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn *websocket.Conn
	Send chan []byte
	Context context.Context 
}

type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"` // Raw JSON to be parsed dynamically
}

type WsResponse struct {
	Ok bool
	Status int
	Data interface{}
}
type SendMessagesFunc func(connection *Connection)

type ContextKey string
const CropPotIDKey ContextKey = "cropPotID"


