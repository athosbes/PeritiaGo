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

// SearchResidualsByTerms checks for left-over directories of uninstalled software based on search terms.
func SearchResidualsByTerms(softwareNames []string) []models.Artifact {
	var artifacts []models.Artifact
	appData, _ := os.UserConfigDir()

	roots := []string{
		os.Getenv("ProgramFiles"),
		os.Getenv("ProgramFiles(x86)"),
		os.Getenv("ProgramData"),
		appData,
	}

	for _, term := range softwareNames {
		if term == "" {
			continue
		}

		for _, root := range roots {
			if root == "" {
				continue
			}
			target := filepath.Join(root, term)
			info, err := os.Stat(target)
			if err == nil {
				artifacts = append(artifacts, models.Artifact{
					Name:        term,
					Type:        "ResidualFile",
					Path:        target,
					Description: "Evidence of uninstalled or existing software remaining folder",
					Timestamp:   info.ModTime().Format(time.RFC3339),
				})
			}
		}
	}
	log.Printf("Found %d residual directories for provided terms\n", len(artifacts))
	return artifacts
}
