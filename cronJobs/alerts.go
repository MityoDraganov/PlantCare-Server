package cronjobs

import (
	"PlantCare/services"
	"PlantCare/websocket"
	"PlantCare/websocket/wsTypes"
	"log"
	"time"
)

// Simulated alert structure
type Alert struct {
	Message   interface{}
	Timestamp time.Time
}

// CheckAndSendAlerts checks for new alerts and sends them to users
func CheckAndSendAlerts(connection wsTypes.Connection) {
	log.Println("Checking for new alerts...")
	userID := connection.Context.Value(wsTypes.CropPotIDKey).(string)
	forecast, err := services.GetIndoorForecast("Sliven", userID)
	if err != nil {
		log.Fatal(err.Error())
	}

	alert := Alert{
		Message:   forecast,
		Timestamp: time.Now(),
	}

	sendAlertToUsers(alert, &connection)

}

// Simulate sending an alert to all users
func sendAlertToUsers(alert Alert, connection *wsTypes.Connection) {
	websocket.SendMessage(connection, wsTypes.ForecastAlert, alert)
}
