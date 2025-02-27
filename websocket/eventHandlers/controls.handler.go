package eventHandlers

import (
	"PlantCare/controllers"
	"PlantCare/dtos"
	"PlantCare/utils"
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func (h *Handler) HandleAttachControl(data json.RawMessage, connection *wsTypes.Connection) {
	potIDStr, ok := connection.Context.Value(wsTypes.CropPotIDKey).(string)
	if !ok {
		fmt.Println("Error while reading PotId")
		return
	}

	var controlDto dtos.AttachControlDto
	err := json.Unmarshal(data, &controlDto)
	if err != nil {
		fmt.Println("Error while unmarshaling control data:", err)
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

	fmt.Println(controlDto.SerialNumber)
	controlDbObject, err := controllers.FindControlBySerialNum(controlDto.SerialNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}

	if controlDbObject == nil {
		fmt.Println("Control not found, adding a new one")
		alert := wsDtos.NotificationDto{
			Title:     utils.StringPtr("Control not found, adding a new one"),
			Data:      nil,
			IsRead:    false,
			Timestamp: time.Now(),
		}
		wsutils.SendMessage(connection, wsTypes.ControlAdded, "", alert)

		controlDbObject, addErr := controllers.AddControl(potIDUint, controlDto)
		if addErr != nil  {
			fmt.Println("Error adding control:", *addErr)
			wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
			alert = wsDtos.NotificationDto{
				Title:     utils.StringPtr("Control added successfully: " + controlDbObject.SerialNumber),
				Data:      nil,
				IsRead:    false,
				Timestamp: time.Now(),
			}
			alert = wsDtos.NotificationDto{
				Title:     utils.StringPtr("Control added successfully"),
				Data:      controlDbObject.SerialNumber,
				IsRead:    false,
				Timestamp: time.Now(),
			}
			wsutils.SendMessage(connection, wsTypes.ControlAdded, "", alert)
		}

		if controlDbObject == nil {
			fmt.Println("Control not found or uninitialized")
			wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
			return
		}
		err = controllers.AttachedStateUpdater(controlDbObject, true)
		if err != nil {
			fmt.Println("Error changing attached state: ", err)
			return
		}

		alert = wsDtos.NotificationDto{}
		controlDriver, err := controllers.FindDriverByControlId(controlDbObject.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println("Error while finding driver: ", err)
			return
		}

		if controlDriver != nil {
			alert.Title = utils.StringPtr("Control connected successfully.")
			wsutils.SendMessage(connection, wsTypes.ControlConnected, "", alert)
			return
		}

		alert.Title = utils.StringPtr("Please provide a driver for the control.")
		controllers.CreateMessage(*cropPotDbObject.ClerkUserID, "Please provide a driver for the control.", "Driver required", wsTypes.DriverRequired, "")
		wsutils.SendMessage(connection, wsTypes.DriverRequired, "", alert)

		userConn, isExisting := connectionManager.ConnManager.GetConnectionByKey(*cropPotDbObject.ClerkUserID)
		if isExisting {
			wsutils.SendMessage(userConn, wsTypes.DriverRequired, "", alert)
		}

	}


}

func (h *Handler) HandleDetachControl(data json.RawMessage, connection *wsTypes.Connection) {
	potIDStr, ok := connection.Context.Value(wsTypes.CropPotIDKey).(string)
	if !ok {
		fmt.Println("Error while reading PotId")
		return
	}

	var controlDto dtos.AttachControlDto
	err := json.Unmarshal(data, &controlDto)
	if err != nil {
		fmt.Println("Error while unmarshaling control data:", err)
		return
	}

	cropPotDbObject, err := controllers.FindCropPotById(potIDStr)
	if err != nil {
		fmt.Println("Pot not found!: " + err.Error())
	}

	fmt.Println(controlDto.SerialNumber)
	controlDbObject, err := controllers.FindControlBySerialNum(controlDto.SerialNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}

	if controlDbObject == nil {
		fmt.Println("Control not found or uninitialized")
		wsutils.SendErrorResponse(connection, http.StatusInternalServerError)
		return
	}
	err = controllers.AttachedStateUpdater(controlDbObject, false)
	if err != nil {
		fmt.Println("Error changing attached state: ", err)
		return
	}

	alert := wsDtos.NotificationDto{}

	userConn, isExisting := connectionManager.ConnManager.GetConnectionByKey(*cropPotDbObject.ClerkUserID)
	if isExisting {
		wsutils.SendMessage(userConn, wsTypes.DriverRequired, "", alert)
	}

	fmt.Println("Control detached: " + controlDto.SerialNumber)
}

// func controlIsCustomCheck(serialNumber string) (bool){
// 	var control models.Control
// 	result := initPackage.Db.First(&control, "serial_number = ?", serialNumber)

// 	return result != nil
// }