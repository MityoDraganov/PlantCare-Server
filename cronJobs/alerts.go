package cronjobs

import (
	"PlantCare/controllers"
	"PlantCare/services"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"encoding/json"
	"fmt"
	"log"
)

// Simulated alert structure

// CheckAndSendAlerts checks for new alerts and sends them to users
func CheckAndSendAlerts(connection wsTypes.Connection) {
	log.Println("Checking for new alerts...")
	userID := connection.Context.Value(wsTypes.CropPotIDKey).(string)


	cropPots, err := controllers.FindPotsByUserId(userID)
	if err != nil {
		fmt.Println("Failed to find crop pots: %v", err)
		return;
	}

	historicalIndoorsWeather := controllers.GetMeasurementsBySensorId(cropPots[0].Sensors[0].ID)

	forecast, err := services.GetIndoorForecast("Sliven", historicalIndoorsWeather)
	if err != nil {
		fmt.Println(err.Error())
	}

	alert := wsTypes.Alert{
		Message:   forecast,
	}

	alertString, err := json.Marshal(alert)
	if err != nil {
		fmt.Println("Failed to marshal alert: %v", err)
	}

	controllers.CreateMessage(userID, string(alertString), "New forecast alert", "", wsTypes.HandleForecastAlert)
	sendAlertToUsers(alert, &connection)

}

// Simulate sending an alert to all users
func sendAlertToUsers(alert wsTypes.Alert, connection *wsTypes.Connection) {
	wsutils.SendMessage(connection, "", wsTypes.HandleForecastAlert, alert)
}
