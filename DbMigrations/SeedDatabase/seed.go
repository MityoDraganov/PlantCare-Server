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
		&models.Condition{},
		&models.Control{},
		&models.Webhook{},
		&models.Update{},
		&models.ActivePeriod{},
		&models.Driver{},
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
				Start: time.Duration(8)*time.Hour + time.Duration(30)*time.Minute, // 8:30 AM
				End:   time.Duration(17) * time.Hour,                             // 5:00 PM
			},
		},
		{
			CropPotID:   cropPots[1].ID,
			SerialNumber: "ctrl_2",
			Alias:       "Light Switch",
			ActivePeriod: models.ActivePeriod{
				Start: time.Duration(6) * time.Hour, // 6:00 AM
				End:   time.Duration(20) * time.Hour, // 8:00 PM
			},
		},
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}

	// Seed Sensors
	sensors := []models.Sensor{
		{
			CropPotID:          cropPots[0].ID,
			SerialNumber:       "sensor_1",
			Alias:              "Temperature Sensor",
			IsOfficial:         true,
			MeasuremntInterval: time.Hour,
			Type: models.TempType,
		},
		{
			CropPotID:          cropPots[0].ID,
			SerialNumber:       "sensor_2",
			Alias:              "Moisture Sensor",
			IsOfficial:         true,
			MeasuremntInterval: 2 * time.Hour,
			Type: models.TempType,
		},
		{
			CropPotID:          cropPots[1].ID,
			SerialNumber:       "sensor_3",
			Alias:              "Light Sensor",
			IsOfficial:         true,
			MeasuremntInterval: time.Hour,
			Type: models.TempType,
		},
	}
	if err := db.Create(&sensors).Error; err != nil {
		return err
	}

	// Seed Conditions
	conditions := []models.Condition{
		{
			ControlID:         controls[0].ID,  // Water Pump control
			DependentSensorID: &sensors[1].ID,  // Moisture Sensor
			On:                30.0,            // If moisture is below 30%, turn on the water pump
			Off:               70.0,            // If moisture is above 70%, turn off the water pump
		},
		{
			ControlID:         controls[1].ID,  // Light Switch control
			DependentSensorID: &sensors[2].ID,  // Light Sensor
			On:                100.0,           // Turn on light if light sensor value is below 100
			Off:               300.0,           // Turn off light if light sensor value is above 300
		},
	}
	if err := db.Create(&conditions).Error; err != nil {
		return err
	}

	// Seed Measurements
	measurements := []models.Measurement{
		{
			SensorID: sensors[0].ID,
			Value:    22.5,
			CreatedAt: time.Now(),
		},
		{
			SensorID: sensors[0].ID,
			Value:    22.5,
			CreatedAt: time.Now(),
		},
		{
			SensorID: sensors[0].ID,
			Value:    21.5,
			CreatedAt: time.Now().AddDate(0, 0, -1), // 1 day ago
		},
		{
			SensorID: sensors[0].ID,
			Value:    22.0,
			CreatedAt: time.Now().AddDate(0, 0, -2),
		},
		{
			SensorID: sensors[0].ID,
			Value:    21,
			CreatedAt: time.Now().AddDate(0, 0, -3),
		},
		{
			SensorID: sensors[0].ID,
			Value:    23,
			CreatedAt: time.Now().AddDate(0, 0, -4),
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

	// Seed Webhooks
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
