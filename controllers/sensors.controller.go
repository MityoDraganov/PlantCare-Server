package controllers

import (
	"PlantCare/dtos"
	"PlantCare/initPackage"
	"PlantCare/models"
	"PlantCare/types"
	"PlantCare/utils"
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/otaManager"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v2"
	"gorm.io/gorm/clause"
)

type SensorConfig struct {
	SerialNumber string `json:"serialNumber"`
	DriverURL    string `json:"driverUrl"`
}



type UpdateDto struct {
	SensorDtos  []dtos.SensorUserRequestDto `json:"sensorDtos"`
	ControlDtos []dtos.ControlUserRequestDto `json:"controlDtos"`
}

func UpdateSensor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating sensor...")
	var updateDto UpdateDto
	driverURLs := make(map[string]string)
	var potId uint
	var sensorConfigs []SensorConfig
	var driverConfig []types.DriverConfig

	err := json.NewDecoder(r.Body).Decode(&updateDto)
	if err != nil {
		log.Println(err)
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

	for _, sensorDto := range updateDto.SensorDtos {
		sensorDbObject, err := findSensorById(uint(sensorDto.ID))
		if err != nil {
			tx.Rollback()
			log.Printf("Sensor not found: %v", err)
			utils.JsonError(w, err.Error(), http.StatusNotFound)
			return
		}

		sensorUpdate := models.Sensor{
			Alias:       &sensorDto.Alias,
			Description: sensorDto.Description,
		}

		result := tx.Model(sensorDbObject).Updates(sensorUpdate).Clauses(clause.Returning{})
		if result.Error != nil {
			tx.Rollback()
			log.Printf("Failed to update sensor: %v", result.Error)
			utils.JsonError(w, result.Error.Error(), http.StatusBadRequest)
			return
		}

		if sensorDto.DriverUrl != "" {
			log.Println("Updating Driver URL...")
			if sensorDbObject.Driver == nil {
				log.Println("Driver not found, creating new driver.")
				driver := models.Driver{
					DownloadUrl: sensorDto.DriverUrl,
				}
				if err := tx.Create(&driver).Error; err != nil {
					log.Printf("Failed to create new driver: %v", err)
					tx.Rollback()
					utils.JsonError(w, "Failed to create new driver", http.StatusInternalServerError)
					return
				}
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
			sensorConfigs = append(sensorConfigs, SensorConfig{
				SerialNumber: sensorDbObject.SerialNumber,
				DriverURL:    sensorDto.DriverUrl,
			})
			driverURLs[sensorDbObject.SerialNumber] = sensorDto.DriverUrl
			potId = sensorDbObject.CropPotID
		}
	}


	fmt.Println("Starting control updates...")
	for _, controlDto := range updateDto.ControlDtos {
		controlDbObject, err := findControlById(uint(controlDto.ID))
		if err != nil {
			tx.Rollback()
			log.Printf("Control not found: %v", err)
			utils.JsonError(w, err.Error(), http.StatusNotFound)
			return
		}
	
		// Update control information
		controlUpdate := models.Control{
			Alias:       controlDto.Alias,
			Description: controlDto.Description,
			MinValue: controlDto.MinValue,
			MaxValue: controlDto.MaxValue,
			DriverUrl: controlDto.DriverUrl,
			DependantSensorSerial: &controlDto.DependantSensorSerial,
		}
	
		// Apply updates to control
		result := tx.Model(controlDbObject).Updates(controlUpdate)//.Clauses(clause.Returning{})
		if result.Error != nil {
			tx.Rollback()
			log.Printf("Failed to update control: %v", result.Error)
			utils.JsonError(w, result.Error.Error(), http.StatusBadRequest)
			return
		}
	
		// Handle Driver URL update
		if controlDto.DriverUrl != "" {
			log.Println("Updating Driver URL for control...")
			if controlDbObject.Driver == nil {
				log.Println("Driver not found, creating new driver for control.")
				driver := models.Driver{
					DownloadUrl: controlDto.DriverUrl,
				}
				if err := tx.Create(&driver).Error; err != nil {
					log.Printf("Failed to create new driver for control: %v", err)
					tx.Rollback()
					utils.JsonError(w, "Failed to create new driver for control", http.StatusInternalServerError)
					return
				}
				controlDbObject.Driver = &driver
			} else {
				log.Println("Driver found, updating URL for control.")
				controlDbObject.Driver.DownloadUrl = controlDto.DriverUrl
			}
	
			// Save the control with updated driver
			if err := tx.Save(controlDbObject).Error; err != nil {
				log.Printf("Failed to save control with updated driver: %v", err)
				tx.Rollback()
				utils.JsonError(w, "Failed to save control with updated driver", http.StatusInternalServerError)
				return
			}
	
			// Update driver config
			driverConfig = append(driverConfig, types.DriverConfig{
				SerialNumber: controlDbObject.SerialNumber,
				DriverURL:    controlDto.DriverUrl,
				MinValue:     *controlDto.MinValue,
				MaxValue:     *controlDto.MaxValue,
	
				DependantSensorSerial: types.DependantSensor{
					SerialNumber: controlDto.DependantSensorSerial,
				},
			})
			potId = controlDbObject.CropPotID
	
			log.Println("Control driverConfig updated")
		}
	}
	

	go func() {
		claims, ok := clerk.SessionClaimsFromContext(r.Context())
		if !ok {
			fmt.Println("Error extracting session claims")
			utils.JsonError(w, "Unauthorized: unable to extract session claims", http.StatusUnauthorized)
			return
		}
		clerkUserID := claims.Subject
		userConn, isExisting := connectionManager.ConnManager.GetConnectionByKey(clerkUserID)

		potIDStr := strconv.FormatUint(uint64(potId), 10)
		connection, ok := connectionManager.ConnManager.GetConnection(potIDStr)
		if !ok {
			err := errors.New("connection not found for pot ID: " + potIDStr + "! Adding update to pendings.")
			if isExisting {
				wsutils.SendMessage(userConn, "", wsTypes.AsyncError, err)
			}
			otaManager.OTAManager.AddOTAPending(potIDStr, driverURLs)
			return
		}
		if isExisting {
			wsutils.SendMessage(userConn, "", wsTypes.AsyncPromise, nil)
		}

		if err := utils.UploadMultipleDrivers(driverURLs, driverConfig, connection); err != nil {
			log.Printf("Failed to upload driver: %v", err)

			if isExisting {
				wsutils.SendMessage(userConn, "", wsTypes.AsyncError, err)
			}
			return
		}

	}()

	// Commit the transaction after the OTA operation
	if err := tx.Commit().Error; err != nil {
		utils.JsonError(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	// Return the updated array of sensorDtos
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updateDto); err != nil {
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
		fmt.Println(err)
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
		ID:           sensor.ID,
		SerialNumber: sensor.SerialNumber,
		Alias:        utils.CoalesceString(sensor.Alias),
		Description:  sensor.Description,
		Measurements: sensor.Measurements,
		DriverUrl:    driverUrl,
		IsAttached:   sensor.IsAttached,
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

func AttachedStateUpdater(entity interface{}, state bool) error {
	switch v := entity.(type) {
	case *models.Sensor:
		v.IsAttached = state
	case *models.Control:
		v.IsAttached = state
	default:
		return errors.New("unsupported entity type")
	}

	// Update the database with the new state
	if err := initPackage.Db.Save(entity).Error; err != nil {
		return err
	}

	return nil
}

func FindDriverBySensorId(sensorId uint) (*models.Driver, error) {
	var driver models.Driver

	// Use Preload to load sensors associated with the driver
	result := initPackage.Db.Preload("Sensors").Joins("JOIN sensors ON sensors.driver_id = drivers.id").Where("sensors.id = ?", sensorId).First(&driver)

	if result.Error != nil {
		return nil, result.Error
	}

	return &driver, nil
}
