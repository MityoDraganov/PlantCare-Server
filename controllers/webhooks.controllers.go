package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func AddWebhook(w http.ResponseWriter, r *http.Request) {
	var webhookDto dtos.AddWebhookDto
	err := json.NewDecoder(r.Body).Decode(&webhookDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	potIdStr := params["potId"]

	potId, err := strconv.ParseUint(potIdStr, 10, 32)
	if err != nil {
		utils.JsonError(w, "Invalid potId", http.StatusBadRequest)
		return
	}

	var subscribedEvents []models.Sensor
	for _, subscribedEvent := range webhookDto.SubscribedEvents {
		subscibedEvent, err := FindSensorBySerialNum(subscribedEvent.SerialNumber)
		if err != nil {
			utils.JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		subscribedEvents = append(subscribedEvents, *subscibedEvent)
	}

	webhook := models.Webhook{
		CropPotID:        uint(potId),
		EndpointUrl:      webhookDto.EndpointUrl,
		SubscribedEvents: subscribedEvents,
		Description:      webhookDto.Description,
	}

	result := initPackage.Db.Create(&webhook)
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	var subscribedEventsDto []dtos.SensorDto

	for _, event := range webhook.SubscribedEvents {
		subscribedEvent := dtos.SensorDto{
			SerialNumber: event.SerialNumber,
			Alias:        utils.CoalesceString(event.Alias),
			Description:  event.Description,
		}

		subscribedEventsDto = append(subscribedEventsDto, subscribedEvent)
	}

	webhookResponse := dtos.WebhookResponse{
		ID:               webhook.ID,
		EndpointUrl:      webhook.EndpointUrl,
		SubscribedEvents: subscribedEventsDto,
		Description:      webhook.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhookResponse)
}

func UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	var webhookDto dtos.UpdateWebhookDto

	// Decode the JSON body into webhookDto
	err := json.NewDecoder(r.Body).Decode(&webhookDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the webhook ID from the URL parameters
	params := mux.Vars(r)
	id := params["webhookId"]

	// Find the existing webhook by ID
	var existingWebhook models.Webhook
	if err := initPackage.Db.Preload("SubscribedEvents").First(&existingWebhook, id).Error; err != nil {
		utils.JsonError(w, "Webhook not found", http.StatusNotFound)
		return
	}

	// Handle subscribed events
	if webhookDto.SubscribedEvents != nil {
		// Clear existing associations
		if err := initPackage.Db.Model(&existingWebhook).Association("SubscribedEvents").Clear(); err != nil {
			utils.JsonError(w, "Failed to clear existing subscribed events", http.StatusInternalServerError)
			return
		}

		var newSubscribedEvents []models.Sensor
		for _, subscribedEvent := range *webhookDto.SubscribedEvents {
			subscribedEventDbo, err := FindSensorBySerialNum(subscribedEvent.SerialNumber)
			if err != nil {
				utils.JsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
			newSubscribedEvents = append(newSubscribedEvents, *subscribedEventDbo)
		}
		// Update with the new associations
		existingWebhook.SubscribedEvents = newSubscribedEvents
	}

	// Update the fields of the existing webhook with provided values
	if webhookDto.EndpointUrl != nil {
		existingWebhook.EndpointUrl = *webhookDto.EndpointUrl
	}
	if webhookDto.Description != nil {
		existingWebhook.Description = webhookDto.Description
	}

	// Save the updated webhook to the database
	if err := initPackage.Db.Save(&existingWebhook).Error; err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	webhookResponse := dtos.WebhookResponse{
		ID:               existingWebhook.ID,
		EndpointUrl:      existingWebhook.EndpointUrl,
		Description:      existingWebhook.Description,
		SubscribedEvents: *webhookDto.SubscribedEvents,
	}

	// Set the response header and encode the updated webhook object
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(webhookResponse); err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	webhookId := params["webhookId"]

	webhook, err := findWebhookById(webhookId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Webhook not found, return a 404 status
			http.Error(w, "Webhook not found", http.StatusNotFound)
		} else {
			// Some other error occurred, return a 500 status
			http.Error(w, "Error finding webhook", http.StatusInternalServerError)
		}
		return
	}

	// Delete the webhook
	if err := initPackage.Db.Delete(webhook).Error; err != nil {
		http.Error(w, "Error deleting webhook", http.StatusInternalServerError)
		return
	}

	// Set the response header and encode the updated webhook object
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func findWebhookById(id string) (*models.Webhook, error) {
	var webhook models.Webhook
	result := initPackage.Db.First(&webhook, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &webhook, nil
}

func GetSubscribedWebhooksForSensor(sensorID uint) ([]models.Webhook, error) {
	var webhooks []models.Webhook
	err := initPackage.Db.Joins("JOIN webhook_sensors ON webhook_sensors.webhook_id = webhooks.id").
		Where("webhook_sensors.sensor_id = ?", sensorID).
		Find(&webhooks).Error
	return webhooks, err
}

// Maps a single Webhook to WebhookResponse DTO
func mapWebhookToDTO(webhook models.Webhook) dtos.WebhookResponse {
	subscribedEvents := []dtos.SensorDto{}

	// Populate subscribedEvents if there are any
	for _, event := range webhook.SubscribedEvents {
		subscribedEvent := dtos.SensorDto{
			SerialNumber: event.SerialNumber,
			Alias:        utils.CoalesceString(event.Alias),
			Description:  utils.CoalesceString(event.Description),
		}
		subscribedEvents = append(subscribedEvents, subscribedEvent)
	}

	return dtos.WebhookResponse{
		ID:               webhook.ID,
		EndpointUrl:      webhook.EndpointUrl,
		Description:      utils.CoalesceString(webhook.Description),
		SubscribedEvents: subscribedEvents, // Empty slice if no events
	}
}

func ToWebhooksDTO(input interface{}) []dtos.WebhookResponse {
	switch v := input.(type) {
	case models.Webhook:
		// If it's a single webhook, wrap it in a slice
		return []dtos.WebhookResponse{mapWebhookToDTO(v)}
	case []models.Webhook:
		// If it's a slice of webhooks, map each webhook to WebhookResponse
		webhookDTOs := make([]dtos.WebhookResponse, len(v))
		for i, webhook := range v {
			webhookDTOs[i] = mapWebhookToDTO(webhook)
		}
		return webhookDTOs
	default:
		// Handle unexpected types by returning an empty slice
		return []dtos.WebhookResponse{}
	}
}
