package controllers

import (
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"PlantCare/websocket/wsTypes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

func GetMessagesForUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here req")
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
		return
	}
	messages, err := FindMessagesByUserId(claims.Subject)
	fmt.Println(messages)
	if err != nil {
		fmt.Println("Error extracting session claims")
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
		return
	}
	initPackage.Db.Model(&models.Message{}).Where("clerk_user_id = ?", claims.Subject).Update("is_read", true)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "All messages marked as read"}`))
}

func CreateMessage(clerkUserId string, data string, title string, statusResponse wsTypes.StatusResponse, event wsTypes.Event) error {
	message := models.Message{
		StatusResponse: statusResponse,
		Event: event,
		Title: &title,
		ClerkUserID: clerkUserId,
		Data: data,
		IsRead: false,
	}


	result := initPackage.Db.Create(&message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}





func FindMessagesByUserId(userId string) ([]models.Message, error) {
	var messages []models.Message
	result := initPackage.Db.
		Where("clerk_user_id = ?", userId).
		Find(&messages)  // No need to preload "Inbox" here

	if result.Error != nil {
		return nil, result.Error
	}

	// Ensure messages is never nil using your utility function
	return utils.ReturnEmptyIfNil(messages), nil
}

