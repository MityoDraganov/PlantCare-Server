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
    var controlsDto []dtos.ControlRequestDto
    err := json.NewDecoder(r.Body).Decode(&controlsDto)
    if err != nil {
        utils.JsonError(w, err.Error(), http.StatusBadRequest)
        return
    }

    tx := initPackage.Db.Begin()
    if tx.Error != nil {
        utils.JsonError(w, tx.Error.Error(), http.StatusInternalServerError)
        return
    }

    for _, controlDto := range controlsDto {
        activePeriodUpdate, err := findActivePeriodById(controlDto.ActivePeriod.ID)
        if err != nil {
            tx.Rollback()
            utils.JsonError(w, err.Error(), http.StatusInternalServerError)
            return
        }

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

        if validStart {
            activePeriodUpdate.Start = startTime
        }

        if validEnd {
            activePeriodUpdate.End = endTime
        }

        if len(controlDto.ActivePeriod.Days) > 0 {
            var daysBitmask uint8
            for _, day := range controlDto.ActivePeriod.Days {
                daysBitmask |= 1 << (day - 1)
            }
            activePeriodUpdate.Days = daysBitmask
        }

        if err := tx.Save(&activePeriodUpdate).Error; err != nil {
            tx.Rollback()
            utils.JsonError(w, err.Error(), http.StatusInternalServerError)
            return
        }

        var sensor *models.Sensor
        if controlDto.Condition.DependentSensor != nil {
            sensor, err = findSensorById(controlDto.Condition.DependentSensor.ID)
            if err != nil {
                tx.Rollback()
                utils.JsonError(w, fmt.Sprintf("Failed to find sensor: %s", err.Error()), http.StatusBadRequest)
                return
            }
        }

		fmt.Println(sensor)

        // conditionUpdate := models.Condition{
        //     On:                   controlDto.Condition.On,
        //     Off:                  controlDto.Condition.Off,
        //     DependentSensorID:    &sensor.ID,
        // }

        controlUpdate := models.Control{
            Alias:        controlDto.Alias,
            Description:  controlDto.Description,
           // Condition:    &conditionUpdate,
        }

        controlSettingsDbObject, err := findControlById(controlDto.ID)
        if err != nil {
            tx.Rollback()
            utils.JsonError(w, err.Error(), http.StatusInternalServerError)
            return
        }

        if err := tx.Model(&controlSettingsDbObject).Updates(controlUpdate).Error; err != nil {
            tx.Rollback()
            utils.JsonError(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    if err := tx.Commit().Error; err != nil {
        utils.JsonError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(controlsDto)
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
func findControlById(id uint) (*models.Control, error) {
	var control models.Control
	result := initPackage.Db.First(&control, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &control, nil
}

// Map Control(s) to ControlDto(s)
func ToControlsDTO(input interface{}) []dtos.ControlDto {
	switch v := input.(type) {
	case models.Control:
		// If it's a single control, map to ControlDto
		return []dtos.ControlDto{mapControlToDTO(v)}
	case []models.Control:
		// If it's a slice of controls, map each control to ControlDto
		controlDTOs := make([]dtos.ControlDto, len(v))
		for i, control := range v {
			controlDTOs[i] = mapControlToDTO(control)
		}
		return controlDTOs
	default:
		// Handle unexpected types
		return nil
	}
}

// Helper function to map a single control to ControlDto
func mapControlToDTO(control models.Control) dtos.ControlDto {


	startStr := utils.DurationToTimeString(control.ActivePeriod.Start)
	endStr := utils.DurationToTimeString(control.ActivePeriod.End)

	activePeriod := dtos.ActivePeriod{
		ID:    control.ActivePeriod.ID,
		Start: startStr,
		End:   endStr,
		Days:  utils.ParseBitmask(control.ActivePeriod.Days),
	}

	return dtos.ControlDto{
		ID:           control.ID,
		SerialNumber: control.SerialNumber,
		Alias:        control.Alias,
		Description:  utils.CoalesceString(control.Description),
		Updates:      control.Updates,
		IsOfficial:   true, // Set this as per your business logic
		// Condition: &dtos.ConditionDto{
		// 	On:  control.Condition.On,
		// 	Off: control.Condition.Off,
		// 	DependentSensor: func() *dtos.SensorDto {
		// 		if control.Condition.DependentSensor != nil {
		// 			dto := MapSensorToDTO(*control.Condition.DependentSensor)
		// 			return &dto
		// 		}
		// 		return nil
		// 	}(),
		// },
		ActivePeriod: activePeriod,
	}

}
