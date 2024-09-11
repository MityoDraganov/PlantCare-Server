package controllers

import (
	"PlantCare/utils" // Make sure to import the utils package
	"PlantCare/websocket/connectionManager"
	wsutils "PlantCare/websocket/wsUtils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// RequestBody defines the structure of the request body
type RequestBody struct {
	GitURL string `json:"gitUrl"`
}

// UploadDriver handles the upload and processing of the driver
func UploadDriver(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the URL
	if reqBody.GitURL == "" {
		http.Error(w, "Git URL is required", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	potIdStr := params["potId"]

	// Define file paths
	driverZipFilePath := "driver.zip"
	repoZipFilePath := "repo.zip"
	driverExtractDir := "./extracted/repo/PlantCare-esp32-main/lib"
	repoExtractDir := "./extracted/repo"

	// Send firmware data via WebSocket
	connManager := connectionManager.GetInstance()

	// Get all connections
	connections := connManager.GetAllConnections()

	// Debugging: Print all connection keys (pot IDs)
	fmt.Println("Connections keys (pot IDs):")
	for key := range connections {
		fmt.Println("Key:", key)
	}

	// Clean up previous extraction directories if they exist
	if err := os.RemoveAll(driverExtractDir); err != nil {
		http.Error(w, "Failed to clean up previous driver extraction directory", http.StatusInternalServerError)
		return
	}
	if err := os.RemoveAll(repoExtractDir); err != nil {
		http.Error(w, "Failed to clean up previous repository extraction directory", http.StatusInternalServerError)
		return
	}

	// Create extraction directories
	if err := os.MkdirAll(driverExtractDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create driver extraction directory", http.StatusInternalServerError)
		return
	}
	if err := os.MkdirAll(repoExtractDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create repository extraction directory", http.StatusInternalServerError)
		return
	}

	// Download the repository ZIP file
	if err := utils.DownloadFile("https://github.com/MityoDraganov/PlantCare-esp32/archive/refs/heads/main.zip", repoZipFilePath); err != nil {
		http.Error(w, "Failed to download repository ZIP file", http.StatusInternalServerError)
		return
	}

	// Download the driver ZIP file
	if err := utils.DownloadFile(reqBody.GitURL, driverZipFilePath); err != nil {
		http.Error(w, "Failed to download driver ZIP file", http.StatusInternalServerError)
		return
	}

	// Unzip the downloaded files
	if err := utils.Unzip(repoZipFilePath, repoExtractDir); err != nil {
		http.Error(w, "Failed to unzip repository file", http.StatusInternalServerError)
		return
	}

	if err := utils.Unzip(driverZipFilePath, driverExtractDir); err != nil {
		http.Error(w, "Failed to unzip driver file", http.StatusInternalServerError)
		return
	}

	// Build the project using PlatformIO
	if err := utils.BuildProject(repoExtractDir); err != nil {
		http.Error(w, "Failed to build project", http.StatusInternalServerError)
		return
	}

	// Locate the compiled firmware
	// Locate the compiled firmware
	firmwarePath := filepath.Join(repoExtractDir, "PlantCare-esp32-main", ".pio", "build", "esp32dev", "firmware.bin")

	// Debugging: Print the firmware path for confirmation
	fmt.Printf("Looking for firmware at: %s\n", firmwarePath)

	// Check if the file exists before reading
	if _, err := os.Stat(firmwarePath); os.IsNotExist(err) {
		http.Error(w, "Compiled firmware not found", http.StatusInternalServerError)
		return
	}

	// Read the compiled firmware
	firmwareData, err := os.ReadFile(firmwarePath)
	if err != nil {
		http.Error(w, "Failed to read compiled firmware", http.StatusInternalServerError)
		return
	}

	connection, ok := connManager.GetConnection(potIdStr)
	if !ok {
		http.Error(w, "Target connection not found", http.StatusNotFound)
		return
	}
	wsutils.SendFirmwareUpdate(connection, firmwareData)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Driver uploaded and processed successfully!")
}
