package controllers

import (
	"encoding/json"
	"net/http"

	"PlantCare/dtos"
	"PlantCare/utils"

	"PlantCare/models"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
)

//cropPotDBObject

func GetCropPotsForUser(w http.ResponseWriter, r *http.Request) {
	var cropPots []dtos.CropPotResponse
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
		return
	}


	//db.Where("user_id = ?", claims.UserID).Find(&cropPots)
	db.Model(&models.CropPot{}).Where("user_id = ?", claims.ID).
		Select("id, alias, watering_interval, last_watered_at, is_archived").
		Find(&cropPots)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPots)
}

func AssignCropPotToUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
		return
	}


	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cropPotDBObject, err := findCropPotById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cropPotDBObject.ClerkUserID = &claims.ID

	db.Save(cropPotDBObject)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPotDBObject)
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

//admin action
func AddCropPot(w http.ResponseWriter, r *http.Request) {

	//claims, ok := clerk.SessionClaimsFromContext(r.Context())
	// if !ok {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	w.Write([]byte(`{"error": "unauthorized"}`))
	// 	return
	// }

	token, err := utils.GenerateSecureToken(32)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cropPot := models.CropPot{
		Token: token,
	}

	cropPotDBObject := db.Create(&cropPot)
	if cropPotDBObject.Error != nil {
		http.Error(w, cropPotDBObject.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPot)
}
