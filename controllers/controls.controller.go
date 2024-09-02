package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"net/http"
)

func UpdateControllSetting(w http.ResponseWriter, r *http.Request) {
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
        activePeriodUpdate,err := findActivePeriodById(controlDto.ActivePeriod.ID)
        if err != nil {
            tx.Rollback()
			utils.JsonError(w, err.Error(), http.StatusInternalServerError)
			return
        }

        activePeriodUpdate.Start = controlDto.ActivePeriod.Start
        activePeriodUpdate.End = controlDto.ActivePeriod.End
        if err := tx.Save(&activePeriodUpdate).Error; err != nil {
			tx.Rollback()
			utils.JsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}


		controlUpdate := models.Control{
			Alias:        controlDto.Alias,
			Description:  controlDto.Description,

			OnCondition:  controlDto.OnCondition,
			OffCondition: controlDto.OffCondition,
		}



		controlSettingsDbObject, err := findControllSettingById(controlDto.ID)
		if err != nil {
			tx.Rollback()
			utils.JsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update main control settings
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
}

func findActivePeriodById(id uint) (*models.ActivePeriod, error) {
	var activePeriod models.ActivePeriod
	result := initPackage.Db.First(&activePeriod, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &activePeriod, nil
}

func findControllSettingById(id uint) (*models.Control, error) {
	var control models.Control
	result := initPackage.Db.First(&control, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &control, nil
}
