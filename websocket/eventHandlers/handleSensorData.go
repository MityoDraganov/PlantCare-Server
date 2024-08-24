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


func (h *Handler) HandleMeasurements(data json.RawMessage, connection *wsTypes.Connection) {
    var sensorDataDto wsDtos.SensorMeasuremntDto

    err := json.Unmarshal(data, &sensorDataDto)
    if err != nil {
        fmt.Println("Error while unmarshaling sensor data:", err)
        return
    }

    fmt.Printf("Handling sensor data: %+v\n", sensorDataDto)

    measurementData := models.Measurement{
        SensorID: sensorDataDto.SensorID,
        Value: sensorDataDto.Value,
    }
    
    measurementDataDbObject := initPackage.Db.Create(&measurementData).Clauses(clause.Returning{})

    if measurementDataDbObject.Error != nil {
       wsutils.SendErrorResponse(connection, http.StatusNotFound)
    }

    fmt.Println(measurementDataDbObject)
    wsutils.SendValidResponse(connection, nil)
}