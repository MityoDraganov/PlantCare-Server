package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/utils"

	"PlantCare/models"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gorilla/mux"


)

type CustomClaims struct {
	UserID         string `json:"user_id"`
	ExternalID     string `json:"external_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	FullName       string `json:"full_name"`
	Username       string `json:"username"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	PrimaryEmail   string `json:"primary_email_address"`
	PrimaryPhone   string `json:"primary_phone_number"`
	PrimaryWeb3    string `json:"primary_web3_wallet"`
	EmailVerified  bool   `json:"email_verified"`
	PhoneVerified  bool   `json:"phone_number_verified"`
	ImageURL       string `json:"image_url"`
	HasImage       bool   `json:"has_image"`
	TwoFactor      bool   `json:"two_factor_enabled"`
	PublicMetadata string `json:"public_metadata"`
	UnsafeMetadata string `json:"unsafe_metadata"`
	SessionActor   string `json:"session_actor"`
}

func GetCropPotsForUser(w http.ResponseWriter, r *http.Request) {
	var cropPots []dtos.CropPotResponse
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
		return
	}

	//db.Where("user_id = ?", claims.UserID).Find(&cropPots)
	initPackage.Db.Model(&models.CropPot{}).Where("user_id = ?", claims.ID).
		Select("id, alias, watering_interval, last_watered_at, is_archived").
		Find(&cropPots)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPots)
}

func AssignCropPotToUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		fmt.Println("Error extracting session claims")
		utils.JsonError(w, "Unauthorized: unable to extract session claims", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)

	cropPotDBObject, err := FindPotByToken(params["token"])
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cropPotDBObject.ClerkUserID != nil {
		utils.JsonError(w, "Crop pot already assigned! Contact support for more information.", http.StatusUnauthorized)
		return
	}

	clerkUserID := claims.Subject
	cropPotDBObject.ClerkUserID = &clerkUserID

	result := initPackage.Db.Save(cropPotDBObject)
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("CropPot after update:", cropPotDBObject)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPotDBObject)
}

func UpdateCropPot(w http.ResponseWriter, r *http.Request) {
	var cropPotDto dtos.CreateCropPot
	err := json.NewDecoder(r.Body).Decode(&cropPotDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	// Find existing CropPot by ID
	cropPotDBObject, err := findCropPotById(id)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the ControlSettings needs to be updated or created
	if cropPotDto.ControlSettings != nil {
		var controlSettings models.ControlSettings

		if cropPotDBObject.ControlSettingsID != nil {
			// Update existing ControlSettings
			if err := initPackage.Db.First(&controlSettings, *cropPotDBObject.ControlSettingsID).Error; err != nil {
				utils.JsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Update the ControlSettings
			initPackage.Db.Model(&controlSettings).Updates(cropPotDto.ControlSettings)
		} else {
			// Create new ControlSettings
			controlSettings = models.ControlSettings{
				WateringInterval: cropPotDto.ControlSettings.WateringInterval,
				LastWateredAt:    cropPotDto.ControlSettings.LastWateredAt,
				CropPotID:        cropPotDBObject.ID,
			}

			if err := initPackage.Db.Create(&controlSettings).Error; err != nil {
				utils.JsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Update CropPot with the new ControlSettingsID
			cropPotDBObject.ControlSettingsID = &controlSettings.ID
		}
	}

	// Update the CropPot
	if err := initPackage.Db.Model(&cropPotDBObject).Updates(cropPotDto).Error; err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPotDBObject)
}


func RemoveCropPot(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	cropPotDBObject, err := findCropPotById(id)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cropPotDBObject.IsArchived = true

	// Save the updated object back to the database
	result := initPackage.Db.Save(&cropPotDBObject)
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// admin action
func AddCropPot(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GenerateSecureToken(32)

	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	cropPot := models.CropPot{
		Token: token,
	}

	cropPotDBObject := initPackage.Db.Create(&cropPot)
	if cropPotDBObject.Error != nil {
		utils.JsonError(w, cropPotDBObject.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPot)
}

func findCropPotById(id string) (*models.CropPot, error) {
	var cropPot models.CropPot
	result := initPackage.Db.First(&cropPot, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cropPot, nil
}

func FindPotByToken(token string) (*models.CropPot, error) {
	var cropPot models.CropPot
	if err := initPackage.Db.Where("token = ?", token).First(&cropPot).Error; err != nil {
		return nil, err
	}
	return &cropPot, nil
}
