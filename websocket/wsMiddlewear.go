package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// MiddlewareFunc represents a WebSocket middleware function
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// ApplyMiddleware applies middleware to the WebSocket handler
func ApplyMiddleware(h http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

// Example middleware that logs WebSocket connections
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New WebSocket connection")
		next.ServeHTTP(w, r)
	}
}

// Use ApplyMiddleware when setting up the route
func SetupWebSocketRoutes(r *mux.Router) {
	ws := r.PathPrefix("/ws/v1").Subrouter()
	ws.HandleFunc("/connect", ApplyMiddleware(HandleConnection, LoggingMiddleware))
}
