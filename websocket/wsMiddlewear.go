package websocket

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"PlantCare/controllers"
	"PlantCare/utils"
	"PlantCare/websocket/connectionManager"
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
func Middleware(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the query parameters.
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
		return
	}

	fmt.Println(token)

	// Authenticate using the token
	cropPotDbObject, err := controllers.FindPotByToken(token)
	if err != nil {
		utils.JsonError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	potIDStr := strconv.FormatUint(uint64(cropPotDbObject.ID), 10)

	ctx := context.WithValue(r.Context(), wsTypes.CropPotIDKey, potIDStr)
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
        Role: wsTypes.PotRole,
	}

	connectionManager.ConnManager.AddConnection(potIDStr, connection)
	defer connectionManager.ConnManager.RemoveConnection(string(cropPotDbObject.ID))

	wsutils.SendValidRequest(&connection, cropPotDbObject)
	// Start handling messages
	go HandleMessages(&connection)
	go wsutils.SendMessages(&connection)

	fmt.Printf("New WebSocket connection established from IP: %s\n", remoteIP)

}

func SetupWebSocketRoutes(r *mux.Router) {
	ws := r.PathPrefix("/v1").Subrouter()

    potWs := ws.PathPrefix("/pots").Subrouter()

    userWs := ws.PathPrefix("/users").Subrouter()
	potWs.HandleFunc("/", Middleware)


    userWs.HandleFunc("/", userWsMiddlewear)

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
            Role: wsTypes.UserRole,
        }
    
        connectionManager.ConnManager.AddConnection(userID, connection)
        //defer connectionManager.ConnManager.RemoveConnection(string(userID))
    
        wsutils.SendValidRequest(&connection, userDbObject)
        // Start handling messages

		messages, err := controllers.FindMessagesByUserId(claims.Subject)
		if err != nil {
			wsutils.SendErrorResponse(&connection, http.StatusInternalServerError)

		}
		for _, message := range messages {
			messageDto := wsDtos.NotificationDto{
				Title: utils.CoalesceString(message.Title),
				Text: message.Text,
				IsRead: message.IsRead,
				Timestamp: message.CreatedAt,

			}
			
			wsutils.SendMessage(&connection, wsTypes.MessageFound, "", messageDto)
		}
        go HandleMessages(&connection)
        go wsutils.SendMessages(&connection)
    
        fmt.Printf("New WebSocket connection established from IP: %s\n", remoteIP)
    
}

// NewInMemoryJWKStore creates a new in-memory JWK store.
func NewInMemoryJWKStore() *InMemoryJWKStore {
	return &InMemoryJWKStore{}
}

