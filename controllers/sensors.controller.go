package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"fmt"
	"net/http"
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
	id := params["sensorId"]

	sensorDbObject, err := findSensorById(id)
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

func FindSensorBySerialNum(serialNumber string) (*models.Sensor, error) {
	var sensorDbObject models.Sensor
	result := initPackage.Db.Where(&models.Sensor{SerialNumber: serialNumber}).First(&sensorDbObject)

	if result.Error != nil {
		return nil, result.Error
	}

	return &sensorDbObject, nil
}

func findSensorById(id string) (*models.Sensor, error) {

	var sensor models.Sensor
	result := initPackage.Db.First(&sensor, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sensor, nil
}
