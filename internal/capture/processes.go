package capture

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// GetRunningProcesses lists currently active volatile memory processes.
// This is critical for catching portable, unauthorized, or unlicensed software running
// actively in RAM that has not left explicit Registry Install traces.
func GetRunningProcesses() []models.Artifact {
	var arts []models.Artifact
	now := time.Now().Format(time.RFC3339)

	// wmic process get Caption,ExecutablePath,ProcessId /format:csv
	out, err := exec.Command("wmic", "process", "get", "Caption,ExecutablePath,ProcessId", "/format:csv").Output()
	if err != nil {
		arts = append(arts, models.Artifact{
			Name:        "Volatile Processes (RAM)",
			Type:        "MemoryProcess",
			Path:        "Memory",
			Description: "Failed to enumerate running processes.",
			Value:       err.Error(),
			Timestamp:   now,
		})
		return arts
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Node") {
			continue
		}

		parts := strings.SplitN(line, ",", 4)
		if len(parts) >= 4 {
			// Expected format from wmic csv: Node,Caption,ExecutablePath,ProcessId
			procName := parts[1]
			procPath := parts[2]
			procID := parts[3]

			if procPath == "" {
				procPath = "Path Restricted/System"
			}

			arts = append(arts, models.Artifact{
				Name:        procName,
				Type:        "MemoryProcess",
				Path:        procPath,
				Description: fmt.Sprintf("Live Execution PID: %s", procID),
				Value:       "Running in Volatile Memory",
				Timestamp:   now,
			})
		}
	}

	return arts
}
