package capture

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// CaptureWMIC runs wmic product get /format:csv and saves it to a file.
func CaptureWMIC(outputsDir string) (string, error) {
	log.Println("Capturing installed software via WMIC...")
	
	// We use the full command to get all bits
	cmd := exec.Command("wmic", "product", "get", "/format:csv")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("[Warning] WMIC capture failed: %v", err)
		return "", err
	}

	outputPath := filepath.Join(outputsDir, "wmic_products.csv")
	err = os.WriteFile(outputPath, output, 0644)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
