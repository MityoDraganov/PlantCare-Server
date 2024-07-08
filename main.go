package main

import (
    "encoding/json"
    "log"
    "math/rand"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    // Pre-populate with some users (optional)
    users = append(users, User{Username: "user1", Email: "user1@example.com", PasswordHash: hashPassword("password1")})
    users = append(users, User{Username: "user2", Email: "user2@example.com", PasswordHash: hashPassword("password2")})

    r.HandleFunc("/users", getUsers).Methods("GET")
    r.HandleFunc("/users/{username}", getUser).Methods("GET")
    r.HandleFunc("/users", createUser).Methods("POST")
    r.HandleFunc("/users/{username}", updateUser).Methods("PUT")
    r.HandleFunc("/users/{username}", deleteUser).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8080", r))
}
