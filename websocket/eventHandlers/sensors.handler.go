package eventHandlers

import (
	"PlantCare/controllers"
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"log"
	"net/http"
	"strconv"

	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"fmt"

	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (h *Handler) HandleMeasurements(data json.RawMessage, connection *wsTypes.Connection) {
	var sensorDataDto []wsDtos.SensorMeasuremntDto

	err := json.Unmarshal(data, &sensorDataDto)
	if err != nil {
		fmt.Println("Error while unmarshaling sensor data:", err)
		return
	}

	fmt.Printf("Handling sensor data: %+v\n", sensorDataDto)

	for _, sensorData := range sensorDataDto {

		sensorDbObject, err := controllers.FindSensorBySerialNum(sensorData.SensorSerialNumber)
		if err != nil {
			wsutils.SendErrorResponse(connection, http.StatusBadRequest)
			return
		}

		measurementData := models.Measurement{
			SensorID: sensorDbObject.ID,
			Value:    sensorData.Value,
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
					Alias:        utils.CoalesceString(sensorDbObject.Alias),
					Description:  sensorDbObject.Description,
				},
				Measurement: measurementData,
			}
			go utils.TriggerWebhook(webhook.EndpointUrl, payload)
		}

		fmt.Println(measurementDataDbObject)
	}
	wsutils.SendValidResponse(connection, nil)
}

func (h *Handler) HandleAttachSensor(data json.RawMessage, connection *wsTypes.Connection) {
	potIDStr, ok := connection.Context.Value(wsTypes.CropPotIDKey).(string)
	if !ok {
		fmt.Println("Error while reading PotId")
		return
	}

	var sensorDto dtos.AttachSensor
	err := json.Unmarshal(data, &sensorDto)
	if err != nil {
		fmt.Println("Error while unmarshaling sensor data:", err)
		return
	}

	potID64, err := strconv.ParseUint(potIDStr, 10, 32)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	potIDUint := uint(potID64)
	cropPotDbObject, err := controllers.FindCropPotById(potIDStr)
	if err != nil {
		log.Fatal("Pot not found!: " + err.Error())
	}

	fmt.Println(sensorDto.SerialNumber)
	sensorDbObject, err := controllers.FindSensorBySerialNum(sensorDto.SerialNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}
	fmt.Println("sensorDbObject")
	fmt.Println(sensorDbObject)
	fmt.Println("err")
	fmt.Println(err)

	if sensorDbObject == nil {
		fmt.Println("Sensor not found, adding a new one")
		alert := wsTypes.Alert{
			Message: "Sensor not found, adding a new one",
		}
		wsutils.SendMessage(connection, wsTypes.SensorAdded, "", alert)

		sensorDbObject, addErr := controllers.AddSensor(potIDUint, sensorDto)
		fmt.Println(addErr)
		fmt.Println(sensorDbObject)
		if addErr != nil {
			fmt.Println("Error adding sensor:", *addErr)
			wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
			return
		}

		alert = wsTypes.Alert{
			Message: "Sensor added successfully: " + sensorDbObject.SerialNumber,
		}
		wsutils.SendMessage(connection, wsTypes.SensorAdded, "", alert)
	}

	if sensorDbObject == nil {
		fmt.Println("Sensor not found or uninitialized")
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}
	err = controllers.AttachedStateUpdater(sensorDbObject, true)
	if err != nil {
		log.Fatal("Error changing attached state: ", err)
		return
	}

	alert := wsTypes.Alert{}
	sensorDriver, err := controllers.FindDriverBySensorId(sensorDbObject.ID)
	fmt.Println(sensorDriver)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatal("Error here 101")
		return
	}

	if sensorDriver != nil {
		alert.Message = "Sensor connected successfully."
		wsutils.SendMessage(connection, wsTypes.SensorConnected, "", alert)
		return
	}

	alert.Message = "Please provide a driver for the sensor."
	fmt.Println("cropPotDbObject.ClerkUserID")
	fmt.Println(*cropPotDbObject.ClerkUserID)
	controllers.CreateMessage(*cropPotDbObject.ClerkUserID, "Please provide a driver for the sensor.", "Driver required")
	wsutils.SendMessage(connection, wsTypes.DriverRequired, "", alert)

	userConn, isExisting := connectionManager.ConnManager.GetConnectionByKey(*cropPotDbObject.ClerkUserID)
	fmt.Println(isExisting)
	fmt.Println(userConn)
	if isExisting {
		wsutils.SendMessage(userConn, wsTypes.DriverRequired, "", alert)
	}

}

func (h *Handler) HandleDetachSensor(data json.RawMessage, connection *wsTypes.Connection) {
	potIDStr, ok := connection.Context.Value(wsTypes.CropPotIDKey).(string)
	if !ok {
		fmt.Println("Error while reading PotId")
		return
	}

	var sensorDto dtos.AttachSensor
	err := json.Unmarshal(data, &sensorDto)
	if err != nil {
		fmt.Println("Error while unmarshaling sensor data:", err)
		return
	}

	cropPotDbObject, err := controllers.FindCropPotById(potIDStr)
	if err != nil {
		log.Fatal("Pot not found!: " + err.Error())
	}

	fmt.Println(sensorDto.SerialNumber)
	sensorDbObject, err := controllers.FindSensorBySerialNum(sensorDto.SerialNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}
	fmt.Println("sensorDbObject")
	fmt.Println(sensorDbObject)
	fmt.Println("err")
	fmt.Println(err)

	if sensorDbObject == nil {
		fmt.Println("Sensor not found or uninitialized")
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}
	err = controllers.AttachedStateUpdater(sensorDbObject, false)
	if err != nil {
		log.Fatal("Error changing attached state: ", err)
		return
	}

	alert := wsTypes.Alert{}

	userConn, isExisting := connectionManager.ConnManager.GetConnectionByKey(*cropPotDbObject.ClerkUserID)
	fmt.Println(isExisting)
	fmt.Println(userConn)
	if isExisting {
		wsutils.SendMessage(userConn, wsTypes.DriverRequired, "", alert)
	}

}
