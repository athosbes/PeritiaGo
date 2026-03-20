package artifacts

import (
	"log"
	"os"
	"path/filepath"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// SearchResiduals checks for left-over directories of uninstalled software.
// It searches common paths referencing software names.
func SearchResiduals(softwareNames []string) []models.Artifact {
	var artifacts []models.Artifact
	appData, _ := os.UserConfigDir()

	// We also search common roots
	roots := []string{
		filepath.Join(os.Getenv("SystemDrive")+"\\", "Program Files"),
		filepath.Join(os.Getenv("SystemDrive")+"\\", "Program Files (x86)"),
		os.Getenv("ProgramData"),
		appData,
	}

	for _, term := range softwareNames {
		if term == "" {
			continue
		}

		for _, root := range roots {
			target := filepath.Join(root, term)
			info, err := os.Stat(target)
			if err == nil {
				artifacts = append(artifacts, models.Artifact{
					Name:        term,
					Type:        "ResidualFile",
					Path:        target,
					Description: "Evidence of uninstalled or existing software remaining folder",
					Timestamp:   info.ModTime().Format("2006-01-02 15:04:05"),
				})
			}
		}
	}
	log.Printf("Found %d residual directories for provided terms\n", len(artifacts))
	return artifacts
}
