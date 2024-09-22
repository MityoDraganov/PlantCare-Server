package dtos

// ForecastDTO represents both the historical and forecasted weather data
type ForecastDTO struct {
	Location Location    `json:"location"`
	Current  *Current    `json:"current,omitempty"` // Optional for historical data
	Forecast ForecastDay `json:"forecast"`
}

// Location represents the location details in the weather response
type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TimezoneID     string  `json:"tz_id"`
	LocalTimeEpoch int64   `json:"localtime_epoch"`
	LocalTime      string  `json:"localtime"`
}

// Current represents the current weather conditions (for forecast data only)
type Current struct {
	LastUpdatedEpoch int64     `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	UV               float64   `json:"uv"`
}

// Condition describes the weather condition details
type Condition struct {
	Text string `json:"text"`
}

// ForecastDay represents the forecast data for a given date
type ForecastDay struct {
	Forecastday []DayForecast `json:"forecastday"`
}

// DayForecast contains daily weather data for both historical and forecast
type DayForecast struct {
	Date       string     `json:"date"`
	DateEpoch  int64      `json:"date_epoch"`
	DayDetails DayDetails `json:"day"`
}

// DayDetails contains summarized data for the day
type DayDetails struct {
	MaxTempC    float64 `json:"maxtemp_c"`
	MinTempC    float64 `json:"mintemp_c"`
	AvgTempC    float64 `json:"avgtemp_c"`
	AvgHumidity float64 `json:"avghumidity"`
}

// GeminiRequest contains past indoor, past outdoor, and future outdoor data
type GeminiRequest struct {
	PastIndoors    SensorMeasurementsSummary `json:"past_indoors"`
	PastOutdoors   ForecastDTO          `json:"past_outdoors"`
	FutureOutdoors ForecastDTO          `json:"future_outdoors"`
}
