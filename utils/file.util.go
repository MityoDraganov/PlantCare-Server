package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"PlantCare/types"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DownloadFile downloads a file from the given URL and saves it to the specified filePath.
func DownloadFile(url, filePath string) error {
	// Perform the HTTP GET request to download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: received HTTP %d", resp.StatusCode)
	}

	// Create the output file where the downloaded content will be saved
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the downloaded content to the output file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Unzip extracts a ZIP file located at src into the destination directory dest.
func Unzip(src, dest string) error {
	// Open the ZIP file
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Loop through each file in the ZIP archive
	for _, f := range r.File {
		fPath := filepath.Join(dest, f.Name)

		// If it's a directory, create it
		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		// Ensure the directory exists for the file
		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		// Open the file for writing
		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Open the file inside the ZIP for reading
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Copy the content from the ZIP file to the output file
		if _, err = io.Copy(outFile, rc); err != nil {
			return err
		}
	}
	return nil
}
func FindSrcDir(driverExtractDir, driverURL string) (string, error) {
	// Extract the last segment of the URL, which is the subfolder name
	urlParts := strings.Split(strings.TrimRight(driverURL, "/"), "/")
	if len(urlParts) == 0 {
		return "", fmt.Errorf("invalid URL: unable to determine driver subdirectory from %s", driverURL)
	}

	subDirName := urlParts[len(urlParts)-1]
	// Construct the path to the expected driver directory
	driverDirPath := filepath.Join(driverExtractDir, subDirName)
	driverDirPath = driverDirPath + "-main"

	var foundFilePath string
	err := filepath.WalkDir(driverDirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Check if the file has .cpp or .c extension
		if !d.IsDir() && (strings.HasSuffix(d.Name(), ".cpp") || strings.HasSuffix(d.Name(), ".c")) {
			foundFilePath = path
			return filepath.SkipDir // Stop walking once we find the file
		}
		return nil
	})
	fmt.Println("foundFilePath")
	fmt.Println(foundFilePath)

	if err != nil {
		return "", fmt.Errorf("error searching for driver file in %s: %v", driverDirPath, err)
	}

	if foundFilePath == "" {
		return "", fmt.Errorf("no .cpp or .c file found in expected path %s", driverDirPath)
	}

	return foundFilePath, nil
}

// CopyDriverLibrary copies the driver library to the specified lib directory.
func CopyDriverLibrary(srcDir, destDir string) error {
	srcDir = filepath.Clean(srcDir)
	destDir = filepath.Clean(destDir)

	err := filepath.WalkDir(srcDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(srcDir, path)
		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		input, err := os.Open(path)
		if err != nil {
			return err
		}
		defer input.Close()

		output, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer output.Close()

		_, err = io.Copy(output, input)
		return err
	})
	return err
}

func FindClassName(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the entire file content
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Regular expression to match a class declaration in C++
	classRegex := regexp.MustCompile(`\bclass\s+([A-Za-z_]\w*)\b`)
	// Additional regex to match a class name based on member function definitions with ClassName::
	memberFunctionRegex := regexp.MustCompile(`\b([A-Za-z_]\w*)::`)

	// First, try to find a class declaration using the original regex
	matches := classRegex.FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return matches[1], nil
	}

	// If no class declaration was found, try finding a member function definition to infer the class name
	memberMatches := memberFunctionRegex.FindStringSubmatch(string(content))
	if len(memberMatches) > 1 {
		return memberMatches[1], nil
	}

	// If no class declaration or member function was found
	return "", fmt.Errorf("no class declaration or member function found in the file")
}


// WriteConfigJSON creates a JSON file in the desired structure
func WriteConfigJSON(configPath string, sensorDriverConfig map[string]string) error {
	// Create the output file
	configFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// Construct the config object
	config := types.Config{}

	// Sample serial numbers (you may adjust based on your logic)
	for serialNumber, className := range sensorDriverConfig {
		sensor := types.Sensor{
			SerialNumber: serialNumber,
			Type:         className,
		}
		// Add sensor to the config
		config.Sensors = append(config.Sensors, sensor)
	}

	// Adding a control entry, as a placeholder (you can populate it dynamically as needed)
	control := types.Control{
		SerialNumber: "",
		Type:         "WaterPump",
		DependantSensor: struct {
			SerialNumber string `json:"serialNumber"`
			MinValue     int    `json:"minValue"`
			MaxValue     int    `json:"maxValue"`
		}{
			SerialNumber: "YKTMgxAKCwE5jNXo", // Example, adjust based on your logic
			MinValue:     0,
			MaxValue:     50,
		},
	}

	// Add control to the config
	config.Controls = append(config.Controls, control)

	// Create a JSON encoder with indentation
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")

	// Write the config structure to the JSON file
	if err := encoder.Encode(config); err != nil {
		return err
	}
	return nil
}