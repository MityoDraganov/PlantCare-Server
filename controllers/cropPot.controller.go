package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/utils"

	"PlantCare/models"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gorilla/mux"
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

    var cropPotResponses []dtos.CropPotResponse
    for _, cropPot := range cropPots {
        var controlsResponse []dtos.ControlDto

        for _, control := range cropPot.Controls {
            // Convert duration to string in format "15:04"
            startStr := durationToTimeString(control.ActivePeriod.Start)
            endStr := durationToTimeString(control.ActivePeriod.End)

            activePeriod := dtos.ActivePeriod{
                ID:    control.ActivePeriod.ID,
                Start: startStr,
                End:   endStr,
				Days: utils.ParseBitmask(control.ActivePeriod.Days),
            }
            controlsResponse = append(controlsResponse, dtos.ControlDto{
                ID:           control.ID,
                SerialNumber: control.SerialNumber,
                Alias:        control.Alias,
                Description:  utils.CoalesceString(control.Description),
                Updates:      control.Updates,
                IsOfficial:   true,

                OnCondition:  control.OnCondition,
                OffCondition: control.OffCondition,
                ActivePeriod: activePeriod,
            })
        }

        // Map SensorData
        var sensorDataResponses []dtos.SensorDto
        for _, sensorData := range cropPot.Sensors {
            sensorDataResponses = append(sensorDataResponses, dtos.SensorDto{
                ID:           sensorData.ID,
                SerialNumber: sensorData.SerialNumber,
                Measurements: sensorData.Measurements,
                Description:  utils.CoalesceString(sensorData.Description),
                Alias:        sensorData.Alias,
				MeasuremntInterval: durationToTimeString(sensorData.MeasuremntInterval),
            })
        }

        webhookResponses := []dtos.WebhookResponse{}

        for _, webhook := range cropPot.Webhooks {
            // Initialize subscribedEvents as an empty slice
            subscribedEvents := []dtos.SensorDto{}

            // Populate subscribedEvents if there are any
            for _, event := range webhook.SubscribedEvents {
                subscribedEvent := dtos.SensorDto{
                    SerialNumber: event.SerialNumber,
                    Alias:        event.Alias,
                    Description:  utils.CoalesceString(event.Description),
                }

                subscribedEvents = append(subscribedEvents, subscribedEvent)
            }

            webhookResponse := dtos.WebhookResponse{
                ID:               webhook.ID,
                EndpointUrl:      webhook.EndpointUrl,
                Description:      utils.CoalesceString(webhook.Description),
                SubscribedEvents: subscribedEvents, // Will be an empty slice if no events
            }

            webhookResponses = append(webhookResponses, webhookResponse)
        }

        cropPotResponse := dtos.CropPotResponse{
            ID:         cropPot.ID,
            Alias:      cropPot.Alias,
            IsArchived: cropPot.IsArchived,
            Controls:   controlsResponse,
            Sensors:    sensorDataResponses,
            Webhooks:   webhookResponses,
        }
        cropPotResponses = append(cropPotResponses, cropPotResponse)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cropPotResponses)
}

// Helper function to convert duration to time string in format "15:04"
func durationToTimeString(d time.Duration) string {
    hours := int(d / time.Hour)
    minutes := int((d % time.Hour) / time.Minute)
    return fmt.Sprintf("%02d:%02d", hours, minutes)
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
	var cropPotDto dtos.CreateCropPot
	err := json.NewDecoder(r.Body).Decode(&cropPotDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	id := params["id"]

	cropPotDBObject, err := FindCropPotById(id)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	initPackage.Db.Model(&cropPotDBObject).Clauses(clause.Returning{}).Updates(cropPotDto)

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
	result := initPackage.Db.First(&cropPot, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cropPot, nil
}

func FindPotByToken(token string) (*models.CropPot, error) {
	var cropPot models.CropPot
	if err := initPackage.Db.Preload("Controls").Where("token = ?", token).First(&cropPot).Error; err != nil {
		return nil, err
	}
	return &cropPot, nil
}

func findPotsByUserId(userId string) ([]models.CropPot, error) {
	var cropPots []models.CropPot
	result := initPackage.Db.
		Preload("Sensors").
		Preload("Sensors.Measurements").
		Preload("Controls").
		Preload("Controls.ActivePeriod").
		Preload("Controls.Updates").
		Preload("Webhooks").
		Preload("Webhooks.SubscribedEvents").
		Where("clerk_user_id = ?", userId).
		Find(&cropPots)

	if result.Error != nil {
		return nil, result.Error
	}

	return cropPots, nil
}

