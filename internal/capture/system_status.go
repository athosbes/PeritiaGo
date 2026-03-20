package capture

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/athosbes/PeritiaGo/internal/models"
	"golang.org/x/sys/windows/registry"
)

// GetSystemStatus extracts system lifecycle details such as original installation time
// and checks for the existence of previous Windows installations (e.g., C:\Windows.old).
func GetSystemStatus() []models.Artifact {
	var arts []models.Artifact
	now := time.Now().Format(time.RFC3339)

	// 1. Check Original Install Date
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err == nil {
		installDateEpoch, _, err := k.GetIntegerValue("InstallDate")
		if err == nil && installDateEpoch > 0 {
			installTime := time.Unix(int64(installDateEpoch), 0)
			arts = append(arts, models.Artifact{
				Name:        "Windows Original Install Date",
				Type:        "SystemStatus",
				Path:        `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\InstallDate`,
				Description: "Date and time the current OS was originally installed or majorly upgraded.",
				Value:       installTime.Format(time.RFC3339),
				Timestamp:   now,
			})
		} else {
			arts = append(arts, models.Artifact{
				Name:        "Windows Original Install Date",
				Type:        "SystemStatus",
				Path:        `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\InstallDate`,
				Description: "Registry key not found or unreadable. System might be heavily anomalous.",
				Value:       "Not Found",
				Timestamp:   now,
			})
		}
		k.Close()
	}

	// 2. Check for Windows.old (Indicates recent upgrade/reinstall without format)
	windowsOldPath := `C:\Windows.old`
	if info, err := os.Stat(windowsOldPath); err == nil && info.IsDir() {
		arts = append(arts, models.Artifact{
			Name:        "Previous Windows Installation Detected",
			Type:        "SystemStatus",
			Path:        windowsOldPath,
			Description: "Presence of C:\\Windows.old indicates the system was recently upgraded, reset, or reinstalled without formatting the drive.",
			Value:       fmt.Sprintf("Created: %s", info.ModTime().Format(time.RFC3339)),
			Timestamp:   now,
		})
	} else {
		arts = append(arts, models.Artifact{
			Name:        "Previous Windows Installation",
			Type:        "SystemStatus",
			Path:        windowsOldPath,
			Description: "No recent in-place upgrades or non-formatting reinstalls detected.",
			Value:       "Not Found",
			Timestamp:   now,
		})
	}

	// 3. System Reset Traces
	resetKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Reset`, registry.QUERY_VALUE)
	if err == nil {
		resetKey.Close()
		arts = append(arts, models.Artifact{
			Name:        "System Reset Traces",
			Type:        "SystemStatus",
			Path:        `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Reset`,
			Description: "Registry keys associated with a full system restoration or 'Reset this PC' feature exist.",
			Value:       "Present",
			Timestamp:   now,
		})
	} else {
		arts = append(arts, models.Artifact{
			Name:        "System Reset Traces",
			Type:        "SystemStatus",
			Path:        `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Reset`,
			Description: "No evidence of 'Reset this PC' feature usage found.",
			Value:       "Not Found",
			Timestamp:   now,
		})
	}

	// 4. Check for Hotfixes (Windows Updates)
	arts = append(arts, GetWindowsUpdates()...)

	return arts
}

func GetWindowsUpdates() []models.Artifact {
	var arts []models.Artifact
	now := time.Now().Format(time.RFC3339)

	out, err := exec.Command("wmic", "qfe", "get", "Caption,Description,HotFixID,InstalledOn", "/format:csv").Output()
	if err != nil {
		arts = append(arts, models.Artifact{
			Name:        "Windows Updates (Hotfixes)",
			Type:        "SystemStatus",
			Path:        "WMIC QFE",
			Description: "Failed to enumerate installed Windows Updates.",
			Value:       err.Error(),
			Timestamp:   now,
		})
		return arts
	}

	lines := strings.Split(string(out), "\n")
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Node") {
			continue // skip empty lines and the CSV header
		}

		parts := strings.SplitN(line, ",", 5)
		if len(parts) >= 5 {
			desc := fmt.Sprintf("%s - %s (Installed: %s)", parts[2], parts[3], parts[4])
			arts = append(arts, models.Artifact{
				Name:        parts[3], // HotFixID like KB123456
				Type:        "WindowsUpdate",
				Path:        parts[1], // Caption/Link
				Description: desc,
				Value:       "Installed",
				Timestamp:   now,
			})
			count++
		}
	}

	if count == 0 {
		arts = append(arts, models.Artifact{
			Name:        "Windows Updates (Hotfixes)",
			Type:        "SystemStatus",
			Path:        "WMIC QFE",
			Description: "No Windows Updates/Hotfixes found.",
			Value:       "Not Found",
			Timestamp:   now,
		})
	}
	return arts
}
