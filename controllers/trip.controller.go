package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"TravelBuddy/middlewares"
	"TravelBuddy/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

func CreateTrip(w http.ResponseWriter, r *http.Request) {
	var trip models.Trip
	err := json.NewDecoder(r.Body).Decode(&trip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the driver's ID from the request context
	claims, ok := r.Context().Value(middlewares.UserContextKey).(*jwt.RegisteredClaims)
	if !ok {
		http.Error(w, "Unable to retrieve claims", http.StatusInternalServerError)
		return
	}
	driverID := claims.Subject

	// Convert driverID to uint and assign to trip
	driverIDUint64, err := strconv.ParseUint(driverID, 10, 64)
	println(driverIDUint64, driverID)
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}
	trip.DriverID = uint(driverIDUint64)

	// Validate the trip
	err = validate.Struct(trip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the trip to the database
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

	trip, err := findTripById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var updatedTrip models.Trip
	err = json.NewDecoder(r.Body).Decode(&updatedTrip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Model(&trip).Updates(updatedTrip)
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

func addPassenger(tripId string, Passenger models.Passenger) (error){
	trip,_ := findTripById(tripId)

	trip.CurrentPassengers = append(trip.CurrentPassengers, Passenger)
	result := db.Save(&trip)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func findTripById(id string) (*models.Trip, error) {
	var trip models.Trip
	result := db.First(&trip, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &trip, nil
}
