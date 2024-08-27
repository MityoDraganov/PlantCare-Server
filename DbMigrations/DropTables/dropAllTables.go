package main

import (
	"log"

	"PlantCare/models"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func InitDBDrop() *gorm.DB {
	dsn := "sqlserver://server:P@ssw0rd123!@localhost:1433?database=Plant_Care"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}

func DropAllTables(db *gorm.DB) error {
	modelsToDrop := []interface{}{
		&models.Measurement{},
		&models.Sensor{},
		&models.Webhook{},
		&models.CropPot{},
		&models.ControlSettings{},
		&models.User{},
	}

	for _, model := range modelsToDrop {
		if err := db.Migrator().DropTable(model); err != nil {
			return err
		}
		log.Printf("Dropped table for model: %T", model)
	}

	return nil
}

func main() {
	db := InitDBDrop()

	if err := DropAllTables(db); err != nil {
		log.Fatal("failed to drop all tables:", err)
	}

	log.Println("All tables dropped successfully!")
}
