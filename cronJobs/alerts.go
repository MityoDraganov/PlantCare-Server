package cronjobs

import (
	"PlantCare/websocket"
	"PlantCare/websocket/wsTypes"
	"log"
	"time"
)

// Simulated alert structure
type Alert struct {
	Message   string
	Timestamp time.Time
}

// CheckAndSendAlerts checks for new alerts and sends them to users
func CheckAndSendAlerts(connection wsTypes.Connection) {
	log.Println("Checking for new alerts...")

	// Simulate fetching alerts (in a real scenario, you might fetch from a database)
	newAlerts := []Alert{
		{Message: "New message in your inbox!", Timestamp: time.Now()},
		{Message: "Don’t forget to complete your task!", Timestamp: time.Now()},
	}

	// Send alerts to all users
	for _, alert := range newAlerts {
		sendAlertToUsers(alert, &connection)
	}

}

// Simulate sending an alert to all users
func sendAlertToUsers(alert Alert, connection *wsTypes.Connection) {
	websocket.SendMessage(connection, wsTypes.ForecastAlert, alert)
}
