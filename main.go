package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"TravelBuddy/controllers"
	"TravelBuddy/middlewares"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"TravelBuddy/models"
)

func InitDB() *gorm.DB {
	dsn := "sqlserver://server:server123@localhost:1433?database=Travel_Buddy"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}

func main() {
	r := mux.NewRouter()

	db := InitDB()

	err := db.AutoMigrate(&models.User{}, &models.Passenger{}, &models.Trip{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to SQL Server!")
	// Pre-populate with some users (optional)

	//tripCollection := client.Database("travelbuddy").Collection("trips")
	//controllers.SetTripCollection(tripCollection)

	r.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{username}", controllers.GetUser).Methods("GET")
	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")

	updateUserHandler := middlewares.LoggerMiddleware(http.HandlerFunc(controllers.UpdateUser))
	r.Handle("/users/{username}", updateUserHandler).Methods("PUT")

	deleteUserHandler := middlewares.LoggerMiddleware(http.HandlerFunc(controllers.DeleteUser))
	r.Handle("/users/{username}", deleteUserHandler).Methods("DELETE")

	r.HandleFunc("/trips", controllers.CreateTrip).Methods("POST")
	r.HandleFunc("/trips", controllers.DeleteTrip).Methods("DELETE")
	r.HandleFunc("/trips", controllers.UpdateTrip).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", r))
}
