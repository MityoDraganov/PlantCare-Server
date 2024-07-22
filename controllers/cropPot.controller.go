package controllers

import (
	"encoding/json"
	"net/http"

	"PlantCare/dtos"
	"PlantCare/middlewares"
	"PlantCare/models"
	"PlantCare/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
)

//cropPotDBObject

func GetCropPotsForUser(w http.ResponseWriter, r *http.Request) {
	var cropPots []dtos.CropPotResponse
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(*utils.CustomClaims)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//db.Where("user_id = ?", claims.UserID).Find(&cropPots)
	db.Model(&models.CropPot{}).Where("user_id = ?", claims.UserID).
		Select("id, alias, watering_interval, last_watered_at, is_archived").
		Find(&cropPots)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPots)
}

func AddCropPot(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(middlewares.ClaimsKey).(*utils.CustomClaims)
	println("userID")
	println(claims.UserID)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var cropPotDto dtos.CreateCropPot
	err := json.NewDecoder(r.Body).Decode(&cropPotDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(cropPotDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cropPot := models.CropPot{
		Alias:            cropPotDto.Alias,
		WateringInterval: cropPotDto.WateringInterval,
		UserID:           claims.UserID,
	}

	cropPotDBObject := db.Create(&cropPot)
	if cropPotDBObject.Error != nil {
		http.Error(w, cropPotDBObject.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPot)
}

func UpdateCropPot(w http.ResponseWriter, r *http.Request) {
	var cropPotDto dtos.CreateCropPot
	err := json.NewDecoder(r.Body).Decode(&cropPotDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	cropPotDBObject, err := findCropPotById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db.Model(&cropPotDBObject).Clauses(clause.Returning{}).Updates(cropPotDto)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPotDBObject)
}

func RemoveCropPot(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	cropPotDBObject, err := findCropPotById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cropPotDBObject.IsArchived = true

	// Save the updated object back to the database
	result := db.Save(&cropPotDBObject)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

func findCropPotById(id string) (*models.CropPot, error) {
	var cropPot models.CropPot
	result := db.First(&cropPot, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cropPot, nil
}
