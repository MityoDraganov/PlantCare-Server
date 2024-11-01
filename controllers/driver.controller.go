package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
)

// EditDriver - Update an existing driver by ID
func GetAllDrivers(w http.ResponseWriter, r *http.Request) {
	// Define a slice to hold multiple drivers
	var drivers []models.Driver

	// Query all drivers from the database
	result := initPackage.Db.Find(&drivers)
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// If no drivers are found, return a 404 error
	if result.RowsAffected == 0 {
		utils.JsonError(w, "No drivers found", http.StatusNotFound)
		return
	}

	var driverDtos []dtos.DriverDto
	for _, driver := range drivers {
		driverDtos = append(driverDtos, ToDriverDTO(driver))
	}

	// Set content type and encode the result to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(driverDtos)
}

// UploadDriver - Create a new driver and upload image to Firebase
func UploadDriver(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		utils.JsonError(w, "Unauthorized: unable to extract session claims", http.StatusUnauthorized)
		return
	}
	clerkUserID := claims.Subject

	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		utils.JsonError(w, "File too large", http.StatusBadRequest)
		return
	}

	// Retrieve fields
	alias := r.FormValue("alias")
	downloadUrl := r.FormValue("downloadUrl")

	// Retrieve file from form data
	file, header, err := r.FormFile("marketplaceBanner")
	if err != nil {
		utils.JsonError(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Define destination path for Firebase (e.g., "drivers/{filename}")
	destinationPath := "drivers/" + header.Filename

	// Upload file to Firebase
	imageUrl, err := utils.UploadFile(file, destinationPath)
	if err != nil {
		fmt.Println(err)
		utils.JsonError(w, "Failed to upload image to Firebase", http.StatusInternalServerError)
		return
	}

	// Create a new Driver record
	driver := models.Driver{
		Alias:                alias,
		DownloadUrl:          downloadUrl,
		MarketplaceBannerUrl: &imageUrl, // Save Firebase URL
		UploadedByUserID:     clerkUserID,
	}

	// Insert driver record in database
	result := initPackage.Db.Create(&driver)
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Send response back with created driver details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(driver)
}

// EditDriver - Update an existing driver by ID
func EditDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["driverId"]

	var updateDto dtos.AddWebhookDto
	err := json.NewDecoder(r.Body).Decode(&updateDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var driverDbObject models.Driver
	result := initPackage.Db.Model(&driverDbObject).Where("id = ?", id).Updates(updateDto).Clauses(clause.Returning{})
	if result.Error != nil {
		if result.RowsAffected == 0 {
			utils.JsonError(w, "Driver not found", http.StatusNotFound)
			return
		}
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(driverDbObject)
}

// DeleteDriver - Delete a driver by ID
func DeleteDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["driverId"]

	var driverDbObject models.Driver
	result := initPackage.Db.Model(&driverDbObject).Where("id = ?", id).Delete(&driverDbObject)
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		utils.JsonError(w, "Driver not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToDriverDTO - Map driver model to DTO
func ToDriverDTO(driver models.Driver) dtos.DriverDto {
	return dtos.DriverDto{
		DownloadUrl:          driver.DownloadUrl,
		MarketplaceBannerUrl: driver.MarketplaceBannerUrl,
		Alias:                driver.Alias,
	}
}
