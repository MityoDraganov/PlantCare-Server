package websocket

import (
	"context"
	"fmt"
	"net/http"

	"PlantCare/controllers"
	"PlantCare/utils"
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsTypes"
	"PlantCare/websocket/wsUtils"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// MiddlewareFunc defines a function type for middleware.
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc
var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connManager = connectionManager.NewConnectionManager() // Initialize the connection manager
)


// ApplyMiddleware applies a series of middleware functions to an HTTP handler function.
func ApplyMiddleware(h http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}


func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}

	connection := &wsTypes.Connection{
		Conn: conn,
		Send: make(chan []byte),
		Context: r.Context(),
	}

	potID, ok := r.Context().Value(wsTypes.CropPotIDKey).(string)
	if ok {
		connManager.AddConnection(potID, connection)
		defer connManager.RemoveConnection(potID)
	}

	go HandleMessages(connection)
	go wsutils.SendMessages(connection)
}


// LoggingMiddleware logs basic information about the HTTP request.
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New WebSocket connection")
		next.ServeHTTP(w, r)
	}
}


// TokenLoggingMiddleware extracts and prints the ?token query parameter from the request.
func AuthMiddlewear(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the query parameters.
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		cropPotDbObject, err := controllers.FindPotByToken(token)
		if err != nil {
			utils.JsonError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), wsTypes.CropPotIDKey, cropPotDbObject.ID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
// SetupWebSocketRoutes sets up the WebSocket routes with the provided middleware.
func SetupWebSocketRoutes(r *mux.Router) {
	ws := r.PathPrefix("/v1").Subrouter()
	ws.HandleFunc("/", ApplyMiddleware(HandleConnection, AuthMiddlewear, LoggingMiddleware))
}

