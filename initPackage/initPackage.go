package initPackage

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var Db *gorm.DB
var validate *validator.Validate

func SetDatabase(database *gorm.DB) {
	Db = database
	validate = validator.New()
}