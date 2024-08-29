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
	"gorm.io/gorm/clause"
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
	for _, subscribedEventSerialNum := range webhookDto.SubscribedEventsSerialNums {
		subscibedEvent, err := FindSensorBySerialNum(subscribedEventSerialNum)
		if err != nil {
			utils.JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		subscribedEvents = append(subscribedEvents, *subscibedEvent)
	}

	webhook := models.Webhook{
		CropPotID:   uint(potId),
		EndpointUrl: webhookDto.EndpointUrl,
		SubscribedEvents: subscribedEvents,
		Description: webhookDto.Description,
	}

	webhookDbObject := initPackage.Db.Create(&webhook)
	if webhookDbObject.Error != nil {
		utils.JsonError(w, webhookDbObject.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhook)
}

func UpdateWebhook (w http.ResponseWriter, r *http.Request) {
	var WebhookDto dtos.AddWebhookDto

	err := json.NewDecoder(r.Body).Decode(&WebhookDto)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	id := params["webhookId"]

	webhookDBObject, err := findWebhookById(id)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	initPackage.Db.Model(&webhookDBObject).Clauses(clause.Returning{}).Updates(WebhookDto)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhookDBObject)
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

	// Successfully deleted, return a 204 status (No Content)
	w.WriteHeader(http.StatusNoContent)
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
