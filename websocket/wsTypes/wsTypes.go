package wsTypes

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Role string
const (
	UserRole Role = "user"
	PotRole  Role = "pot"
)

type Event string
const (
	ForecastAlert Event = "forecastAlert"
)

type Connection struct {
	Conn    *websocket.Conn
	Send    chan []byte
	Context context.Context
	IP      string
	Role    Role
}

type Message struct {
	Event Event          `json:"event"`
	Data  json.RawMessage `json:"data"` // Raw JSON to be parsed dynamically
	Timestamp time.Time `json:"timestamp"`
}

type WsResponse struct {
	Ok     bool
	Status int
	Data   interface{}
}
type SendMessagesFunc func(connection *Connection)

type ContextKey string

const CropPotIDKey ContextKey = "cropPotID"
const UserIDKey ContextKey = "ClerkID"
