package controllers

import (
	"PlantCare/dtos"
	"PlantCare/models"
)

func ToDriverDTO(driver models.Driver) dtos.DriverDto {
	return dtos.DriverDto{
		DownloadUrl: driver.DownloadUrl,
	}
}
