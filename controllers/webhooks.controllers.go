package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

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

	webhook := models.Webhook{
		CropPotID:   uint(potId),
		EndpointUrl: webhookDto.EndpointUrl,
		SubscribedEvents: webhookDto.SubscribedEvents,
	}

	webhookDbObject := initPackage.Db.Create(&webhook)
	if webhookDbObject.Error != nil {
		utils.JsonError(w, webhookDbObject.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhook)
}


func RemoveWebhook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	webhookId := params["webhookId"]

	initPackage.Db.Delete(models.Webhook{}, webhookId)
}

func findWebhookById(id string) (*models.Webhook, error) {
	var webhook models.Webhook
	result := initPackage.Db.First(&webhook, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &webhook, nil
}
