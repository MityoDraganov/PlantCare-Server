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
)

func UpdateControls(w http.ResponseWriter, r *http.Request) {
	var controlDtos []dtos.ControlRequestDto
	err := json.NewDecoder(r.Body).Decode(&controlDtos)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx := initPackage.Db.Begin()
	if tx.Error != nil {
		utils.JsonError(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	for _, controlDto := range controlDtos {
		// Fetch the existing ActivePeriod record
		activePeriodUpdate, err := findActivePeriodById(controlDto.ActivePeriod.ID)
		if err != nil {
			tx.Rollback()
			utils.JsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse Start and End times
		var startTime, endTime time.Duration
		var validStart, validEnd bool

		if controlDto.ActivePeriod.Start != "" {
			t, err := time.Parse("15:04", controlDto.ActivePeriod.Start)
			if err != nil {
				tx.Rollback()
				utils.JsonError(w, fmt.Sprintf("Invalid start time format: %s", err.Error()), http.StatusBadRequest)
				return
			}
			startTime = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
			validStart = true
		}

		if controlDto.ActivePeriod.End != "" {
			t, err := time.Parse("15:04", controlDto.ActivePeriod.End)
			if err != nil {
				tx.Rollback()
				utils.JsonError(w, fmt.Sprintf("Invalid end time format: %s", err.Error()), http.StatusBadRequest)
				return
			}
			endTime = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
			validEnd = true
		}

		// Update ActivePeriod with parsed times if valid
		if validStart {
			activePeriodUpdate.Start = startTime
		}

		if validEnd {
			activePeriodUpdate.End = endTime
		}

		// Handle updating days
		if len(controlDto.ActivePeriod.Days) > 0 {
			var daysBitmask uint8
			for _, day := range controlDto.ActivePeriod.Days {
				daysBitmask |= 1 << (day - 1) // Shift 1 to the left by day-1 places
			}
			activePeriodUpdate.Days = daysBitmask
		}

		if err := tx.Save(&activePeriodUpdate).Error; err != nil {
			tx.Rollback()
			utils.JsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}


		sensor, err := findSensorById(string(controlDto.Condition.DependentSensor.ID))
		if err != nil {
			utils.JsonError(w, fmt.Sprintf("Invalid end time format: %s", err.Error()), http.StatusBadRequest)
			return
		}

		conditionUpdate := models.Condition{
			On: controlDto.Condition.On,
			Off: controlDto.Condition.Off,
			DependentSensor: sensor,
		}
	

		// Update main control settings
		controlUpdate := models.Control{
			Alias:        controlDto.Alias,
			Description:  controlDto.Description,
			Condition: 	 &conditionUpdate,
		}

		controlSettingsDbObject, err := findControllSettingById(controlDto.ID)
		if err != nil {
			tx.Rollback()
			utils.JsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update control settings
		if err := tx.Model(&controlSettingsDbObject).Updates(controlUpdate).Error; err != nil {
			tx.Rollback()
			utils.JsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)	
	json.NewEncoder(w).Encode(controlDtos)
}

// findActivePeriodById fetches an ActivePeriod by its ID
func findActivePeriodById(id uint) (*models.ActivePeriod, error) {
	var activePeriod models.ActivePeriod
	result := initPackage.Db.First(&activePeriod, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &activePeriod, nil
}

// findControllSettingById fetches a Control by its ID
func findControllSettingById(id uint) (*models.Control, error) {
	var control models.Control
	result := initPackage.Db.First(&control, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &control, nil
}
