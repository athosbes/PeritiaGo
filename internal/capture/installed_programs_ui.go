package capture

import (
	"log"
	"path/filepath"
	"time"

	"github.com/athosbes/PeritiaGo/internal/executor"
)

// OpenAppWizAndCapture opens the Add/Remove Programs control panel,
// waits for it to render, and takes a screenshot.
func OpenAppWizAndCapture(outputsDir string) (string, error) {
	log.Println("Opening Control Panel: appwiz.cpl...")
	// control appwiz.cpl via secure executor
	_, err := executor.Start("control", "appwiz.cpl")
	if err != nil {
		return "", err
	}

	// Wait 5 seconds for the window to open and populate
	log.Println("Waiting 5 seconds for window to populate...")
	time.Sleep(5 * time.Second)

	screenshotPath := filepath.Join(outputsDir, "screenshots", "programas_instalados.png")
	err = CaptureScreen(screenshotPath)
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
	// control system via secure executor
	_, err := executor.Start("control", "system")
	if err != nil {
		return "", err
	}

	// Wait 5 seconds for the window to open and populate
	log.Println("Waiting 5 seconds for system info window to populate...")
	time.Sleep(5 * time.Second)

	screenshotPath := filepath.Join(outputsDir, "screenshots", "dados_maquina.png")
	err = CaptureScreen(screenshotPath)
	if err != nil {
		log.Printf("Failed to capture system info screen: %v", err)
		return "", err
	}

	return screenshotPath, nil
}
