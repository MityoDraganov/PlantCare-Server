package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

// GeminiService represents a service for interacting with the Gemini API.
type GeminiService struct {
	apiEndpoint         string
	serviceAccountFile  string
	httpClient          *http.Client
}

// NewGeminiService creates a new instance of GeminiService.
func NewGeminiService(apiEndpoint, serviceAccountFile string) (*GeminiService, error) {
	httpClient := &http.Client{}
	return &GeminiService{
		apiEndpoint:        apiEndpoint,
		serviceAccountFile: serviceAccountFile,
		httpClient:         httpClient,
	}, nil
}

// GetAccessToken retrieves an access token using the service account.
func (g *GeminiService) GetAccessToken(ctx context.Context) (string, error) {
	credentials, err := idtoken.NewTokenSource(ctx, g.apiEndpoint, option.WithCredentialsFile(g.serviceAccountFile))
	if err != nil {
		return "", fmt.Errorf("failed to create token source: %w", err)
	}

	token, err := credentials.Token()
	if err != nil {
		return "", fmt.Errorf("failed to obtain token: %w", err)
	}

	return token.AccessToken, nil
}

// Predict makes a prediction based on the provided input data.
func (g *GeminiService) Predict(ctx context.Context, inputData interface{}) (interface{}, error) {
	accessToken, err := g.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"instances": inputData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %w", err)
	}

	req, err := http.NewRequest("POST", g.apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %d %s", resp.StatusCode, body)
	}

	var predictions interface{}
	if err := json.Unmarshal(body, &predictions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal predictions: %w", err)
	}

	return predictions, nil
}
