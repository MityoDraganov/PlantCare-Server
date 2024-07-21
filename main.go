package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"PlantCare/controllers"
	"PlantCare/middlewares"
	"PlantCare/models"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "sqlserver://server:server123@localhost:1433?database=Plant_Care"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}

func main() {
	r := mux.NewRouter()

	db := InitDB()
	controllers.SetDatabase(db)

	err := db.AutoMigrate(&models.User{}, &models.CropPot{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to SQL Server!")

	// Public routes
	publicRoutes := r.PathPrefix("/auth").Subrouter()
	publicRoutes.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	publicRoutes.HandleFunc("/login", controllers.LoginUser).Methods("POST")

	// Apply middleware to all other routes
	protectedRoutes := r.PathPrefix("/").Subrouter()
	protectedRoutes.Use(middlewares.TokenMiddleware)

	// User routes
	protectedRoutes.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	protectedRoutes.HandleFunc("/users/{username}", controllers.GetUser).Methods("GET")
	protectedRoutes.HandleFunc("/users/{username}", controllers.UpdateUser).Methods("PUT")
	protectedRoutes.HandleFunc("/users/{username}", controllers.DeleteUser).Methods("DELETE")

	// Crop Pot routes
	protectedRoutes.HandleFunc("/cropPots", controllers.AddCropPot).Methods("POST")
	protectedRoutes.HandleFunc("/cropPots/{id}", controllers.UpdateCropPot).Methods("PUT")


	log.Fatal(http.ListenAndServe(":8080", r))
	fmt.Println("Server listening on port 8080 !")
}
