package controllers

import (
    "PlantCare/dtos"
    "PlantCare/models"
    "encoding/json"
    "fmt"
    "net/http"
)

func ClerkUserRegister(w http.ResponseWriter, r *http.Request) {
    var clerkResponse dtos.ClerkCreateUserDto
    err := json.NewDecoder(r.Body).Decode(&clerkResponse)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Println("response")
    fmt.Println(clerkResponse)

	fmt.Println("response 1")
    //fmt.Println(&clerkResponse.ID)

    user := models.User{
        ClerkID: clerkResponse.Data.ID,
        IsAdmin_: false,
    }
    userDbObject := db.Create(&user)

    if userDbObject.Error != nil {
        fmt.Println(userDbObject.Error.Error())
    }

    fmt.Println(&userDbObject)
}