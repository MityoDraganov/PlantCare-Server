package utils

import (
	"PlantCare/types"
	"PlantCare/utils/firebaseUtil"
	"PlantCare/websocket/connectionManager"
	"PlantCare/websocket/wsDtos"
	"PlantCare/websocket/wsTypes"
	wsutils "PlantCare/websocket/wsUtils"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// main repo url
const repoURL = "https://github.com/MityoDraganov/PlantCare-esp32/archive/refs/heads/controls.zip"

func UploadMultipleDrivers(sensorDriverURLs map[string]string, driverConfig []types.DriverConfig , potConn *wsTypes.Connection) error {
	driverZipFilePath := "./driver.zip"
	repoZipFilePath := "./repo.zip"
	driverExtractDir := "./extracted/repo/PlantCare-esp32-controls/src/drivers"
	repoExtractDir := "./extracted/repo"
	configJsonDir := "./extracted/repo/PlantCare-esp32-controls/src"
	sensorDriverConfig := make(map[string]string)
	controlConfig := []types.DriverJsonConfig{}
	// --SENSORS--

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
	for serialNumber, driverURL := range sensorDriverURLs {
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

	// --CONTROLS--
	

	for _, control := range driverConfig {
		zipUrl, err := convertGitHubURLToZip(control.DriverURL, "main")
		if err != nil {
			return err
		}

		// Download and extract the control files
		if err := DownloadFile(zipUrl, driverZipFilePath); err != nil {
			return err
		}
		if err := Unzip(driverZipFilePath, driverExtractDir); err != nil {
			return err
		}

		// Find the control source directory and class name
		controlFilePath, err := FindSrcDir(driverExtractDir, control.DriverURL)
		if err != nil {
			return err
		}

		className, err := FindClassName(controlFilePath)
		if err != nil {
			return err
		}

		// Store the control configuration based on the control serial number
		controlConfig = append(controlConfig, types.DriverJsonConfig{
			SerialNumber:    control.SerialNumber,
			DriverURL:       control.DriverURL,
			DependantSensor: control.DependantSensorSerial,
			MinValue:        control.MinValue,
			MaxValue:        control.MaxValue,
			Classname:       className,
		})
	}





	// Write the configuration JSON file
	configPath := filepath.Join(configJsonDir, "config.json")
	if err := WriteConfigJSON(configPath, sensorDriverConfig, controlConfig); err != nil {
		return err
	}


	// Prepare the PlatformIO OTA command
	firmwarePath := filepath.Join(repoExtractDir, "PlantCare-esp32-controls")
	cmd := exec.Command("pio", "run")
	cmd.Dir = firmwarePath
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("OTA upload failed: %s\nstdout: %s\nstderr: %s", err.Error(), out.String(), errOut.String())
	}

	destinationPath := "firmwareUpdates/" + time.Now().Format("20060102150405") + ".bin"
	file, err := os.Open(firmwarePath + "/.pio/build/esp32dev/firmware.bin")
	if err != nil {
		return err
	}
	firmwareUrl, err := firebaseUtil.UploadFile(file, destinationPath)
	if err != nil {
		return err
	}

	message := wsDtos.FirmwareCommand{
		Command:     string(wsTypes.FirmwareUpdate),
		DownloadUrl: firmwareUrl,
	}

	notification := wsDtos.NotificationDto{
		Data: message,
	}

	fmt.Println("Sending firmware message")
	if err := wsutils.SendMessage(potConn, "", wsTypes.FirmwareUpdate, notification); err != nil {
		fmt.Println("Failed to send firmware update message:", err)
		return err
	}

	connectionManager.ConnManager.RemoveConnectionByInstance(potConn)
	return nil
}

// UploadDriver handles the upload and processing of the driver
func UploadDriver(GitURL string, potIdStr string) *error {

	zipUrl, err := convertGitHubURLToZip(GitURL, "controls")
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
	driverExtractDir := "./extracted/repo/PlantCare-esp32-controls/lib"
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
	if err := DownloadFile("https://github.com/MityoDraganov/PlantCare-esp32/archive/refs/heads/controls.zip", repoZipFilePath); err != nil {
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
