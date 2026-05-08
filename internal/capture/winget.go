package capture

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/athosbes/PeritiaGo/internal/executor"
	"github.com/athosbes/PeritiaGo/internal/filesystem"
)

// CaptureWinget runs winget list and saves it to a file.
func CaptureWinget(outputsDir string) (string, error) {
	log.Println("Capturing installed software via Winget...")

	// winget is not in SystemBinaries, we should use its absolute path if known
	// or find it securely. For now, we'll assume it might be in the path but
	// we'll try to find its absolute path to be safe.
	wingetPath := "winget"
	if p, err := exec.LookPath("winget"); err == nil {
		wingetPath = p
	}

	output, err := executor.Execute(wingetPath, "list", "--nowarn", "--ignore-warnings")
	if err != nil {
		log.Printf("[Warning] Winget capture failed: %v", err)
		return "", err
	}

	outputPath := filepath.Join(outputsDir, "winget_list.csv")
	err = filesystem.WriteFile(outputPath, output, 0644)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
