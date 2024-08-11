package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

func UpdateControllSettings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["controllSettingsId"]

	// Decode the incoming JSON request into the DTO
	var controlSettingsDto dtos.ControlSettingsDto
	err := json.NewDecoder(r.Body).Decode(&controlSettingsDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the existing ControlSettings object in the database
	controlSettingsDbObject, err := findControllSettingsById(id)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	initPackage.Db.Model(&controlSettingsDbObject).Updates(controlSettingsDto)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(controlSettingsDbObject)
}


// general

func findControllSettingsById(id string) (*models.ControlSettings, error) {
	var controllSettings models.ControlSettings
	result := initPackage.Db.First(&controllSettings, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &controllSettings, nil
}