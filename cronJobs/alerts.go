package cronjobs

import (
	"PlantCare/controllers"
	"PlantCare/services"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"encoding/json"
	"log"
)

// Simulated alert structure

// CheckAndSendAlerts checks for new alerts and sends them to users
func CheckAndSendAlerts(connection wsTypes.Connection) {
	log.Println("Checking for new alerts...")
	userID := connection.Context.Value(wsTypes.CropPotIDKey).(string)
	forecast, err := services.GetIndoorForecast("Sliven", userID)
	if err != nil {
		log.Fatal(err.Error())
	}

	alert := wsTypes.Alert{
		Message:   forecast,
	}

	alertString, err := json.Marshal(alert)
	if err != nil {
		log.Fatalf("Failed to marshal alert: %v", err)
	}

	controllers.CreateMessage(userID, string(alertString))
	sendAlertToUsers(alert, &connection)

}

// Simulate sending an alert to all users
func sendAlertToUsers(alert wsTypes.Alert, connection *wsTypes.Connection) {
	wsutils.SendMessage(connection, "", wsTypes.HandleForecastAlert, alert)
}
