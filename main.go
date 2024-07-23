package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"PlantCare/controllers"


	"PlantCare/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"


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
	r := mux.NewRouter()

	db := InitDB()
	controllers.SetDatabase(db)

	err := db.AutoMigrate(&models.User{}, &models.CropPot{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to SQL Server!")



	// Apply middleware to all other routes
	protectedRoutes := r.PathPrefix("/api/v1").Subrouter()

	authMiddleware := clerkhttp.WithHeaderAuthorization()
	protectedRoutes.Use(authMiddleware)

	// Crop Pot routes
	protectedRoutes.HandleFunc("/cropPots", controllers.GetCropPotsForUser).Methods("GET")
	protectedRoutes.HandleFunc("/cropPots/{token}", controllers.AssignCropPotToUser).Methods("POST")
	protectedRoutes.HandleFunc("/cropPots/{id}", controllers.UpdateCropPot).Methods("PUT")
	protectedRoutes.HandleFunc("/cropPots/{id}", controllers.RemoveCropPot).Methods("DELETE")
	
	//ADMIN ACTIONS
	r.HandleFunc("/cropPots", controllers.AddCropPot).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
	fmt.Println("Server listening on port 8080 !")
}
