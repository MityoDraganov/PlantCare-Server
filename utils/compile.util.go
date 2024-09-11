package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

// BuildProject builds the PlatformIO project located at projectDir.
func BuildProject(projectDir string) error {

	projectDir = filepath.Join(projectDir, "PlantCare-esp32-main")


    if err := printPIOVersion(); err != nil {
        return err
    }

    cmd := exec.Command("pio", "run")
    cmd.Dir = projectDir

    var out bytes.Buffer
    var errOut bytes.Buffer

    cmd.Stdout = &out
    cmd.Stderr = &errOut

    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("build failed: %s\nstdout: %s\nstderr: %s", err.Error(), out.String(), errOut.String())
    }

    fmt.Printf("Build output: %s\n", out.String())
    return nil
}

// printPIOVersion prints the version of PlatformIO for debugging purposes.
func printPIOVersion() error {
    cmd := exec.Command("pio", "--version")

    var out bytes.Buffer
    var errOut bytes.Buffer

    cmd.Stdout = &out
    cmd.Stderr = &errOut

    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("failed to get PlatformIO version: %s", errOut.String())
    }

    fmt.Printf("PlatformIO version: %s\n", out.String())
    return nil
}
