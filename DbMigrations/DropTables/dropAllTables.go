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
		&models.User{},
		&models.CropPot{},
		&models.Sensor{}, 
		&models.Measurement{}, 
		&models.Control{}, 
		&models.Webhook{}, 
		&models.Update{},
		&models.ActivePeriod{},
		&models.Condition{},
	}

	for _, model := range modelsToDrop {
		if err := db.Migrator().DropTable(model); err != nil {
			return err
		}
		log.Printf("Dropped table for model: %T", model)
	}

	if err := db.Migrator().DropTable("webhook_sensors"); err != nil {
		return err
	}
	log.Println("Dropped table: webhook_sensors")

	return nil
}

func main() {
	db := InitDBDrop()

	if err := DropAllTables(db); err != nil {
		log.Fatal("failed to drop all tables:", err)
	}

	log.Println("All tables dropped successfully!")
}
