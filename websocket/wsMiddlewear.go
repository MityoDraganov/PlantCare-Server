package websocket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"PlantCare/controllers"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/otaManager"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	"PlantCare/websocket/wsUtils"

	"github.com/clerk/clerk-sdk-go/v2"

	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Middleware handles authentication, connection, and logging in a single function.
func PotMiddleware(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the query parameters.
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
		return
	}

	// Authenticate using the token
	cropPotDbObject, err := controllers.FindPotByToken(token)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cropPotDbObject.Status = models.StatusOnline
	if err := initPackage.Db.Save(&cropPotDbObject).Error; err != nil {
		fmt.Println("Error while updating pot connection status:", err)
		return
	}

	potIDStr := strconv.FormatUint(uint64(cropPotDbObject.ID), 10)

	ctx := context.WithValue(r.Context(), wsTypes.CropPotIDKey, potIDStr)
	r = r.WithContext(ctx)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}

	remoteIP := r.RemoteAddr

	connection := wsTypes.Connection{
		Conn:    conn,
		Send:    make(chan []byte),
		Context: r.Context(),
		IP:      remoteIP,
		Role:    wsTypes.PotRole,
	}

	// Add connection to the manager
	connectionManager.ConnManager.AddConnection(potIDStr, connection)
	isOtaPending := otaManager.OTAManager.IsOTAPending(potIDStr)
	if isOtaPending {
		go func() {
			pendingOta, ok := otaManager.OTAManager.GetPendingOTA(potIDStr)
			if !ok {
				err := errors.New("error with the pending ota")
				fmt.Println(err)

				utils.JsonError(w, err.Error(), http.StatusBadRequest)
				return
			}

			connection, ok := connectionManager.ConnManager.GetConnection(potIDStr)
			if !ok {
				err := errors.New("connection not found for pot ID: " + potIDStr)
				fmt.Println(err)

				utils.JsonError(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := utils.UploadMultipleDrivers(pendingOta.DriverURLs, connection); err != nil {
				fmt.Printf("Failed to upload driver: %v", err)
			}
		}()
	}

	wsutils.SendValidRequest(&connection, controllers.ToCropPotResponseDTO(*cropPotDbObject))

	ownerConnection, exists := connectionManager.ConnManager.GetConnectionByOwner(*cropPotDbObject.ClerkUserID)
	if exists {

		wsutils.SendMessage(ownerConnection, "", wsTypes.UpdatedPot, controllers.ToCropPotResponseDTO(*cropPotDbObject))
	}

	// Start handling messages
	go HandleMessages(&connection, cropPotDbObject)
	go wsutils.SendMessages(&connection)

	fmt.Printf("New WebSocket connection established from IP: %s\n", remoteIP)

}

func SetupWebSocketRoutes(r *mux.Router) {
	ws := r.PathPrefix("/v1").Subrouter()

	potWs := ws.PathPrefix("/pots").Subrouter()

	userWs := ws.PathPrefix("/users").Subrouter()
	potWs.HandleFunc("/", PotMiddleware)

	userWs.HandleFunc("/", userWsMiddlewear)

}

// UserWsMiddlewear verifies Clerk session tokens for WebSocket connections and extracts user ID.
func userWsMiddlewear(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
		return
	}

	// Initialize JWK Store and JWKS Client
	jwkStore := NewInMemoryJWKStore()
	config := &clerk.ClientConfig{}
	config.Key = clerk.String(os.Getenv(("CLERK_API_KEY")))
	jwksClient := jwks.NewClient(config)

	// Attempt to get the JSON Web Key from the store.
	jwk := jwkStore.GetJWK()
	if jwk == nil {
		// Decode the session JWT to find the key ID.
		claims, err := jwt.Decode(r.Context(), &jwt.DecodeParams{
			Token: token,
		})
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Fetch the JSON Web Key from Clerk.
		jwk, err = jwt.GetJSONWebKey(r.Context(), &jwt.GetJSONWebKeyParams{
			KeyID:      claims.KeyID,
			JWKSClient: jwksClient,
		})
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Store the JWK for future use.
		jwkStore.SetJWK(jwk)
	}

	// Verify the session token.
	_, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
		Token: token,
		JWK:   jwk,
	})
	if err != nil {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract the user ID from the JWT token.
	claims, err := jwt.Decode(r.Context(), &jwt.DecodeParams{
		Token: token,
	})
	if err != nil {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	userID := claims.Subject

	userDbObject, err := controllers.FindUserById(userID)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}
	if userDbObject == nil {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), wsTypes.CropPotIDKey, userID)
	r = r.WithContext(ctx)

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}

	remoteIP := r.RemoteAddr

	connection := wsTypes.Connection{
		Conn:    conn,
		Send:    make(chan []byte),
		Context: r.Context(),
		IP:      remoteIP,
		Role:    wsTypes.UserRole,
	}

	connectionManager.ConnManager.AddConnection(userID, connection)

	messages, err := controllers.FindMessagesByUserId(claims.Subject)
	if err != nil {
		wsutils.SendErrorResponse(&connection, http.StatusInternalServerError)

	}
	for _, message := range messages {
		messageDto := wsDtos.NotificationDto{
			Title:     utils.CoalesceString(message.Title),
			Text:      message.Text,
			IsRead:    message.IsRead,
			Timestamp: message.CreatedAt,
		}

		wsutils.SendMessage(&connection, wsTypes.MessageFound, "", messageDto)
	}
	go HandleMessages(&connection, nil)
	go wsutils.SendMessages(&connection)

	fmt.Printf("New WebSocket connection established from IP: %s\n", remoteIP)

}

// NewInMemoryJWKStore creates a new in-memory JWK store.
func NewInMemoryJWKStore() *InMemoryJWKStore {
	return &InMemoryJWKStore{}
}

// JWKStore defines an interface for storing and retrieving JSON Web Keys.
type JWKStore interface {
	GetJWK() *clerk.JSONWebKey
	SetJWK(*clerk.JSONWebKey)
}

// InMemoryJWKStore is a simple in-memory implementation of JWKStore.
type InMemoryJWKStore struct {
	jwk *clerk.JSONWebKey
}

// GetJWK retrieves the JSON Web Key.
func (s *InMemoryJWKStore) GetJWK() *clerk.JSONWebKey {
	return s.jwk
}

// SetJWK stores the JSON Web Key.
func (s *InMemoryJWKStore) SetJWK(jwk *clerk.JSONWebKey) {
	s.jwk = jwk
}
