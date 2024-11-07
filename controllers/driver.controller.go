package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
)

// EditDriver - Update an existing driver by ID
func GetAllDrivers(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		utils.JsonError(w, "Unauthorized: unable to extract session claims", http.StatusUnauthorized)
		return
	}
	clerkUserID := claims.Subject

	var drivers []models.Driver

	// Query all drivers from the database
	result := initPackage.Db.Where("is_marketplace_featured = ?", true).Find(&drivers)
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
		isClerkUser := driver.UploadedByUserID == clerkUserID
		driverDtos = append(driverDtos, ToDriverDTO(driver, isClerkUser))
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
		UploadedByUserID:          clerkUserID,
		IsMarketplaceFeatured: true,
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

// EditDriver - Update an existing driver by ID, including optional file update
func EditDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["driverId"]

	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		utils.JsonError(w, "Unauthorized: unable to extract session claims", http.StatusUnauthorized)
		return
	}
	clerkUserID := claims.Subject

	// Parse form data including file if available
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		utils.JsonError(w, "File too large", http.StatusBadRequest)
		return
	}

	// Get form values for fields to update
	alias := r.FormValue("alias")
	downloadUrl := r.FormValue("downloadUrl")

	var updatedFields = map[string]interface{}{
		"Alias":       alias,
		"DownloadUrl": downloadUrl,
	}

	// Check if a new file is included in the request
	file, header, err := r.FormFile("marketplaceBanner")
	if err == nil { // File is provided
		defer file.Close()
		destinationPath := "drivers/" + header.Filename

		// Upload the new file to Firebase
		imageUrl, err := utils.UploadFile(file, destinationPath)
		if err != nil {
			utils.JsonError(w, "Failed to upload image to Firebase", http.StatusInternalServerError)
			return
		}
		updatedFields["MarketplaceBannerUrl"] = imageUrl
	}

	// Fetch the existing driver record
	var driverDbObject models.Driver
	result := initPackage.Db.First(&driverDbObject, "id = ?", id)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			utils.JsonError(w, "Driver not found", http.StatusNotFound)
			return
		}
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Ensure the clerk user is the original uploader
	if driverDbObject.UploadedByUserID != clerkUserID {
		utils.JsonError(w, "Unauthorized: you do not have permission to edit this driver", http.StatusForbidden)
		return
	}

	// Update the driver record with new fields and optionally new file URL
	result = initPackage.Db.Model(&driverDbObject).Where("id = ?", id).Updates(updatedFields).Clauses(clause.Returning{})
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the updated driver information
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

type ClerkUserResponse struct {
	Username       string `json:"username"`
	PrimaryEmailID string `json:"primary_email_address_id"`
	EmailAddresses []struct {
		ID           string `json:"id"`
		EmailAddress string `json:"email_address"`
	} `json:"email_addresses"`
}

func ToDriverDTO(driver models.Driver, IsUploader bool) dtos.DriverDto {
	url := fmt.Sprintf("https://api.clerk.dev/v1/users/%s", driver.UploadedByUserID)

	// Create a new HTTP request with Clerk API Key for authorization
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CLERK_API_KEY"))

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Parse the JSON response into our struct
	var user ClerkUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Errorf("failed to parse response: %w", err)
	}

	var primaryEmail string
	for _, email := range user.EmailAddresses {
		if email.ID == user.PrimaryEmailID {
			primaryEmail = email.EmailAddress
			break
		}
	}

	userDto := dtos.UserResponseDto{
		Username: user.Username,
		Email:    primaryEmail,
	}


	
	return dtos.DriverDto{
		Id:                   driver.ID,
		DownloadUrl:          driver.DownloadUrl,
		MarketplaceBannerUrl: *utils.CoalesceString(driver.MarketplaceBannerUrl),
		Alias:                driver.Alias,
		User:                 userDto,
		IsUploader:           IsUploader,
	}
}
