package controllers

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var db *gorm.DB
var validate *validator.Validate

func SetDatabase(database *gorm.DB) {
	db = database
	validate = validator.New()
}
