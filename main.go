package main

import (
	"log"
	"net/http"

	"TravelBuddy/controllers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Pre-populate with some users (optional)

	r.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{username}", controllers.GetUser).Methods("GET")
	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/users/{username}", controllers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{username}", controllers.DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
