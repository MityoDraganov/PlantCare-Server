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

	for _ ,controlDto := range controlDtos {
		controlSettingsDbObject, err := findControllSettingById(controlDto.ID)
		if err != nil {
			utils.JsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		initPackage.Db.Model(&controlSettingsDbObject).Updates(controlDto)
		initPackage.Db.Model(&controlSettingsDbObject).Updates(controlDto)
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
}


// general

func findControllSettingById(id uint) (*models.Control, error) {
	var control models.Control
	result := initPackage.Db.First(&control, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &control, nil
}