package eventHandlers

import (
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	"encoding/json"
	"fmt"
)

func (h *Handler) HandleAttachControl(data json.RawMessage, connection *wsTypes.Connection) {
	var ControlDto wsDtos.ControlDto

	err := json.Unmarshal(data, &ControlDto)
    if err != nil {

		
		fmt.Println("Error while unmarshaling sensor data:", err)
        return
    }

	isCustom := sensorIsCustomCheck(ControlDto.SerialNumber)
	fmt.Println(isCustom)
}

func controlIsCustomCheck(serialNumber string) (bool){
	var control models.Control
	result := initPackage.Db.First(&control, "serial_number = ?", serialNumber)

	return result != nil
}