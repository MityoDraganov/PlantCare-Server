package services

import (
	"PlantCare/dtos"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Struct to hold structured data for the model
type ModelInput struct {
	PastIndoors    dtos.SensorMeasurementsSummary `json:"past_indoors"`
	PastOutdoors   []OutdoorDay                   `json:"past_outdoors"`
	FutureOutdoors []OutdoorDay                   `json:"future_outdoors"`
}

// Struct to hold day-specific outdoor data
type OutdoorDay struct {
	Date        string  `json:"date"`
	MaxTempC    float64 `json:"max_temp_c"`
	MinTempC    float64 `json:"min_temp_c"`
	AvgTempC    float64 `json:"avg_temp_c"`
	AvgHumidity float64 `json:"avg_humidity"`
}

// Predict makes a prediction based on the provided input data.
func Predict(inputData dtos.GeminiRequest) (*string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Specify the model to use
	model := client.GenerativeModel("gemini-1.5-flash")

	// Convert past and future outdoor data to a structured format
	var pastOutdoors []OutdoorDay
	for _, forecast := range inputData.PastOutdoors.Forecast.Forecastday {
		pastOutdoors = append(pastOutdoors, OutdoorDay{
			Date:        forecast.Date,
			MaxTempC:    forecast.DayDetails.MaxTempC,
			MinTempC:    forecast.DayDetails.MinTempC,
			AvgTempC:    forecast.DayDetails.AvgTempC,
			AvgHumidity: forecast.DayDetails.AvgHumidity,
		})
	}

	var futureOutdoors []OutdoorDay
	for _, forecast := range inputData.FutureOutdoors.Forecast.Forecastday {
		futureOutdoors = append(futureOutdoors, OutdoorDay{
			Date:        forecast.Date,
			MaxTempC:    forecast.DayDetails.MaxTempC,
			MinTempC:    forecast.DayDetails.MinTempC,
			AvgTempC:    forecast.DayDetails.AvgTempC,
			AvgHumidity: forecast.DayDetails.AvgHumidity,
		})
	}

	// Create the final input structure for the model
	modelInput := ModelInput{
		PastIndoors:    inputData.PastIndoors,
		PastOutdoors:   pastOutdoors,
		FutureOutdoors: futureOutdoors,
	}

	// Convert the input to JSON format
	inputJSON, err := json.MarshalIndent(modelInput, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %w", err)
	}

	// Generate content using the model with the JSON input
	resp, err := model.GenerateContent(ctx, genai.Text((
		"Analyze the past indoor and outdoor temperature data and provide predictions for future indoor temperatures. " +
			"Format the response as a JSON array with each entry containing the date and a 'predictions' object that includes 'temp' as the predicted temperature. " +
			"Example format:\n[\n{\n\"date\": \"YYYY-MM-DD\",\n\"predictions\": {\n\"temp\": 22.5,\n(other values can follow)\n}\n}\n].\n" +
			"When generating the result, return only a single JSON in this format without any extra explanation." +
			"\nData:\n" + string(inputJSON))))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("no response received from model")
	}

	respString := extractResponseContent(resp)
	fmt.Println("Model Response:", respString)
	// Extract valid JSON from the response string
	validJSON, err := extractJSON(respString)
	if err != nil {
		return nil, fmt.Errorf("failed to extract valid JSON: %w", err)
	}

	return &validJSON, nil
}

// Extract content from the response
func extractResponseContent(resp *genai.GenerateContentResponse) string {
	var content strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content.WriteString(fmt.Sprintf("%v\n", part))
			}
		}
	}
	return content.String()
}

// Extract and validate JSON from a string
func extractJSON(input string) (string, error) {
	// Regular expression to match a valid JSON array or object
	jsonRegex := regexp.MustCompile(`(?s)(\[[\s\S]*?\{.*?\}[\s\S]*?\]|\{.*?\})`)


	// Find all JSON matches in the input string
	matches := jsonRegex.FindAllString(input, -1)

	if len(matches) == 0 {
		return "", fmt.Errorf("no valid JSON found in input")
	}

	// Assume the first match is the valid JSON we need
	jsonStr := strings.TrimSpace(matches[0])

	// Validate that the extracted string is valid JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return "", fmt.Errorf("invalid JSON format: %w", err)
	}

	return jsonStr, nil
}
