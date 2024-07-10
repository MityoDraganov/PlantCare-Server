package models

//user model
type User struct {
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string
}