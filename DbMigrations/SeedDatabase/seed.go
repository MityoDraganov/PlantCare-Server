package main

import (
	"log"
	"time"

	"PlantCare/models"

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
    err := db.AutoMigrate(
		&models.User{}, 
		&models.CropPot{}, 
		&models.Sensor{}, 
		&models.Measurement{}, 
		&models.ControlSettings{}, 
		&models.Webhook{},
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

    users := []models.User{
        {ClerkID: "user_2jod4hRuJ9nqUIzftpaTTNWLVxv", IsAdmin_: false},
    }
    if err := db.Create(&users).Error; err != nil {
        return err
    }

    // Create ControlSettings
    controlSettings := []models.ControlSettings{
        {WateringInterval: 30},
        {WateringInterval: 60},
    }
    if err := db.Create(&controlSettings).Error; err != nil {
        return err
    }

	timeNow := time.Now()

    // Create CropPots
    cropPots := []models.CropPot{
        {
            Token:            "pot_1",
            Alias:            "Herb Garden",
            LastWateredAt:    &timeNow,
            IsArchived:       false,
            ClerkUserID:      &users[0].ClerkID,
            ControlSettingsID: &controlSettings[0].ID,
        },
        {
            Token:            "pot_2",
            Alias:            "Vegetable Bed",
            LastWateredAt:    &timeNow,
            IsArchived:       false,
            ClerkUserID:      &users[0].ClerkID,
            ControlSettingsID: &controlSettings[1].ID,
        },
    }
    if err := db.Create(&cropPots).Error; err != nil {
        return err
    }

    // Create Sensors
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

    // Create Measurements
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

    return nil
}

func main() {
	db := InitDBSeed()

	if err := SeedDatabase(db); err != nil {
		log.Fatal("failed to seed database:", err)
	}

	log.Println("Database seeded successfully!")
}
