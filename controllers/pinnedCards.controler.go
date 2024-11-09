// Import additional packages at the top if needed
package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetCanvasesByUser retrieves all canvases and their pinned cards for crop pots belonging to the current user.
func GetCanvasesByUser(w http.ResponseWriter, r *http.Request) {
	// Assume user ID is available in context; replace with your actual user ID retrieval method
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		utils.JsonError(w, "Unauthorized: unable to extract session claims", http.StatusUnauthorized)
		return
	}
	clerkUserID := claims.Subject

	var canvases []models.Canvas
	if err := initPackage.Db.
		Joins("JOIN crop_pots ON crop_pots.id = canvases.crop_pot_id").
		Where("crop_pots.user_id = ?", clerkUserID).
		Preload("PinnedCards").
		Find(&canvases).Error; err != nil {
		utils.JsonError(w, "Failed to retrieve canvases", http.StatusInternalServerError)
		return
	}

	// Convert canvases to DTOs
	var canvasDtos []dtos.CanvasDto
	for _, canvas := range canvases {
		canvasDtos = append(canvasDtos, ToCanvasDTO(canvas))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(canvasDtos)
}

// CreateCanvas creates a new Canvas with associated PinnedCards.
func CreateCanvas(w http.ResponseWriter, r *http.Request) {
	var canvasDto dtos.CanvasDto
	err := json.NewDecoder(r.Body).Decode(&canvasDto)
	if err != nil {
		log.Println(err)
		utils.JsonError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Convert DTO to Model with linked PinnedCards

	canvas := models.Canvas{
		CropPotID: canvasDto.CropPotID,
	}

	result := initPackage.Db.Create(&canvas).Clauses(clause.Returning{})
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	pinnedCards, _ := toPinnedCardsModel(canvasDto.PinnedCards, canvas.ID)

	canvas.PinnedCards = pinnedCards

	result = initPackage.Db.Updates(&canvas).Clauses(clause.Returning{})
	if result.Error != nil {
		utils.JsonError(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(canvas)
}

// UpdateCanvas updates an existing Canvas and its PinnedCards.
func UpdateCanvas(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("canvasId") // Canvas ID passed as a query parameter
	var canvasDto dtos.CanvasDto

	err := json.NewDecoder(r.Body).Decode(&canvasDto)
	if err != nil {
		log.Println(err)
		utils.JsonError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Find existing canvas
	var canvas models.Canvas
	if err := initPackage.Db.Preload("PinnedCards").First(&canvas, id).Error; err != nil {
		utils.JsonError(w, "Canvas not found", http.StatusNotFound)
		return
	}

	// Delete old PinnedCards and create new ones
	initPackage.Db.Model(&canvas).Association("PinnedCards").Clear()
	newPinnedCards, _ := toPinnedCardsModel(canvasDto.PinnedCards, canvas.ID)
	canvas.CropPotID = canvasDto.CropPotID
	canvas.PinnedCards = newPinnedCards

	if err := initPackage.Db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&canvas).Error; err != nil {
		utils.JsonError(w, "Failed to update canvas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(canvas)
}

// DeleteCanvas deletes a Canvas and its associated PinnedCards.
func DeleteCanvas(id uint) *error {
	if err := initPackage.Db.Delete(&models.Canvas{}, id).Error; err != nil {
		return &err
	}

	return nil
}

// Helper function to convert DTO PinnedCards to Model PinnedCards
func toPinnedCardsModel(dtos []dtos.PinnedCardDto, canvasId uint) ([]models.PinnedCard, error) {
	var cards []models.PinnedCard
	for _, dto := range dtos {
		card := models.PinnedCard{
			CanvasID:      canvasId,
			Title:         utils.CoalesceString(&dto.Title),
			Icon:          utils.CoalesceString(&dto.Icon),
			SensorID:      dto.SensorID,
			StartLocation: dto.StartLocation,
			Width:         dto.Width,
			Height:        dto.Height,
			Type: dto.Type,
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func ToPinnedCardsDTO(cards []models.PinnedCard) []dtos.PinnedCardDto {
	var cardsDto []dtos.PinnedCardDto

	// Check if there are any cards
	if len(cards) == 0 {
		fmt.Println("No cards to convert")
		return []dtos.PinnedCardDto{} // Return an empty slice
	}

	for _, card := range cards {
		cardsDto = append(cardsDto, dtos.PinnedCardDto{
			Title:    *utils.CoalesceString(card.Title),
			Icon:     *utils.CoalesceString(card.Icon),
			SensorID: card.SensorID,
			StartLocation: card.StartLocation,
			Width:    card.Width,
			Height:   card.Height,
			Type: card.Type,
		})
	}
	return cardsDto
}

func ToCanvasDTO(canvas models.Canvas) dtos.CanvasDto {
	return dtos.CanvasDto{
		CropPotID:   canvas.CropPotID,
		PinnedCards: ToPinnedCardsDTO(canvas.PinnedCards),
	}
}
