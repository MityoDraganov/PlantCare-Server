package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"PlantCare/controllers"
	"PlantCare/utils"
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsTypes"
	"PlantCare/websocket/wsUtils"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connManager = connectionManager.GetInstance() // Initialize the connection manager
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

	connection := &wsTypes.Connection{
		Conn:    conn,
		Send:    make(chan []byte),
		Context: r.Context(),
		IP:      remoteIP,
	}

	connManager.AddConnection(potIDStr, connection)
	defer connManager.RemoveConnection(string(cropPotDbObject.ID))

	wsutils.SendValidRequest(connection, cropPotDbObject)
	// Start handling messages
	go HandleMessages(connection)
	go wsutils.SendMessages(connection)

	fmt.Printf("New WebSocket connection established from IP: %s\n", remoteIP)

}

func SetupWebSocketRoutes(r *mux.Router) {
	ws := r.PathPrefix("/v1").Subrouter()
	ws.HandleFunc("/", Middleware)
}


func toJSON(data interface{}) json.RawMessage {
	if data == nil {
		return nil
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return nil
	}
	return dataBytes
}