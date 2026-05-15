package capture

import (
	"fmt"
	"os"
	"time"

	"github.com/athosbes/PeritiaGo/internal/models"
	"golang.org/x/sys/windows/registry"
)

// GetSystemStatus gathers information about OS, version, architecture and status.
func GetSystemStatus() []models.Artifact {
	var arts []models.Artifact
	now := time.Now().Format(time.RFC3339)

	// OS Info from Registry
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		pn, _, _ := k.GetStringValue("ProductName")
		arts = append(arts, models.Artifact{
			Name:        "Operating System",
			Type:        "SystemStatus",
			Path:        "Registry",
			Description: "OS Product Name",
			Value:       pn,
			Timestamp:   now,
		})

		ver, _, _ := k.GetStringValue("DisplayVersion")
		if ver == "" {
			ver, _, _ = k.GetStringValue("ReleaseId")
		}
		arts = append(arts, models.Artifact{
			Name:        "OS Version",
			Type:        "SystemStatus",
			Path:        "Registry",
			Description: "OS Release/Version ID",
			Value:       ver,
			Timestamp:   now,
		})

		ubr, _, _ := k.GetIntegerValue("UBR")
		arts = append(arts, models.Artifact{
			Name:        "OS Build Revision (UBR)",
			Type:        "SystemStatus",
			Path:        "Registry",
			Description: "Update Build Revision",
			Value:       fmt.Sprintf("%d", ubr),
			Timestamp:   now,
		})
	}

	// Environment
	arts = append(arts, models.Artifact{
		Name:        "Architecture",
		Type:        "SystemStatus",
		Path:        "Environment",
		Description: "System Architecture",
		Value:       os.Getenv("PROCESSOR_ARCHITECTURE"),
		Timestamp:   now,
	})

	return arts
}

// Win32_QuickFixEngineering represents the WMI class for Windows updates.
type Win32_QuickFixEngineering struct {
	Caption     string
	Description string
	HotFixID    string
	InstalledOn string
}

func GetWindowsUpdates() []models.Artifact {
	var arts []models.Artifact
	now := time.Now().Format(time.RFC3339)

	var dst []Win32_QuickFixEngineering
	query := "SELECT Caption, Description, HotFixID, InstalledOn FROM Win32_QuickFixEngineering"
	err := QueryWMI(query, &dst)
	if err != nil {
		arts = append(arts, models.Artifact{
			Name:        "Windows Updates (Hotfixes)",
			Type:        "SystemStatus",
			Path:        "WMI QuickFixEngineering",
			Description: "Failed to enumerate installed Windows Updates via WMI.",
			Value:       err.Error(),
			Timestamp:   now,
		})
		return arts
	}

	for _, q := range dst {
		desc := fmt.Sprintf("%s - %s (Installed: %s)", q.Caption, q.Description, q.InstalledOn)
		arts = append(arts, models.Artifact{
			Name:        q.HotFixID,
			Type:        "WindowsUpdate",
			Path:        q.Caption,
			Description: desc,
			Value:       "Installed",
			Timestamp:   now,
		})
	}

	if len(dst) == 0 {
		arts = append(arts, models.Artifact{
			Name:        "Windows Updates (Hotfixes)",
			Type:        "SystemStatus",
			Path:        "WMI QFE",
			Description: "No Windows Updates/Hotfixes found.",
			Value:       "Not Found",
			Timestamp:   now,
		})
	}
	return arts
}
