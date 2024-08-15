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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetCropPotsForUser(w http.ResponseWriter, r *http.Request) {
    claims, ok := clerk.SessionClaimsFromContext(r.Context())
    if !ok {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte(`{"error": "unauthorized"}`))
        return
    }
    cropPots, err := findPotsByUserId(claims.Subject)
    if err != nil {
        fmt.Println("Error extracting session claims")
        utils.JsonError(w, "Pots not found!", http.StatusNotFound)
        return
    }

    // Map the crop pots to the DTOs
    var cropPotResponses []dtos.CropPotResponse
    for _, cropPot := range cropPots {
        var controlSettingsResponse *dtos.ControlSettingsResponse
        if cropPot.ControlSettings != nil {
            controlSettingsResponse = &dtos.ControlSettingsResponse{
                WateringInterval: cropPot.ControlSettings.WateringInterval,
            }
        }

        // Map SensorData
        var sensorDataResponses []dtos.SensorDataResponse
        for _, sensorData := range cropPot.SensorDatas {
            sensorDataResponses = append(sensorDataResponses, dtos.SensorDataResponse{
				CreatedAt: sensorData.CreatedAt,
                Temperature: sensorData.Temperature,
                Moisture:    sensorData.Moisture,
                WaterLevel:  sensorData.WaterLevel,
                SunExposure: sensorData.SunExposure,
            })
        }

        // Map CustomSensorData
        var customSensorDataResponses []dtos.CustomSensorDataResponse
        for _, customSensorField := range cropPot.CustomSensorFields {
            for _, customSensorData := range customSensorField.CustomSensorData {
                customSensorDataResponses = append(customSensorDataResponses, dtos.CustomSensorDataResponse{
                    FieldAlias: customSensorField.FieldAlias,
                    DataValue:  customSensorData.DataValue,
                })
            }
        }

        cropPotResponse := dtos.CropPotResponse{
            ID:                cropPot.ID,
            Alias:             cropPot.Alias,
            LastWateredAt:     cropPot.LastWateredAt,
            IsArchived:        cropPot.IsArchived,
            ControlSettings:   controlSettingsResponse,
            SensorData:        sensorDataResponses,
            CustomSensorData:  customSensorDataResponses,
        }
        cropPotResponses = append(cropPotResponses, cropPotResponse)
    }

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

	fmt.Println("CropPot after update:", cropPotDBObject)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPotDBObject)
}

func UpdateCropPot(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		fmt.Println("Error extracting session claims")
		utils.JsonError(w, "Unauthorized: unable to extract session claims", http.StatusUnauthorized)
		return
	}

	var cropPotDto dtos.CreateCropPot
	err := json.NewDecoder(r.Body).Decode(&cropPotDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	id := params["id"]

	cropPotDBObject, err := FindCropPotById(id, claims.Subject)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	initPackage.Db.Model(&cropPotDBObject).Clauses(clause.Returning{}).Updates(cropPotDto)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cropPotDBObject)
}

func RemoveCropPot(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		fmt.Println("Error extracting session claims")
		utils.JsonError(w, "Unauthorized: unable to extract session claims", http.StatusUnauthorized)
		return
	}
	
	params := mux.Vars(r)
	id := params["potId"]

	cropPotDBObject, err := FindCropPotById(id, claims.Subject)
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

func FindCropPotById(id string, userId string) (*models.CropPot, error) {
	
	var cropPot models.CropPot
	result := initPackage.Db.Scopes(userScope(userId)).First(&cropPot, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cropPot, nil
}

func FindPotByToken(token string) (*models.CropPot, error) {
	var cropPot models.CropPot
	if err := initPackage.Db.Preload("ControlSettings").Where("token = ?", token).First(&cropPot).Error; err != nil {
		return nil, err
	}
	return &cropPot, nil
}

func findPotsByUserId(userId string) ([]models.CropPot, error) {
    var cropPots []models.CropPot
    result := initPackage.Db.Scopes(userScope(userId)).
        Preload("SensorDatas").
        Preload("CustomSensorFields.CustomSensorData").
        Preload("ControlSettings").
        Where("clerk_user_id = ?", userId).
        Find(&cropPots)

    if result.Error != nil {
        return nil, result.Error
    }

    return cropPots, nil
}


func userScope(userId string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("clerk_user_id = ?", userId)
    }
}