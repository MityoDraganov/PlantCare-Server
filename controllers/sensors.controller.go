package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"time"
)

func UpdateSensor(w http.ResponseWriter, r *http.Request) {
	var sensorDtos []dtos.SensorUserRequestDto

	// Decode request body into an array of sensorDtos
	err := json.NewDecoder(r.Body).Decode(&sensorDtos)
	if err != nil {
		log.Println(err) // Use log.Println for non-fatal errors
		utils.JsonError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	tx := initPackage.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			utils.JsonError(w, fmt.Sprintf("Transaction error: %v", r), http.StatusInternalServerError)
		}
	}()

	for _, sensorDto := range sensorDtos {
		sensorDbObject, err := findSensorById(uint(sensorDto.ID))
		if err != nil {
			tx.Rollback()
			log.Printf("Sensor not found: %v", err) // Improved logging
			utils.JsonError(w, err.Error(), http.StatusNotFound)
			return
		}
	
		var intervalUpdateTime time.Duration
		if sensorDto.MeasurementInterval != "" {
			t, err := time.Parse("15:04", sensorDto.MeasurementInterval)
			if err != nil {
				tx.Rollback()
				log.Printf("Invalid measurement interval format: %s", err.Error()) // Improved logging
				utils.JsonError(w, fmt.Sprintf("Invalid measurement interval format: %s", err.Error()), http.StatusBadRequest)
				return
			}
			intervalUpdateTime = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
		}
	
		sensorUpdate := models.Sensor{
			Alias:              &sensorDto.Alias, 
			Description:        sensorDto.Description,
			MeasuremntInterval: intervalUpdateTime,
		}
	
		result := tx.Model(sensorDbObject).Updates(sensorUpdate).Clauses(clause.Returning{})
		if result.Error != nil {
			tx.Rollback()
			log.Printf("Failed to update sensor: %v", result.Error) // Improved logging
			utils.JsonError(w, result.Error.Error(), http.StatusBadRequest)
			return
		}
	
		// Handle DriverUrl updates safely
		if sensorDto.DriverUrl != "" {
			log.Println("Updating Driver URL...") // Log before attempting to update driver
			if sensorDbObject.Driver == nil {
				log.Println("Driver not found, creating new driver.")
				driver := models.Driver{
					
					DownloadUrl: sensorDto.DriverUrl,
				}
				if err := tx.Create(&driver).Error; err != nil {
					log.Printf("Failed to create new driver: %v", err) // Improved logging
					tx.Rollback()
					utils.JsonError(w, "Failed to create new driver", http.StatusInternalServerError)
					return
				}
				fmt.Printf("Driver object: %+v\n", driver)

				sensorDbObject.Driver = &driver
				tx.Save(sensorDbObject)
			} else {
				log.Println("Driver found, updating URL.")
				sensorDbObject.Driver.DownloadUrl = sensorDto.DriverUrl
				if err := tx.Save(sensorDbObject.Driver).Error; err != nil {
					log.Printf("Failed to update driver URL: %v", err)
					tx.Rollback()
					utils.JsonError(w, "Failed to update driver URL", http.StatusInternalServerError)
					return
				}
			}
		}
	}


	if err := tx.Commit().Error; err != nil {
		utils.JsonError(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	// Return the updated array of sensorDtos
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sensorDtos); err != nil {
		utils.JsonError(w, err.Error(), http.StatusInternalServerError)
	}
}


func AddSensor(potId uint, sensorDto dtos.AttachSensor) (*models.Sensor, *error) {
	sensor := models.Sensor{
		CropPotID:    potId,
		SerialNumber: sensorDto.SerialNumber,
		IsOfficial:   false,

		Alias:       sensorDto.Alias,
		Description: sensorDto.Description,

		MeasuremntInterval: time.Hour,
	}

	result := initPackage.Db.Create(&sensor).Clauses(clause.Returning{})
	if result.Error != nil {
		return nil, &result.Error
	}

	return &sensor, nil
}

func GetMeasurementsBySensorId(id uint) dtos.SensorMeasurementsSummary {
	sensor, err := findSensorById(id)
	if err != nil {
		log.Fatal(err)
	}

	SensorMeasurementsSummaryDto := dtos.SensorMeasurementsSummary{
		SensorType:   sensor.Type,
		Measurements: sensor.Measurements,
	}
	return SensorMeasurementsSummaryDto
}

func FindSensorBySerialNum(serialNumber string) (*models.Sensor, error) {
	var sensorDbObject models.Sensor
	result := initPackage.Db.Where(&models.Sensor{SerialNumber: serialNumber}).First(&sensorDbObject)

	if result.Error != nil {
		return nil, result.Error
	}

	return &sensorDbObject, nil
}
func findSensorById(id uint) (*models.Sensor, error) {
	var sensor models.Sensor
	result := initPackage.Db.Preload("Measurements").Preload("Driver").First(&sensor, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sensor, nil
}

// Maps a single Sensor to SensorDto
func MapSensorToDTO(sensor models.Sensor) dtos.SensorDto {
	var driverUrl string

	if sensor.Driver != nil {
		driverUrl = sensor.Driver.DownloadUrl
	}

	return dtos.SensorDto{
		ID:                  sensor.ID,
		SerialNumber:        sensor.SerialNumber,
		Alias:               *sensor.Alias,
		Description:         sensor.Description,
		MeasurementInterval: utils.DurationToTimeString(sensor.MeasuremntInterval),
		Measurements:        sensor.Measurements,
		DriverUrl:           driverUrl,
	}
}

// Converts a single Sensor or a slice of Sensors to a slice of SensorDto
func ToSensorsDTO(input interface{}) []dtos.SensorDto {
	switch v := input.(type) {
	case models.Sensor:
		// If it's a single sensor, wrap it in a slice
		return []dtos.SensorDto{MapSensorToDTO(v)}
	case []models.Sensor:
		// If it's a slice of sensors, map each sensor to SensorDto
		sensorDTOs := make([]dtos.SensorDto, len(v))
		for i, sensor := range v {
			sensorDTOs[i] = MapSensorToDTO(sensor)
		}
		return sensorDTOs
	default:
		// Handle unexpected types by returning an empty slice
		return []dtos.SensorDto{}
	}
}

func ChangeAttachedState(sensor *models.Sensor) error {
	// Toggle the IsAttached field
	sensor.IsAttached = !sensor.IsAttached

	// Update the database with the new state
	if err := initPackage.Db.Save(sensor).Error; err != nil {
		return err
	}

	return nil
}

func FindDriverBySensorId(sensorId uint) (*models.Driver, error) {
	var driver models.Driver
	result := initPackage.Db.First(&driver, "sensor_id = ?", sensorId)

	if result.Error != nil {
		return nil, result.Error
	}

	return &driver, nil
}
