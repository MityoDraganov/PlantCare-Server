package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"

	"TravelBuddy/controllers"
    "TravelBuddy/middlewares"
)

func main() {
	r := mux.NewRouter()

	// Pre-populate with some users (optional)

	r.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{username}", controllers.GetUser).Methods("GET")
	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")

    updateUserHandler := middlewares.LoggerMiddleware(http.HandlerFunc(controllers.UpdateUser))
	r.Handle("/users/{username}", updateUserHandler).Methods("PUT")

    deleteUserHandler := middlewares.LoggerMiddleware(http.HandlerFunc(controllers.DeleteUser))
	r.Handle("/users/{username}", deleteUserHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
