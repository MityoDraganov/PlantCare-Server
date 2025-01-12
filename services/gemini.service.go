package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeneratePrediction makes a prediction based on the provided prompt and input data.
func GeneratePrediction(prompt string, inputData interface{}) (*string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	var inputJSON string
	if inputData != nil {
		inputBytes, err := json.MarshalIndent(inputData, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input data: %w", err)
		}
		inputJSON = string(inputBytes)
	}

	// Generate content using the model with the JSON input
	model := client.GenerativeModel("gemini-1.5-flash")
	var promptWithInput string
	if inputJSON != "" {
		promptWithInput = prompt + "\nData:\n" + inputJSON
	} else {
		promptWithInput = prompt
	}
	resp, err := model.GenerateContent(ctx, genai.Text(promptWithInput))
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
