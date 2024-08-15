package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"PlantCare/controllers"
	"PlantCare/initPackage"
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
	clerk.SetKey("sk_test_gy7eUfUIotA7K6RXGOa0VJBUclqUhHRSmOI6zqriDU")
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	authMiddleware := clerkhttp.RequireHeaderAuthorization()
	api.Use(authMiddleware)

	db := InitDB()
	initPackage.SetDatabase(db)

	err := db.AutoMigrate(&models.User{}, &models.CropPot{}, &models.SensorData{}, &models.ControlSettings{}, &models.CustomSensorField{}, &models.CustomSensorData{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}


	r.HandleFunc("/users/clerk/register", controllers.ClerkUserRegister)

	// PROTECTED ROUTES

	pots := api.PathPrefix("/cropPots").Subrouter()
	pots.HandleFunc("/assign/{token}", controllers.AssignCropPotToUser).Methods("POST")


	pots.HandleFunc("", controllers.GetCropPotsForUser).Methods("GET")
	pots.HandleFunc("/{potId}", controllers.UpdateCropPot).Methods("PUT")
	pots.HandleFunc("/{potId}", controllers.RemoveCropPot).Methods("DELETE")

	pots.HandleFunc("/controlls/{controllSettingsId}", controllers.UpdateControllSettings).Methods("PUT")

	// WEBSOCKET CONNECTIONS
	websocket.SetupWebSocketRoutes(r)

	// ADMIN ACTIONS
	adminR := r.NewRoute().Subrouter()
	adminR.HandleFunc("/cropPots", controllers.AddCropPot).Methods("POST")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Start the server
	fmt.Println("Server listening on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Worked!")
}