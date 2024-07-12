package controllers

import (
	"TravelBuddy/dtos"
	"TravelBuddy/models"
	"encoding/json"
	"fmt"

	"net/http"

	"TravelBuddy/utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDto dtos.CreateUserDto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the DTO
	err = validate.Struct(userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if password and rePassword match
	if userDto.Password != userDto.RePassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Check password length
	if len(userDto.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Map DTO to model
	user := models.User{
		Username:     userDto.Username,
		Email:        userDto.Email,
		PasswordHash: string(hashedPassword),
	}

	// Save the user to the database
	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWT(user.Username, user.Email)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := dtos.AuthResponse{
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var userDto dtos.LoginUserDto
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := findUserByEmail(userDto.Email).First(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userDto.Password))
	if err != nil {
		fmt.Println("Password does not match.")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	result := db.Find(&users)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User
	result := db.First(&user, "username = ?", params["username"])
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Retrieve the existing user
	var user models.User
	result := findUserById(params["id"]).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Decode incoming data into a new user instance
	var userDto dtos.CreateUserDto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update fields if they are provided
	if userDto.Username != "" {
		user.Username = userDto.Username
	}
	if userDto.Email != "" {
		user.Email = userDto.Email
	}
	if userDto.Password != "" && userDto.Password == userDto.RePassword {
		if len(userDto.Password) < 8 {
			http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Save the updates to the database
	result = db.Save(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	result := db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

func findUserById(id string) *gorm.DB {
	var user models.User
	result := db.First(&user, "id = ?", id)

	return result
}

func findUserByEmail(email string) *gorm.DB {
	var user models.User
	result := db.First(&user, "email = ?", email)

	return result
}
