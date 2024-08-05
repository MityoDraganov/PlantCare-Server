package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"PlantCare/controllers"
	"PlantCare/models"
	"PlantCare/websocket"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

func InitDB() *gorm.DB {
	dsn := "sqlserver://server:P@ssw0rd123!@localhost:1433?database=Plant_Care"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}


func main() {
	// Initialize the router
	r := mux.NewRouter()

	// Apply the path prefix
	api := r.PathPrefix("/api/v1").Subrouter()

	clerk.SetKey("sk_test_gy7eUfUIotA7K6RXGOa0VJBUclqUhHRSmOI6zqriDU")

	// Initialize the database
	db := InitDB()
	controllers.SetDatabase(db)

	// Auto migrate the database models
	err := db.AutoMigrate(&models.User{}, &models.CropPot{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to SQL Server!")

	authMiddleware := clerkhttp.WithHeaderAuthorization()
	api.Use(authMiddleware)

	// PUBLIC ROUTES
	r.HandleFunc("/users/clerk/register", controllers.ClerkUserRegister)

	// PROTECTED ROUTES
	api.HandleFunc("/cropPots", controllers.GetCropPotsForUser).Methods("GET")
	api.HandleFunc("/cropPots/assign/{token}", controllers.AssignCropPotToUser).Methods("POST")
	api.HandleFunc("/cropPots/{id}", controllers.UpdateCropPot).Methods("PUT")
	api.HandleFunc("/cropPots/{id}", controllers.RemoveCropPot).Methods("DELETE")

	// WEBSOCKET CONNECTIONS
	websocket.SetupWebSocketRoutes(r)

	// ADMIN ACTIONS
	r.HandleFunc("/cropPots", controllers.AddCropPot).Methods("POST")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", handler))
	fmt.Println("Server listening on port 8080!")
}
