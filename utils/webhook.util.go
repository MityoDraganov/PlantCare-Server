package utils

import (
	"PlantCare/websocket/wsDtos"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func TriggerWebhook(endpointUrl string, payload wsDtos.WebhookPayload) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
	}

	resp, err := http.Post(endpointUrl, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Failed to send webhook to %s: %v", endpointUrl, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Webhook to %s returned non-OK status: %v", endpointUrl, resp.Status)
	}
}