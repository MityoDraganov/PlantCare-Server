package controllers

import (
	"PlantCare/initPackage"
	"PlantCare/models"
)

func FindSensorBySerialNum(serialNumber string)(*models.Sensor, error) {
	var sensorDbObject models.Sensor
	result := initPackage.Db.Where(&models.Sensor{SerialNumber: serialNumber}).First(&sensorDbObject)

	if result.Error != nil {
        return nil, result.Error
    }

    return &sensorDbObject, nil
}