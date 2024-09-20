package services

import (
	"PlantCare/dtos"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Fetch forecast data for the given location and number of days
func GetWeatherForecast(apiKey string, location string, days int) (*dtos.WeatherForecast, error) {
	// Construct the URL with placeholders
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=%d&aqi=no&alerts=no", apiKey, location, days)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch weather data: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into the WeatherForecast struct
	var forecast dtos.WeatherForecast
	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return nil, err
	}

	return &forecast, nil
}

// Estimate indoor temperature based on outdoor temperature and a known differential
func PredictIndoorTemperature(outdoorTemp float64, differential float64) float64 {
	return outdoorTemp + differential
}
