package eventHandlers

import (
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"net/http"

	"encoding/json"
	"fmt"

	"gorm.io/gorm/clause"
)


func (h *Handler) HandleUpdateSensorData(data json.RawMessage, connection *wstypes.Connection) {
    var sensorDataDto wsDtos.SensorDataDto
    cropPotID, ok := connection.Context.Value(wstypes.CropPotIDKey).(uint)
    if !ok {
        response, _ := json.Marshal("Error")
		connection.Send <- response
		return
	}
    err := json.Unmarshal(data, &sensorDataDto)
    if err != nil {
        fmt.Println("Error while unmarshaling sensor data:", err)
        return
    }

    fmt.Printf("Handling sensor data: %+v\n", sensorDataDto)

    sensorData := models.SensorData{
        Temperature: sensorDataDto.Temperature,
        Moisture: sensorDataDto.Moisture,    
        WaterLevel: sensorDataDto.WaterLevel,
        SunExposure: sensorDataDto.SunExposure,

        CropPotID: cropPotID,
    }
    
    sensorDataDbObject := initPackage.Db.Create(&sensorData).Clauses(clause.Returning{})

    if sensorDataDbObject.Error != nil {
       wsutils.SendErrorResponse(connection, http.StatusNotFound)
    }

    fmt.Println(sensorDataDbObject)
    wsutils.SendValidResponse(connection, sensorDataDbObject)
}