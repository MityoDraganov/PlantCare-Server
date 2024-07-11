package main

import (
	"fmt"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"TravelBuddy/controllers"
	"TravelBuddy/middlewares"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := mux.NewRouter()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://server:server123@cluster0.jmfnonn.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
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
