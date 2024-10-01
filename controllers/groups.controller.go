package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"net/http"
)


func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var groupDto dtos.GroupRequestDto
	err := json.NewDecoder(r.Body).Decode(&groupDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	cropPotDbObject, err := FindCropPotById(string(groupDto.CropPotID))

	group := models.Group{
		CropPots: []models.CropPot{*cropPotDbObject},
	}

	result := initPackage.Db.Create(&group)
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	GroupResponse := dtos.GroupResponsetDto{
		CropPots: ToCropPotResponse(group.CropPots),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GroupResponse)
}