package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"PlantCare/types"
	"image"
	"image/jpeg"
	"mime/multipart"

	"github.com/disintegration/imaging"
)

// Predict makes a prediction based on the provided input data and image file.
func PredictPlantHealth(image multipart.File) (*types.ModelOutput, error) {

	// Read and encode the image file
    processedImage, err := preProcessImage(image)
    if err != nil {
        return nil, fmt.Errorf("failed to pre-process image: %w", err)
    }

    // Encode the pre-processed image to a base64 string
    processedImageData, err := encodeImage(processedImage)
    if err != nil {
        return nil, fmt.Errorf("failed to encode pre-processed image: %w", err)
    }

	// Create the input JSON with both text and image data
	inputJSON, err := json.Marshal(map[string]interface{}{
		"image": processedImageData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %w", err)
	}

	promt := "Diagnose the plant. Retrieve the output in a json format {percentageHealthy: 0 to 100, plantName: string, certentyPercantage: 0 to 100, isPlantRecognised: boolean}" + string(inputJSON)

	resp, err := GeneratePrediction(promt, nil)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no response received from model")
	}

	var modelOutput types.ModelOutput
	err = json.Unmarshal([]byte(*resp), &modelOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal model output: %w", err)
	}

	return &modelOutput, nil
}

// readAndEncodeImage reads an image from a multipart.File and encodes it to a base64 string.
func encodeImage(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// Return the base64-encoded string
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}


func preProcessImage(file multipart.File) (image.Image, error) {
    // Decode the image file
    img, _, err := image.Decode(file)
    if err != nil {
        return nil, fmt.Errorf("failed to decode image: %w", err)
    }

    // Resize the image to a suitable dimension (e.g., 224x224)
    resizedImage := imaging.Resize(img, 224, 224, imaging.Lanczos)

    // Apply basic enhancements (optional)
    // You can explore techniques like brightness/contrast adjustment, noise reduction

    return resizedImage, nil
}