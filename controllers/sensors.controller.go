package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
)

func UpdateSensor(w http.ResponseWriter, r *http.Request) {
	var sensorDto dtos.SensorRequestDto

	// Decode the JSON body into webhookDto
	err := json.NewDecoder(r.Body).Decode(&sensorDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the webhook ID from the URL parameters
	params := mux.Vars(r)
	stringId := params["sensorId"]
	id, err := strconv.ParseUint(stringId, 10, 32)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sensorDbObject, err := findSensorById(uint(id))
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var intervalUpdateTime time.Duration
	if sensorDto.MeasurementInterval != "" {
		t, err := time.Parse("15:04", sensorDto.MeasurementInterval)
		if err != nil {
			utils.JsonError(w, fmt.Sprintf("Invalid start time format: %s", err.Error()), http.StatusBadRequest)
			return
		}
		intervalUpdateTime = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
	}

	sensorUpdate := models.Sensor{
		Alias:              sensorDto.Alias,
		Description:        sensorDto.Description,
		MeasuremntInterval: intervalUpdateTime,
	}

	result := initPackage.Db.Model(sensorDbObject).Updates(sensorUpdate).Clauses(clause.Returning{})
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	// Set the response header and encode the updated webhook object
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sensorDto); err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetMeasurementsBySensorId(id uint) dtos.SensorMeasurementsSummary{
	sensor, err := findSensorById(id)
	if err != nil {
		log.Fatal(err)
	}

	SensorMeasurementsSummaryDto := dtos.SensorMeasurementsSummary{
		SensorType: sensor.Type,
		Measurements: sensor.Measurements,
	}
	return SensorMeasurementsSummaryDto
}

func FindSensorBySerialNum(serialNumber string) (*models.Sensor, error) {
	var sensorDbObject models.Sensor
	result := initPackage.Db.Where(&models.Sensor{SerialNumber: serialNumber}).First(&sensorDbObject)

	if result.Error != nil {
		return nil, result.Error
	}

	return &sensorDbObject, nil
}

func findSensorById(id uint) (*models.Sensor, error) {

	var sensor models.Sensor
	result := initPackage.Db.Preload("Measurements").First(&sensor, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sensor, nil
}

// Maps a single Sensor to SensorResponseDto
func MapSensorToDTO(sensor models.Sensor) dtos.SensorResponseDto {
	return dtos.SensorResponseDto{
		ID:                  sensor.ID,
		SerialNumber:        sensor.SerialNumber,
		Alias:               sensor.Alias,
		Description:         sensor.Description,
		MeasurementInterval: utils.DurationToTimeString(sensor.MeasuremntInterval),
		Measurements:        sensor.Measurements,
	}
}

// Converts a single Sensor or a slice of Sensors to a slice of SensorResponseDto
func ToSensorsDTO(input interface{}) []dtos.SensorResponseDto {
	switch v := input.(type) {
	case models.Sensor:
		// If it's a single sensor, wrap it in a slice
		return []dtos.SensorResponseDto{MapSensorToDTO(v)}
	case []models.Sensor:
		// If it's a slice of sensors, map each sensor to SensorResponseDto
		sensorDTOs := make([]dtos.SensorResponseDto, len(v))
		for i, sensor := range v {
			sensorDTOs[i] = MapSensorToDTO(sensor)
		}
		return sensorDTOs
	default:
		// Handle unexpected types by returning an empty slice
		return []dtos.SensorResponseDto{}
	}
}

