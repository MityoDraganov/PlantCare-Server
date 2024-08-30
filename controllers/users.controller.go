package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/utils"

	"PlantCare/models"
	"encoding/json"
	"net/http"
)

func ClerkUserRegister(w http.ResponseWriter, r *http.Request) {
    var clerkResponse dtos.ClerkCreateUserDto
    err := json.NewDecoder(r.Body).Decode(&clerkResponse)
    if err != nil {
        utils.JsonError(w, err.Error(), http.StatusBadRequest)
        return
    }

    user := models.User{
        ClerkID: clerkResponse.Data.ID,
        IsAdmin_: false,
    }
    userDbObject := initPackage.Db.Create(&user)

    if userDbObject.Error != nil {
        utils.JsonError(w, userDbObject.Error.Error(), http.StatusInternalServerError)
    }
}