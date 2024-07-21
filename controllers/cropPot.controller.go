package controllers

import (
	"encoding/json"
	"net/http"

	"PlantCare/dtos"
	"PlantCare/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
)

//cropPotDBObject

func AddCropPot(w http.ResponseWriter, r *http.Request) {
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
		Alias: cropPotDto.Alias,
		WateringInterval: cropPotDto.WateringInterval,
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

// func RemoveCroPot(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
// }

// func addCropPot(cropPotId string, Passenger models.Passenger) (error){
// 	cropPot,_ := findCropPotById(cropPotId)

// 	cropPot.CurrentPassengers = append(cropPot.CurrentPassengers, Passenger)
// 	result := db.Save(&cropPot)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

func findCropPotById(id string) (*models.CropPot, error) {
	var cropPot models.CropPot
	result := db.First(&cropPot, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cropPot, nil
}
