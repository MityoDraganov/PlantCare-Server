package services

import (
	"PlantCare/controllers"
	"PlantCare/dtos"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Fetch historical weather data for the given location and date
func getHistoricalWeather(apiKey string, location string, date time.Time) (*dtos.ForecastDTO, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/history.json?key=%s&q=%s&dt=%s",
		apiKey, location, date.Format("2006-01-02"))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch historical weather data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var historicalWeather dtos.ForecastDTO
	err = json.Unmarshal(body, &historicalWeather)
	if err != nil {
		return nil, err
	}

	return &historicalWeather, nil
}

// Fetch forecast data for the given location and number of days
func getWeatherForecast(apiKey string, location string) (*dtos.ForecastDTO, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=3&aqi=no&alerts=yes",
		apiKey, location)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch weather data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var forecast dtos.ForecastDTO
	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return nil, err
	}

	return &forecast, nil
}

func GetIndoorForecast(location string, userId string) (*string, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	// -1 = yesterday
	historicalDate := time.Now().AddDate(0, 0, -1)
	historicalOutdoorsWeather, err := getHistoricalWeather(apiKey, location, historicalDate)
	if err != nil {
		return nil, err
	}

	futureOutdoorsWeather, err := getWeatherForecast(apiKey, location)
	if err != nil {
		return nil, err
	}

	cropPots, err := controllers.FindPotsByUserId(userId)
	if err != nil {
		return nil, err
	}

	historicalIndoorsWeather := controllers.GetMeasurementsBySensorId(cropPots[0].Sensors[0].ID)

	// Call Predict with the correct structure
	indoorForecast, err := Predict(dtos.GeminiRequest{
		PastIndoors:    historicalIndoorsWeather,
		PastOutdoors:   *historicalOutdoorsWeather,
		FutureOutdoors: *futureOutdoorsWeather,
	})
	if err != nil {
		return nil, err
	}

	//return indoorForecast, nil
	return indoorForecast, nil
}
