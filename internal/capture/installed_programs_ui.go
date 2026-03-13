package capture

import (
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

// OpenAppWizAndCapture opens the Add/Remove Programs control panel,
// waits for it to render, and takes a screenshot.
func OpenAppWizAndCapture(outputsDir string) (string, error) {
	log.Println("Opening Control Panel: appwiz.cpl...")
	cmd := exec.Command("control", "appwiz.cpl")
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Wait 5 seconds for the window to open and populate
	log.Println("Waiting 5 seconds for window to populate...")
	time.Sleep(5 * time.Second)

	screenshotPath := filepath.Join(outputsDir, "screenshots", "programas_instalados.png")
	err := CaptureScreen(screenshotPath)
	if err != nil {
		log.Printf("Failed to capture screen: %v", err)
		return "", err
	}

	return screenshotPath, nil
}

// OpenSystemInfoAndCapture opens the System Information control panel,
// waits for it to render, and takes a screenshot.
func OpenSystemInfoAndCapture(outputsDir string) (string, error) {
	log.Println("Opening System Information...")
	// 'control system' works on Win7, 10, 11
	cmd := exec.Command("control", "system")
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Wait 5 seconds for the window to open and populate
	log.Println("Waiting 5 seconds for system info window to populate...")
	time.Sleep(5 * time.Second)

	screenshotPath := filepath.Join(outputsDir, "screenshots", "dados_maquina.png")
	err := CaptureScreen(screenshotPath)
	if err != nil {
		log.Printf("Failed to capture system info screen: %v", err)
		return "", err
	}

	return screenshotPath, nil
}
