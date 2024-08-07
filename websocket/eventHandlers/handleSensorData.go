package eventHandlers

import (
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
  

	"encoding/json"
	"fmt"
)

type Handler struct{
}

func (h *Handler) HandleUpdateSensorData(data json.RawMessage, connection *wstypes.Connection) {
    var sensorData wsDtos.SensorDataDto
    cropPotID, ok := connection.Context.Value(wstypes.CropPotIDKey).(uint)
    if !ok {
        response, _ := json.Marshal("Error")
		connection.Send <- response
		return
	}
    err := json.Unmarshal(data, &sensorData)
    if err != nil {
        fmt.Println("Error while unmarshaling sensor data:", err)
        return
    }

    fmt.Printf("Handling sensor data: %+v\n", sensorData)

    sensorDataDbObject := models.SensorData{
        Temperature: sensorData.Temperature,
        Moisture: sensorData.Moisture,    
        WaterLevel: sensorData.WaterLevel,
        SunExposure: sensorData.SunExposure,

        CropPotID: cropPotID,
    }
    
    initPackage.Db.Create(sensorDataDbObject)

    response, _ := json.Marshal("Sensor data updated!")
	connection.Send <- response
}