package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/utils"

	"PlantCare/models"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SetAllPotsOffline() error {
	db := initPackage.Db
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Model(&models.CropPot{}).
		Update("status", models.StatusOffline).Error; err != nil {
		return fmt.Errorf("failed to set all pots to offline: %w", err)
	}
	return nil
}

func GetCropPotsForUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
		return
	}
	cropPots, err := FindPotsByUserId(claims.Subject)
	if err != nil {
		fmt.Println("Error extracting session claims")
	}

	cropPotResponses := ToCropPotResponse(cropPots)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPotResponses)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPotDBObject)
}

func UpdateCropPot(w http.ResponseWriter, r *http.Request) {
	var cropPotDto dtos.CropPotRequest

	err := json.NewDecoder(r.Body).Decode(&cropPotDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["potId"]

	// Find the crop pot by ID
	cropPotDBObject, err := FindCropPotById(id)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var intervalUpdateTime time.Duration
	if cropPotDto.MeasurementInterval != "" {
		t, err := time.Parse("15:04", cropPotDto.MeasurementInterval)
		if err != nil {
			log.Printf("Invalid measurement interval format: %s", err.Error())
			utils.JsonError(w, fmt.Sprintf("Invalid measurement interval format: %s", err.Error()), http.StatusBadRequest)
			return
		}
		intervalUpdateTime = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
	}

	potUpdate := models.CropPot{
		Alias:              cropPotDto.Alias,
		IsPinned:           cropPotDto.IsPinned,
		MeasuremntInterval: intervalUpdateTime,
	}


	// Update CropPot with the new values
	initPackage.Db.Model(&cropPotDBObject).Updates(potUpdate)

	// Check if isPinned is set to true and if there is no existing Canvas for this CropPot
	if cropPotDto.IsPinned && cropPotDBObject.Canvas.ID == 0 {
		// Create a new Canvas for the CropPot
		newCanvas := models.Canvas{
			CropPotID: cropPotDBObject.ID,
		}

		// Save the new Canvas to the database
		if err := initPackage.Db.Create(&newCanvas).Error; err != nil {
			utils.JsonError(w, "Failed to create Canvas", http.StatusInternalServerError)
			return
		}
		cropPotDBObject.Canvas = newCanvas
	}

    if !cropPotDto.IsPinned {
        if cropPotDBObject.Canvas.ID != 0 {
			fmt.Println("Deleted canvas")
           DeleteCanvas(cropPotDBObject.Canvas.ID )
        }
        // Ensure the IsPinned flag is updated
        cropPotDBObject.IsPinned = false
    }


	if err := initPackage.Db.Save(&cropPotDBObject).Clauses(clause.Returning{}).Error; err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPotDBObject)
}

func RemoveCropPot(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["potId"]

	cropPotDBObject, err := FindCropPotById(id)
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

func FindCropPotById(id string) (*models.CropPot, error) {

	var cropPot models.CropPot
	result := initPackage.Db.
		Preload("Sensors").
		Preload("Sensors.Measurements").
		Preload("Canvas").
		Preload("Canvas.PinnedCards").
		First(&cropPot, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cropPot, nil
}

func FindPotByToken(token string) (*models.CropPot, error) {
	var cropPot models.CropPot
	if err := initPackage.Db.
		Where("token = ?", token).First(&cropPot).Error; err != nil {
		return nil, err
	}
	return &cropPot, nil
}

func FindPotsByUserId(userId string) ([]models.CropPot, error) {
	var cropPots []models.CropPot
	result := initPackage.Db.
		Preload("Canvas").
		Preload("Canvas.PinnedCards").
		Preload("Sensors").
		Preload("Sensors.Measurements").
		Preload("Sensors.Driver").
		Preload("Controls").
		Preload("Webhooks").
		Preload("Webhooks.SubscribedEvents").
		Where("clerk_user_id = ?", userId).
		Find(&cropPots)
	if result.Error != nil {
		return nil, result.Error
	}

	return cropPots, nil
}

func ToCropPotResponse(data interface{}) []dtos.CropPotResponse {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Slice {
		var dtosArray []dtos.CropPotResponse
		for i := 0; i < val.Len(); i++ {
			cropPot := val.Index(i).Interface().(models.CropPot)
			dtosArray = append(dtosArray, ToCropPotResponseDTO(cropPot))
		}
		return dtosArray
	} else if val.Kind() == reflect.Struct {
		cropPot := val.Interface().(models.CropPot)
		return []dtos.CropPotResponse{ToCropPotResponseDTO(cropPot)}
	}
	return nil
}

func ToCropPotResponseDTO(cropPot models.CropPot) dtos.CropPotResponse {
	return dtos.CropPotResponse{
		ID:                  cropPot.ID,
		Alias:               cropPot.Alias,
		IsArchived:          cropPot.IsArchived,
		IsPinned:            cropPot.IsPinned,
		Controls:            ToControlsDTO(cropPot.Controls),
		Sensors:             ToSensorsDTO(cropPot.Sensors),
		Webhooks:            ToWebhooksDTO(cropPot.Webhooks),
		Status:              cropPot.Status,
		MeasurementInterval: utils.DurationToTimeString(cropPot.MeasuremntInterval),
		Canvas:              ToCanvasDTO(cropPot.Canvas),
	}
}
