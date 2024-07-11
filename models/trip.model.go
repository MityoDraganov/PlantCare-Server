package models


//user model
type Trip struct {
	driver User
	currentPassengers []User

	startPoint string `json:"startPoint"`
	endPoint string `json:"endPoint"`
	maxAvaliableSeats int `json:"maxAvaliableSeats"`
	description string  `json:"description"`
}