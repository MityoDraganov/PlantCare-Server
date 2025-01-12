package controllers

import (
	"PlantCare/services"
	"PlantCare/utils"
	"PlantCare/utils/firebaseUtil"
	"encoding/json"
	"fmt"
	"net/http"
)

func DiagnoseMeasuringGroup(w http.ResponseWriter, r *http.Request) {

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max file size
	if err != nil {
		http.Error(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form
	file, header, err := r.FormFile("picture")
	if err != nil {
		utils.JsonError(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Define destination path for Firebase (e.g., "drivers/{filename}")
	destinationPath := "plantPictures/diagnosisPictures/" + header.Filename

	plantHealth, err := services.PredictPlantHealth(file)

	// Upload file to Firebase
	imageUrl, err := firebaseUtil.UploadFile(file, destinationPath)
	if err != nil {
		fmt.Println(err)
		utils.JsonError(w, "Failed to upload image to Firebase", http.StatusInternalServerError)

		// Respond with success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	response := struct {
		PlantHealth *string `json:"plantHealth"`
		ImageUrl string `json:"imageUrl"`
	}{
		PlantHealth: plantHealth,
		ImageUrl: imageUrl,
	}

	// Respond with the image URL
	w.Header().Set("Content-Type", "application/json")
	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.JsonError(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(responseBytes)
}
