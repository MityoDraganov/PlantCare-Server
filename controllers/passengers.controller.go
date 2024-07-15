package controllers

import (
	"TravelBuddy/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RequestTrip(w http.ResponseWriter, r *http.Request) {
	var passenger models.Passenger
	err := json.NewDecoder(r.Body).Decode(&passenger)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(passenger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passenger.State = models.StatePending

	result := db.Create(&passenger)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passenger)
}

func DenyTrip(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	passenger, err := findPassengerById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Add logic to handle denying the trip, for example:
	passenger.State = models.StateDenied
	result := db.Save(&passenger)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func AcceptPassenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripId := params["tripId"]
	passengerId := params["passengerId"]

	passenger, err := findPassengerById(passengerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	
	passenger.State = models.StateActive
	result := db.Save(&passenger)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	addPassenger(tripId, *passenger)

	w.WriteHeader(http.StatusOK)
}

func findPassengerById(id string) (*models.Passenger, error) {
	var passenger models.Passenger
	result := db.First(&passenger, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &passenger, nil
}
