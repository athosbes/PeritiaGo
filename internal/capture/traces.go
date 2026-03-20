package capture

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// SearchResidualTraces looks for folders that might belong to uninstalled software.
func SearchResidualTraces() []models.Artifact {
	var traces []models.Artifact

	dirsToScan := []string{
		os.Getenv("ProgramFiles"),
		os.Getenv("ProgramFiles(x86)"),
		os.Getenv("ProgramData"),
		filepath.Join(os.Getenv("AppData"), "..", "Local"),
		os.Getenv("AppData"),
	}

	// This is a simple heuristic: folders that don't match common system folders
	// or are known to be residual. In a real scenario, this would be compared
	// against the list of currently installed software.

	for _, root := range dirsToScan {
		if root == "" {
			continue
		}

		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				path := filepath.Join(root, entry.Name())
				// For now, we just record all non-hidden folders in these locations
				// as potential software traces for the peritus to analyze.
				// In a refined version, we'd cross-reference with 'GetInstalledSoftware'

				info, err := entry.Info()
				if err != nil {
					continue
				}

				traces = append(traces, models.Artifact{
					Name:        entry.Name(),
					Type:        "ResidualFolder",
					Path:        path,
					Description: "Potential residual or active software folder",
					Timestamp:   info.ModTime().Format(time.RFC3339),
				})
			}
		}
	}

	log.Printf("Found %d potential residual traces\n", len(traces))
	return traces
}
