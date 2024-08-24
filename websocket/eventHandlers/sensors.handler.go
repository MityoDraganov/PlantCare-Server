package eventHandlers

import (
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	"fmt"

	"encoding/json"
)


func (h *Handler) HandleAttachSensor(data json.RawMessage, connection *wsTypes.Connection) {
	var sensorDto wsDtos.SensorDTO

	err := json.Unmarshal(data, &sensorDto)
    if err != nil {

		
		fmt.Println("Error while unmarshaling sensor data:", err)
        return
    }

	isCustom := sensorIsCustomCheck(sensorDto.SerialNumber)
	fmt.Println(isCustom)
}

func sensorIsCustomCheck(serialNumber string) (bool){
	var sensor models.Sensor
	result := initPackage.Db.First(&sensor, "serial_number = ?", serialNumber)

	return result != nil
}