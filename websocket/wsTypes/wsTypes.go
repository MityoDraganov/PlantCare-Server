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
	MLRole   Role = "ml"
)

type Event string

const (
	HandleForecastAlert Event = "ForecastNotification"
	HandleAttachSensor  Event = "handleAttachSensor"
	HandleDetachSensor  Event = "handleDetachSensor"
	UpdatedPot          Event = "updatedPot"

	NotificationAlert Event = "notificationAlert"

	AsyncError   Event = "asyncError"
	AsyncPromise Event = "asyncPromise"
	HandleSensorDataRequest Event = "HandleSensorDataRequest"
	FirmwareUpdate Event = "FirmwareUpdate"

	UndiagnosedMeasurement Event = "UndiagnosedMeasurement"
)

type StatusResponse string

const (
	SensorConnected StatusResponse = "sensorConnected"
	SensorDetached  StatusResponse = "sensorDetached"
	SensorNotFound  StatusResponse = "sensorNotFound"
	SensorAdded     StatusResponse = "sensorAdded"
	DriverRequired  StatusResponse = "driverRequired"

	MessageFound StatusResponse = "messageFound"
)

type Connection struct {
	Conn    *websocket.Conn
	Send    chan []byte
	Context context.Context
	IP      string
	Role    Role
}

type Message struct {
	StatusResponse *StatusResponse `json:"statusResponse,omitempty"`
	Event          *Event          `json:"event,omitempty"`
	Data           json.RawMessage `json:"data"`
	Timestamp      time.Time       `json:"timestamp"`
	IsRead         bool            `json:"isRead"`
}

type WsResponse struct {
	Ok     bool
	Status int
	Data   interface{}
}

type SendMessagesFunc func(connection *Connection)

type ContextKey string

const (
	CropPotIDKey ContextKey = "cropPotID"
	UserIDKey    ContextKey = "ClerkID"
)

type Alert struct {
	Message interface{}
}

