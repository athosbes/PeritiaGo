package capture

import (
	"fmt"
	"time"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// GetRunningProcesses lists currently active volatile memory processes.
// This is critical for catching portable, unauthorized, or unlicensed software running
// actively in RAM that has not left explicit Registry Install traces.
// Win32_Process represents the WMI class for system processes.
type Win32_Process struct {
	Caption        string
	ExecutablePath *string // Use pointer because it can be null for system processes
	ProcessId      uint32
}

// GetRunningProcesses lists currently active volatile memory processes.
func GetRunningProcesses() []models.Artifact {
	var arts []models.Artifact
	now := time.Now().Format(time.RFC3339)

	var dst []Win32_Process
	query := "SELECT Caption, ExecutablePath, ProcessId FROM Win32_Process"
	err := QueryWMI(query, &dst)
	if err != nil {
		arts = append(arts, models.Artifact{
			Name:        "Volatile Processes (RAM)",
			Type:        "MemoryProcess",
			Path:        "Memory",
			Description: "Failed to enumerate running processes via WMI.",
			Value:       err.Error(),
			Timestamp:   now,
		})
		return arts
	}

	for _, p := range dst {
		procPath := "Path Restricted/System"
		if p.ExecutablePath != nil {
			procPath = *p.ExecutablePath
		}

		arts = append(arts, models.Artifact{
			Name:        p.Caption,
			Type:        "MemoryProcess",
			Path:        procPath,
			Description: fmt.Sprintf("Live Execution PID: %d", p.ProcessId),
			Value:       "Running in Volatile Memory",
			Timestamp:   now,
		})
	}

	return arts
}
