package controllers

import (
	"encoding/json"
	"net/http"

	"TravelBuddy/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateTrip(w http.ResponseWriter, r *http.Request) {
	var trip models.Trip
	err := json.NewDecoder(r.Body).Decode(&trip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(trip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&trip)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trip)
}

func UpdateTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var trip models.Trip
	result := findTripById(id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	var updatedTrip models.Trip
	err := json.NewDecoder(r.Body).Decode(&updatedTrip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result = db.Model(&trip).Updates(updatedTrip)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTrip)
}

func DeleteTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result := db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

func findTripById(id string) *gorm.DB {
	var trip models.Trip
	result := db.First(&trip, "id = ?", id)

	return result
}
