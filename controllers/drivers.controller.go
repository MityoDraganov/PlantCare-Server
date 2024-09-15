package controllers

import (
	"PlantCare/utils" // Make sure to import the utils package
	"PlantCare/websocket/connectionManager"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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



	connection, ok := connManager.GetConnection(potIdStr)
	if !ok {
		http.Error(w, "Target connection not found", http.StatusNotFound)
		return
	}
	if err := uploadFirmwareOTA(repoExtractDir, connection.IP); err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload firmware: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Driver uploaded and processed successfully!")
}

func uploadFirmwareOTA(repoExtractDir string, esp32IP string) error {
	firmwarePath := filepath.Join(repoExtractDir, "PlantCare-esp32-main")
	fmt.Println(firmwarePath)
	// Extract the first part of the IP address
	ipParts := strings.Split(esp32IP, ":")

	// Hardcode the port to 8266
	otaAddress := ipParts[0] + ":8266"
	fmt.Println(otaAddress)

	// Prepare the PlatformIO OTA command
	cmd := exec.Command("pio", "run", "-e", "esp32dev_ota", "--target", "upload", "-v")
	cmd.Dir = firmwarePath

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("OTA upload failed: %s\nstdout: %s\nstderr: %s", err.Error(), out.String(), errOut.String())
	}

	fmt.Printf("OTA upload output: %s\n", out.String())
	return nil
}
