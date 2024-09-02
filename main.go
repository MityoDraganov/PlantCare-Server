package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"PlantCare/controllers"
	"PlantCare/initPackage"
	"PlantCare/middlewears"
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

	err := db.AutoMigrate(
		&models.User{},
		&models.CropPot{},
		&models.Sensor{}, 
		&models.Measurement{}, 
		&models.Control{}, 
		&models.Webhook{}, 
		&models.Update{},
		&models.ActivePeriod{},
		
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	r.HandleFunc("/users/clerk/register", controllers.ClerkUserRegister)

	// PROTECTED ROUTES

	// --CROP POTS--
	api.HandleFunc("/cropPots/assign/{token}", controllers.AssignCropPotToUser).Methods("POST")
	api.HandleFunc("/cropPots", controllers.GetCropPotsForUser).Methods("GET")

	pots := api.PathPrefix("/cropPots").Subrouter()
	pots.Use(middlewears.PotMiddleware)
	pots.HandleFunc("/{potId}", controllers.UpdateCropPot).Methods("PUT")
	pots.HandleFunc("/{potId}", controllers.RemoveCropPot).Methods("DELETE")

	// --CONTROLS--
	controls := api.PathPrefix("/controls").Subrouter()
	controls.HandleFunc("", controllers.UpdateControllSetting).Methods("PUT")

	// --WEBHOOKS--
	webhooks := api.PathPrefix("/webhooks").Subrouter()
	webhooks.Use(middlewears.PotMiddleware)

	webhooks.HandleFunc("/{potId}", controllers.AddWebhook).Methods("POST")
	webhooks.HandleFunc("/{potId}/{webhookId}", controllers.UpdateWebhook).Methods("PUT")
	webhooks.HandleFunc("/{potId}/{webhookId}", controllers.DeleteWebhook).Methods("DELETE")

	// --ADMIN ACTIONS--
	adminR := r.NewRoute().Subrouter()
	adminR.HandleFunc("/cropPots", controllers.AddCropPot).Methods("POST")

	// WEBSOCKET CONNECTIONS
	websocket.SetupWebSocketRoutes(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "192.168.0.120"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	fmt.Println("Server listening on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
