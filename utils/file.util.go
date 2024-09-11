package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
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

func FindSrcDir(rootDir string) (string, error) {
	var srcDir string

	err := filepath.WalkDir(rootDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && strings.HasSuffix(entry.Name(), "src") {
			srcDir = path
			return filepath.SkipDir // Stop searching as soon as we find the src directory
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	if srcDir == "" {
		return "", fmt.Errorf("src directory not found in %s", rootDir)
	}
	return srcDir, nil
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