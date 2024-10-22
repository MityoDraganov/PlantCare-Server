package utils

import (
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsTypes"
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func UploadMultipleDrivers(driverURLs []string, potConn *wsTypes.Connection) error {
	// Prepare paths for drivers and repositories
	driverZipFilePath := "./driver.zip"
	repoZipFilePath := "./repo.zip"
	driverExtractDir := "./extracted/repo/PlantCare-esp32-main/src/drivers"
	repoExtractDir := "./extracted/repo"

	// Clean up any previous artifacts
	if err := cleanUp(driverExtractDir, repoExtractDir, driverZipFilePath, repoZipFilePath); err != nil {
		return err
	}

	// Create directories for extraction
	if err := os.MkdirAll(driverExtractDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(repoExtractDir, os.ModePerm); err != nil {
		return err
	}

	// Download the main repo ZIP
	if err := DownloadFile("https://github.com/MityoDraganov/PlantCare-esp32/archive/refs/heads/main.zip", repoZipFilePath); err != nil {
		return err
	}

	// Unzip the main repository
	if err := Unzip(repoZipFilePath, repoExtractDir); err != nil {
		return err
	}

	// Download and extract all the drivers sequentially for the single pot
	for _, driverURL := range driverURLs {
		// Convert the GitHub URL to the ZIP file download link
		zipUrl, err := convertGitHubURLToZip(driverURL)
		if err != nil {
			return err
		}

		// Download the driver ZIP file
		if err := DownloadFile(zipUrl, driverZipFilePath); err != nil {
			return err
		}

		// Extract the driver to the correct directory
		if err := Unzip(driverZipFilePath, driverExtractDir); err != nil {
			fmt.Println(err)
			return err
		}
	}



	// Perform the OTA upload using the repoExtractDir, which now includes all the drivers
	if err := uploadFirmwareOTA(repoExtractDir, potConn.IP); err != nil {
		customErr := errors.New("failed to upload driver OTA! Check pot connectivity")
		return customErr
	}

	// Clean up after the upload process
	if err := cleanUp(driverExtractDir, repoExtractDir, driverZipFilePath, repoZipFilePath); err != nil {
		return err
	}

	return nil
}



// UploadDriver handles the upload and processing of the driver
func UploadDriver(GitURL string, potIdStr string) *error {

	zipUrl, err := convertGitHubURLToZip(GitURL)
	if err != nil {
		return &err
	}

	connection, ok := connectionManager.ConnManager.GetConnection(potIdStr)
	if !ok {
		err := errors.New("connection not found")
		fmt.Println(err)
		return &err
	}


	driverZipFilePath := "./driver.zip"
	repoZipFilePath := "./repo.zip"
	driverExtractDir := "./extracted/repo/PlantCare-esp32-main/lib"
	repoExtractDir := "./extracted/repo"

	if err := cleanUp(driverExtractDir, repoExtractDir, driverZipFilePath, repoZipFilePath	); err != nil {
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
	if err := DownloadFile("https://github.com/MityoDraganov/PlantCare-esp32/archive/refs/heads/main.zip", repoZipFilePath); err != nil {
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


	if err := uploadFirmwareOTA(repoExtractDir, connection.IP); err != nil {
		return &err
	}

	if err := cleanUp(driverExtractDir, repoExtractDir, driverZipFilePath, repoZipFilePath	); err != nil {
        return &err
    }

	return nil
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

func convertGitHubURLToZip(gitURL string) (string, error) {
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

    // Construct the zip download URL
    zipURL := fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/main.zip", owner, repo)
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