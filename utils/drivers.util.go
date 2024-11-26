package utils

import (
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// main repo url
const repoURL = "https://github.com/MityoDraganov/PlantCare-esp32/archive/refs/heads/production.zip"
const firmwareUpdateURL = "https://firebasestorage.googleapis.com/v0/b/plantcare-436309.appspot.com/o/firmwareUpdates%2Ffirmware.bin?alt=media&token=3a373153-8f1c-467e-b829-9c323b5de3b1"

func UploadMultipleDrivers(driverURLs map[string]string, potConn *wsTypes.Connection) error {
	driverZipFilePath := "./driver.zip"
	repoZipFilePath := "./repo.zip"
	driverExtractDir := "./extracted/repo/PlantCare-esp32-production/src/drivers"
	repoExtractDir := "./extracted/repo"
	configJsonDir := "./extracted/repo/PlantCare-esp32-production/src"
	sensorDriverConfig := make(map[string]string)

	// Clean up any previous artifacts
	if err := cleanUp(driverExtractDir, repoExtractDir, driverZipFilePath, repoZipFilePath); err != nil {
		return err
	}

	if err := os.MkdirAll(driverExtractDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(repoExtractDir, os.ModePerm); err != nil {
		return err
	}

	if err := DownloadFile(repoURL, repoZipFilePath); err != nil {
		return err
	}

	if err := Unzip(repoZipFilePath, repoExtractDir); err != nil {
		return err
	}

	// Process each driver URL
	for serialNumber, driverURL := range driverURLs {
		zipUrl, err := convertGitHubURLToZip(driverURL, "main")
		if err != nil {
			return err
		}

		if err := DownloadFile(zipUrl, driverZipFilePath); err != nil {
			return err
		}

		if err := Unzip(driverZipFilePath, driverExtractDir); err != nil {
			return err
		}

		// Find the driver source directory containing the C++ class
		driverFilePath, err := FindSrcDir(driverExtractDir, driverURL)
		if err != nil {
			return err
		}

		className, err := FindClassName(driverFilePath)
		if err != nil {
			return err
		}

		// Store the configuration based on the serial number
		sensorDriverConfig[serialNumber] = className
	}

	// Write the configuration JSON file
	configPath := filepath.Join(configJsonDir, "config.json")
	if err := WriteConfigJSON(configPath, sensorDriverConfig); err != nil {
		return err
	}

	message := wsDtos.FirmwareCommand{
		Command: string(wsTypes.FirmwareUpdate),
		DownloadUrl: firmwareUpdateURL,
    }

	fmt.Println("Sending firmware message")
	if err := wsutils.SendMessage(potConn, "", wsTypes.FirmwareUpdate, message); err != nil {
        fmt.Println("Failed to send firmware update message:", err)
        return err
    }

	connectionManager.ConnManager.RemoveConnectionByInstance(potConn)
	return nil
}

// UploadDriver handles the upload and processing of the driver
func UploadDriver(GitURL string, potIdStr string) *error {

	zipUrl, err := convertGitHubURLToZip(GitURL, "production")
	if err != nil {
		return &err
	}

	// connection, ok := connectionManager.ConnManager.GetConnection(potIdStr)
	// if !ok {
	// 	err := errors.New("connection not found")
	// 	fmt.Println(err)
	// 	return &err
	// }

	driverZipFilePath := "./driver.zip"
	repoZipFilePath := "./repo.zip"
	driverExtractDir := "./extracted/repo/PlantCare-esp32-production/lib"
	repoExtractDir := "./extracted/repo"

	if err := cleanUp(driverExtractDir, repoExtractDir, driverZipFilePath, repoZipFilePath); err != nil {
		return &err
	}

	// Create extraction directories
	if err := os.MkdirAll(driverExtractDir, os.ModePerm); err != nil {
		return &err

	}
	if err := os.MkdirAll(repoExtractDir, os.ModePerm); err != nil {
		return &err

	}

	// Download the repository ZIP file
	if err := DownloadFile("https://github.com/MityoDraganov/PlantCare-esp32/archive/refs/heads/production.zip", repoZipFilePath); err != nil {
		return &err
	}

	// Download the driver ZIP file
	if err := DownloadFile(zipUrl, driverZipFilePath); err != nil {
		return &err

	}

	// Unzip the downloaded files
	if err := Unzip(repoZipFilePath, repoExtractDir); err != nil {
		return &err

	}

	if err := Unzip(driverZipFilePath, driverExtractDir); err != nil {
		fmt.Println(err)
		return &err

	}

	// if err := uploadFirmwareOTA(repoExtractDir, connection.IP); err != nil {
	// 	return &err
	// }

	if err := cleanUp(driverExtractDir, repoExtractDir, driverZipFilePath, repoZipFilePath); err != nil {
		return &err
	}

	return nil
}

func convertGitHubURLToZip(gitURL string, branch string) (string, error) {
	parsedURL, err := url.Parse(gitURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Ensure the URL is a GitHub repo URL
	if !strings.Contains(parsedURL.Host, "github.com") {
		return "", fmt.Errorf("invalid GitHub URL")
	}

	// Split the URL path to extract owner and repository
	pathParts := strings.Split(parsedURL.Path, "/")
	if len(pathParts) < 3 {
		return "", fmt.Errorf("URL should be in the format https://github.com/<owner>/<repo>")
	}

	owner := pathParts[1]
	repo := pathParts[2]

	// Default to "main" branch if no branch is specified
	if branch == "" {
		branch = "main"
	}

	// Construct the zip download URL for the specified branch
	zipURL := fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.zip", owner, repo, branch)
	return zipURL, nil
}

func cleanUp(paths ...string) error {
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("Error removing %s: %v\n", path, err)
			return err
		}
	}
	return nil
}
