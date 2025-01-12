package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

// Struct to hold structured data for the model
type ModelOutput struct {
	Status            string `json:"status"`
	PercentageHealthy uint8  `json:"percentageHealthy"`
}

// Predict makes a prediction based on the provided input data and image file.
func PredictPlantHealth(imagePath string) (*string, error) {

	// Read and encode the image file
	imageData, err := readAndEncodeImage(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read and encode image: %w", err)
	}

	// Create the input JSON with both text and image data
	inputJSON, err := json.Marshal(map[string]interface{}{
		"image": imageData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %w", err)
	}

	promt := "retrieve the output in a json format {status: \"healthy\" or \"stressed\", percentageHealthy: 0 to 100}" + string(inputJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	resp, err := GeneratePrediction(promt, nil)

	if resp == nil {
		return nil, fmt.Errorf("no response received from model")
	}

	return resp, nil
}

// readAndEncodeImage reads an image file and encodes it to a base64 string.
func readAndEncodeImage(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
