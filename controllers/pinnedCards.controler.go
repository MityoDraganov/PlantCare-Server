package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm/clause"
)

func CreateCard(w http.ResponseWriter, r *http.Request) {
	var cardDto dtos.PinnedCard

	err := json.NewDecoder(r.Body).Decode(&cardDto)
	if err != nil {
		log.Println(err)
		utils.JsonError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	result := initPackage.Db.Model(&models.PinnedCard{}).Create(cardDto)
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated array of sensorDtos
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateCard updates an existing PinnedCard based on its ID.
func UpdateCard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id") // Assuming the ID is passed as a query parameter.
	var cardDto dtos.PinnedCard

	// Decode the request body into the DTO.
	err := json.NewDecoder(r.Body).Decode(&cardDto)
	if err != nil {
		log.Println(err)
		utils.JsonError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	var cardDbObject models.PinnedCard
	if err := initPackage.Db.First(&cardDbObject, id).Error; err != nil {
		utils.JsonError(w, "Card not found", http.StatusNotFound)
		return
	}

	result := initPackage.Db.Model(cardDbObject).Updates(cardDto).Clauses(clause.Returning{})
	if result.Error != nil {
		log.Printf("Failed to update sensor: %v", result.Error)
		utils.JsonError(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cardDbObject); err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteCard deletes a PinnedCard based on its ID.
func DeleteCard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id") // Assuming the ID is passed as a query parameter.

	// Delete the pinned card in the database.
	if err := initPackage.Db.Delete(&models.PinnedCard{}, id).Error; err != nil {
		utils.JsonError(w, "Card not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Respond with no content status.
}
