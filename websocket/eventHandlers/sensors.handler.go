package eventHandlers

import (
	"PlantCare/controllers"
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"net/http"
	"strconv"
	"time"

	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"fmt"

	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Interface for handling sensor data
type SensorDataInterface interface {
    GetSensors() []wsDtos.SensorMeasuremntDto
}

// Struct for sensor measurement data

// Struct implementing the interface for incoming data
type SensorData struct {
    Sensors []wsDtos.SensorMeasuremntDto `json:"sensors"`
}

func (sd *SensorData) GetSensors() []wsDtos.SensorMeasuremntDto {
    return sd.Sensors
}



func (h *Handler) HandleMeasurements(data json.RawMessage, connection *wsTypes.Connection) {
	fmt.Println(data)
    var sensorData SensorData
    var sensorDataDto []wsDtos.SensorMeasuremntDto
    var measurements []models.Measurement

	err := json.Unmarshal(data, &sensorData)
    if err != nil {
        fmt.Println("Error while unmarshaling sensor data:", err)
        return
    }

    // Store the extracted data from sensorData.Sensors into sensorDataDto
    sensorDataDto = sensorData.GetSensors()

	role := connection.Role
	potIdStr := connection.Context.Value(wsTypes.CropPotIDKey).(string)
	cropPotDBObject, err := controllers.FindCropPotById(potIdStr)
	if err != nil {
		fmt.Println("Error finding crop pot by id:", err)
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}

	potId, err := strconv.ParseUint(potIdStr, 10, 32)
	if err != nil {
		fmt.Println("Error converting potIdStr to uint:", err)
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}
	measurementGroup := models.MeasurementGroup{
		CropPotID: uint(potId),
	}

	if err := initPackage.Db.Create(&measurementGroup).Clauses(clause.Returning{}).Error; err != nil {
		fmt.Println("Error creating measurement group:", err)
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
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
			SensorID:           sensorDbObject.ID,
			Value:              sensorData.Value,
			Role:               role,
			MeasurementGroupID: measurementGroup.ID,
		}

		measurementDataDbObject := initPackage.Db.Create(&measurementData).Clauses(clause.Returning{})

		if measurementDataDbObject.Error != nil {
			wsutils.SendErrorResponse(connection, http.StatusNotFound)
		}
		measurements = append(measurements, measurementData)

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

	ownerConn, isExisting := connectionManager.ConnManager.GetConnectionByKey(*cropPotDBObject.ClerkUserID)
	if isExisting {
		fmt.Println("Owner connection found")

		measurementResponse := dtos.MeasurementGroupDto{
			MeasurementGroupID: measurementGroup.ID,
			CropPotID:          measurementGroup.CropPotID,
			Measurements:       measurements,
			HealthStatus:       nil,
		}

		measurementResponseString, err := json.Marshal(measurementResponse)
		if err != nil {
			fmt.Printf("Failed to marshal alert: %v\n", err)
		}

		// messageDto := wsDtos.NotificationDto{
		//     Title:     utils."New measurement",
		//     Data:      measurementResponse,
		//     IsRead:    false,
		//     Timestamp: time.Now(),
		// }

		controllers.CreateMessage(*cropPotDBObject.ClerkUserID, string(measurementResponseString), "New measurement", wsTypes.MessageFound, wsTypes.UndiagnosedMeasurement)
		notification := wsDtos.NotificationDto{
			Title:     utils.StringPtr("New measurement"),
			Data:      measurementResponse,
			IsRead:    false,
			Timestamp: time.Now(),
		}
		wsutils.SendMessage(ownerConn, wsTypes.MessageFound, wsTypes.UndiagnosedMeasurement, notification)
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
		fmt.Println(err.Error())
		return
	}
	potIDUint := uint(potID64)
	cropPotDbObject, err := controllers.FindCropPotById(potIDStr)
	if err != nil {
		fmt.Println("Pot not found!: " + err.Error())
	}

	fmt.Println(sensorDto.SerialNumber)
	sensorDbObject, err := controllers.FindSensorBySerialNum(sensorDto.SerialNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}

	if sensorDbObject == nil {
		fmt.Println("Sensor not found, adding a new one")
		alert := wsDtos.NotificationDto{
			Title:     utils.StringPtr("Sensor not found, adding a new one"),
			Data:      nil,
			IsRead:    false,
			Timestamp: time.Now(),
		}
		wsutils.SendMessage(connection, wsTypes.SensorAdded, "", alert)

		sensorDbObject, addErr := controllers.AddSensor(potIDUint, sensorDto)
		if addErr != nil {
			fmt.Println("Error adding sensor:", *addErr)
			wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
			alert = wsDtos.NotificationDto{
				Title:     utils.StringPtr("Sensor added successfully: " + sensorDbObject.SerialNumber),
				Data:      nil,
				IsRead:    false,
				Timestamp: time.Now(),
			}
			alert = wsDtos.NotificationDto{
				Title:     utils.StringPtr("Sensor added successfully"),
				Data:      sensorDbObject.SerialNumber,
				IsRead:    false,
				Timestamp: time.Now(),
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
			fmt.Println("Error changing attached state: ", err)
			return
		}

		alert = wsDtos.NotificationDto{}
		sensorDriver, err := controllers.FindDriverBySensorId(sensorDbObject.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println("Error while finding driver: ", err)
			return
		}

		if sensorDriver != nil {
			alert.Title = utils.StringPtr("Sensor connected successfully.")
			wsutils.SendMessage(connection, wsTypes.SensorConnected, "", alert)
			return
		}

		alert.Title = utils.StringPtr("Please provide a driver for the sensor.")
		controllers.CreateMessage(*cropPotDbObject.ClerkUserID, "Please provide a driver for the sensor.", "Driver required", wsTypes.DriverRequired, "")
		wsutils.SendMessage(connection, wsTypes.DriverRequired, "", alert)

		userConn, isExisting := connectionManager.ConnManager.GetConnectionByKey(*cropPotDbObject.ClerkUserID)
		if isExisting {
			wsutils.SendMessage(userConn, wsTypes.DriverRequired, "", alert)
		}

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
		fmt.Println("Pot not found!: " + err.Error())
	}

	fmt.Println(sensorDto.SerialNumber)
	sensorDbObject, err := controllers.FindSensorBySerialNum(sensorDto.SerialNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}

	if sensorDbObject == nil {
		fmt.Println("Sensor not found or uninitialized")
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}
	err = controllers.AttachedStateUpdater(sensorDbObject, false)
	if err != nil {
		fmt.Println("Error changing attached state: ", err)
		return
	}

	alert := wsDtos.NotificationDto{}

	userConn, isExisting := connectionManager.ConnManager.GetConnectionByKey(*cropPotDbObject.ClerkUserID)
	if isExisting {
		wsutils.SendMessage(userConn, wsTypes.DriverRequired, "", alert)
	}

}
