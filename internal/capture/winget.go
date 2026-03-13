package capture

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// CaptureWinget runs winget list and saves it to a file.
func CaptureWinget(outputsDir string) (string, error) {
	log.Println("Capturing installed software via Winget...")
	
	// winget list --nowarn --ignore-warnings
	cmd := exec.Command("winget", "list", "--nowarn", "--ignore-warnings")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("[Warning] Winget capture failed: %v", err)
		return "", err
	}

	outputPath := filepath.Join(outputsDir, "winget_list.csv")
	err = os.WriteFile(outputPath, output, 0644)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
