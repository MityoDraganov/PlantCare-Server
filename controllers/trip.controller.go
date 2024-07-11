package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"TravelBuddy/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var tripCollection *mongo.Collection
var validate *validator.Validate

func SetTripCollection(collection *mongo.Collection) {
	tripCollection = collection
	validate = validator.New()
}

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

	_, err = tripCollection.InsertOne(context.TODO(), trip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trip)
}

func UpdateTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updatedTrip models.Trip
	err := json.NewDecoder(r.Body).Decode(&updatedTrip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": updatedTrip,
	}

	result, err := tripCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Trip not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTrip)
}

func DeleteTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	filter := bson.M{"id": id}

	result, err := tripCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Trip not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}
