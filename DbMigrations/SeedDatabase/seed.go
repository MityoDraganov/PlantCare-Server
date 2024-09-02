package main

import (
	"PlantCare/models"
	"log"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func InitDBSeed() *gorm.DB {
	dsn := "sqlserver://server:P@ssw0rd123!@localhost:1433?database=Plant_Care"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}

func SeedDatabase(db *gorm.DB) error {
	// AutoMigrate the models to keep schema in sync
	err := db.AutoMigrate(
		&models.User{},
		&models.CropPot{},
		&models.Sensor{},
		&models.Measurement{},
		&models.Control{},
		&models.Update{},
		&models.Webhook{},
		&models.ActivePeriod{},
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// Seed Users
	users := []models.User{
		{ClerkID: "user_2jod4hRuJ9nqUIzftpaTTNWLVxv", IsAdmin_: false},
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}

	// Seed CropPots
	cropPots := []models.CropPot{
		{
			Token:       "pot_1",
			Alias:       "Herb Garden",
			IsArchived:  false,
			ClerkUserID: &users[0].ClerkID,
		},
		{
			Token:       "pot_2",
			Alias:       "Vegetable Bed",
			IsArchived:  false,
			ClerkUserID: &users[0].ClerkID,
		},
	}
	if err := db.Create(&cropPots).Error; err != nil {
		return err
	}

	// Seed Controls
	controls := []models.Control{
		{
			CropPotID:   cropPots[0].ID,
			SerialNumber: "ctrl_1",
			Alias:       "Water Pump",
			ActivePeriod: models.ActivePeriod{
				Start:     time.Duration(8)*time.Hour + time.Duration(30)*time.Minute, // 8:30 AM
				End:       time.Duration(17)*time.Hour, // 5:00 PM
			},
		},
		{
			CropPotID:   cropPots[1].ID,
			SerialNumber: "ctrl_2",
			Alias:       "Light Switch",
			ActivePeriod: models.ActivePeriod{
				Start:     time.Duration(6)*time.Hour, // 6:00 AM
			End:       time.Duration(20)*time.Hour, // 8:00 PM
			},
		},
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}

	// Seed Sensors
	sensors := []models.Sensor{
		{
			CropPotID:    cropPots[0].ID,
			SerialNumber: "sensor_1",
			Alias:        "Temperature Sensor",
			IsOfficial:   true,
		},
		{
			CropPotID:    cropPots[1].ID,
			SerialNumber: "sensor_2",
			Alias:        "Moisture Sensor",
			IsOfficial:   true,
		},
	}
	if err := db.Create(&sensors).Error; err != nil {
		return err
	}

	// Seed Measurements
	measurements := []models.Measurement{
		{
			SensorID: sensors[0].ID,
			Value:    22.5,
		},
		{
			SensorID: sensors[1].ID,
			Value:    65.0,
		},
	}
	if err := db.Create(&measurements).Error; err != nil {
		return err
	}

	// Seed Updates for Controls
	updates := []models.Update{
		{ControlID: controls[0].ID},
		{ControlID: controls[1].ID},
	}
	if err := db.Create(&updates).Error; err != nil {
		return err
	}

	subscribedEvents := []models.Sensor{
		sensors[0],
	}
	webhooks := []models.Webhook{
		{
			CropPotID:        cropPots[0].ID,
			EndpointUrl:      "https://webhook.site/7c7240ca-1f9d-4c84-b19e-2419c386d715",
			SubscribedEvents: subscribedEvents,
		},
	}
	if err := db.Create(&webhooks).Error; err != nil {
		return err
	}

	return nil
}

func main() {
	db := InitDBSeed()

	if err := SeedDatabase(db); err != nil {
		log.Fatal("failed to seed database:", err)
	}

	log.Println("Database seeded successfully!")
}
