package eventHandlers

import (
	"PlantCare/controllers"
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"net/http"

	"encoding/json"
	"fmt"

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

