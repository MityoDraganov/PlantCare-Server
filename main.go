package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"PlantCare/utils/firebaseUtil"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"PlantCare/controllers"
	cronjobs "PlantCare/cronJobs"
	"PlantCare/initPackage"
	"PlantCare/middlewears"
	"PlantCare/models"
	"PlantCare/websocket"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cronjobs.StartCronJobs()

	clerk.SetKey(os.Getenv("CLERK_API_KEY"))
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	const credentialsFile = "./firebase_service_account.json"
	_, err = firebaseUtil.InitializeApp(credentialsFile)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	authMiddleware := clerkhttp.RequireHeaderAuthorization()
	api.Use(authMiddleware)

	db := InitDB()
	initPackage.SetDatabase(db)

	err = db.AutoMigrate(
		&models.User{},
		&models.CropPot{},
		&models.Sensor{},
		&models.Driver{},
		&models.Measurement{},
		&models.Condition{},
		&models.Control{},
		&models.Webhook{},
		&models.Update{},
		&models.ActivePeriod{},
		&models.Message{},
		&models.Canvas{},
		&models.PinnedCard{},
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// Set all crop pots to offline at startup
	if err := controllers.SetAllPotsOffline(); err != nil {
		log.Fatal("failed to set crop pots to offline on startup:", err)
	}

	r.HandleFunc("/users/clerk/register", controllers.ClerkUserRegister).Methods("POST");

	// PUBLIC ROUTES
	//test route
	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Server is running"}`))
		}).Methods("GET")
	// PROTECTED ROUTES

	// --INBOX--
	api.HandleFunc("/inbox", controllers.GetMessagesForUser).Methods("GET")
	api.HandleFunc("/inbox", controllers.MarkAllAsRead).Methods("PUT")

	// --CROP POTS--
	api.HandleFunc("/cropPots/assign/{token}", controllers.AssignCropPotToUser).Methods("POST")
	api.HandleFunc("/cropPots", controllers.GetCropPotsForUser).Methods("GET")

	pots := api.PathPrefix("/cropPots").Subrouter()
	pots.Use(middlewears.PotMiddleware)
	pots.HandleFunc("/{potId}", controllers.UpdateCropPot).Methods("PUT")
	pots.HandleFunc("/{potId}", controllers.RemoveCropPot).Methods("DELETE")

	// --CONTROLS--
	controls := api.PathPrefix("/controls").Subrouter()
	controls.HandleFunc("", controllers.UpdateControls).Methods("PUT")

	// --SENSORS--
	sensors := api.PathPrefix("/sensors").Subrouter()
	sensors.HandleFunc("/{sensorId}", controllers.UpdateSensor).Methods("PUT")
	sensors.HandleFunc("", controllers.UpdateSensor).Methods("PUT")

	// --WEBHOOKS--
	webhooks := api.PathPrefix("/webhooks").Subrouter()
	webhooks.Use(middlewears.PotMiddleware)

	webhooks.HandleFunc("/{potId}", controllers.AddWebhook).Methods("POST")
	webhooks.HandleFunc("/{potId}/{webhookId}", controllers.UpdateWebhook).Methods("PUT")
	webhooks.HandleFunc("/{potId}/{webhookId}", controllers.DeleteWebhook).Methods("DELETE")

	// --DRIVERS--
	drivers := api.PathPrefix("/drivers").Subrouter()
	drivers.HandleFunc("", controllers.GetAllDrivers).Methods("GET")
	drivers.HandleFunc("", controllers.UploadDriver).Methods("POST")
	drivers.HandleFunc("/{driverId}", controllers.EditDriver).Methods("PUT")
	drivers.HandleFunc("/{driverId}", controllers.DeleteDriver).Methods("DELETE")

	//	--CANVAS
	canvas := api.PathPrefix("/canvas").Subrouter()
	canvas.Use(middlewears.PotMiddleware)
	canvas.HandleFunc("", controllers.GetCanvasesByUser).Methods("GET")
	canvas.HandleFunc("/{potId}", controllers.UpdateCanvas).Methods("PUT")
	//canvas.HandleFunc("/{potId}", controllers.DeleteCanvas).Methods("DELETE")

	// --ADMIN ACTIONS--
	adminR := r.NewRoute().Subrouter()
	adminR.HandleFunc("/cropPots", controllers.AddCropPot).Methods("POST")

	// WEBSOCKET CONNECTIONS
	websocket.SetupWebSocketRoutes(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "192.168.0.120", "https://plantscare.sytes.net", "http://plantscare.sytes.net"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	fmt.Println("Server listening on port 8000!")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
