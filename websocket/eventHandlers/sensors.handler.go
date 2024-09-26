package eventHandlers

import (
	"PlantCare/controllers"
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"net/http"

	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"fmt"

	"encoding/json"

	"gorm.io/gorm/clause"
)

func (h *Handler) HandleMeasurements(data json.RawMessage, connection *wsTypes.Connection) {
    var sensorDataDto wsDtos.SensorMeasuremntDto

    err := json.Unmarshal(data, &sensorDataDto)
    if err != nil {
        fmt.Println("Error while unmarshaling sensor data:", err)
        return
    }

    fmt.Printf("Handling sensor data: %+v\n", sensorDataDto)

    sensorDbObject, err := controllers.FindSensorBySerialNum(sensorDataDto.SensorSerialNumber)
    if err != nil {
		wsutils.SendErrorResponse(connection, http.StatusBadRequest)
		return
	}

    measurementData := models.Measurement{
        SensorID: sensorDbObject.ID,
        Value: sensorDataDto.Value,
    }
    
    measurementDataDbObject := initPackage.Db.Create(&measurementData).Clauses(clause.Returning{})

    if measurementDataDbObject.Error != nil {
       wsutils.SendErrorResponse(connection, http.StatusNotFound)
    }

    webhooks, err := controllers.GetSubscribedWebhooksForSensor(sensorDbObject.ID)
	if err != nil {
        wsutils.SendErrorResponse(connection, http.StatusBadRequest)
	}

    for _, webhook := range webhooks {
        payload := wsDtos.WebhookPayload{
            Sensor: dtos.SensorDto{
                SerialNumber: sensorDbObject.SerialNumber,
                Alias: sensorDbObject.Alias,
                Description: sensorDbObject.Description,
            },
            Measurement: measurementData,
        }
		go utils.TriggerWebhook(webhook.EndpointUrl, payload)
	}

    fmt.Println(measurementDataDbObject)
    wsutils.SendValidResponse(connection, nil)
}


func (h *Handler) HandleAttachSensor(data json.RawMessage, connection *wsTypes.Connection) {
	var sensorDto dtos.AttachSensor

	potID, ok := connection.Context.Value(wsTypes.CropPotIDKey).(string)
	fmt.Print(potID)
	if !ok{
		fmt.Println("Error while reading PotId")
		return
	}
	
	err := json.Unmarshal(data, &sensorDto)
	if err != nil {
		fmt.Println("Error while unmarshaling sensor data:", err)
		return
	}

	// Try to find the sensor by serial number
	sensorDbObject, err := controllers.FindSensorBySerialNum(sensorDto.SerialNumber)
	if err != nil {
		fmt.Println("Sensor not found, adding a new one")

		// Add the sensor if not found
		addedSensor, addErr := controllers.AddSensor(0, sensorDto) 
		if addErr != nil {
			fmt.Println("Error adding sensor:", *addErr)
			return
		}

		// Inform the user and provide the option to add a driver
		alert := wsTypes.Alert{
			Message: addedSensor,
		}
		wsutils.SendMessage(connection, wsTypes.SensorAdded, alert)
		return
	}

	// If the sensor exists, check if a driver is provided
	isDriverProvided, err := isDriverProvided(sensorDbObject.ID)
	if err != nil {
		fmt.Println("Error checking driver:", err)
		return
	}

	// If no driver is provided, prompt the user to add one
	alert := wsTypes.Alert{}
	if !*isDriverProvided {
		alert.Message = "Please provide a driver for the sensor."
		wsutils.SendMessage(connection, wsTypes.DriverRequired, alert)
	} else {
		alert.Message = "Sensor connected successfully."
		wsutils.SendMessage(connection, wsTypes.SensorConnected, alert)
	}


}


func isDriverProvided(sensorId uint) (*bool, error){
	var driver models.Driver
	result := initPackage.Db.First(&driver, "sensor_id = ?", sensorId)
	if result.Error != nil {
		return nil, result.Error
	}
	isProvided := result != nil
	return &isProvided , nil
}