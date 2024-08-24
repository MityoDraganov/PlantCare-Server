// seed.go
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
    err := db.AutoMigrate(&models.User{}, &models.CropPot{}, &models.SensorData{}, &models.ControlSettings{})
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

    // Create SensorData
    measurement := []models.Measurement{
        {
            Temperature: 22.5,
            Moisture:    75.0,
            WaterLevel:  50.0,
            SunExposure: 80.0,
            CropPotID:   cropPots[0].ID,
        },
        {
            Temperature: 21.0,
            Moisture:    65.0,
            WaterLevel:  45.0,
            SunExposure: 70.0,
            CropPotID:   cropPots[1].ID,
        },
    }
    if err := db.Create(&sensorData).Error; err != nil {
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
