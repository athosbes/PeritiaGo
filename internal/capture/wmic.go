package capture

import (
	"encoding/csv"
	"log"
	"path/filepath"

	"github.com/athosbes/PeritiaGo/internal/filesystem"
)

// Win32_Product represents the WMI class for installed software.
type Win32_Product struct {
	Name        string
	Vendor      string
	Version     string
	InstallDate string
}

// CaptureInstalledSoftwareWMI captures installed software via native WMI (replacing WMIC).
func CaptureInstalledSoftwareWMI(outputsDir string) (string, error) {
	log.Println("Capturing installed software via Native WMI...")

	var dst []Win32_Product
	query := "SELECT Name, Vendor, Version, InstallDate FROM Win32_Product"
	err := QueryWMI(query, &dst)
	if err != nil {
		log.Printf("[Warning] WMI software capture failed: %v", err)
		return "", err
	}

	csvPath := filepath.Join(outputsDir, "softwares_wmi.csv")
	f, err := filesystem.Create(csvPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Write Header
	writer.Write([]string{"Name", "Vendor", "Version", "InstallDate"})

	for _, p := range dst {
		writer.Write([]string{p.Name, p.Vendor, p.Version, p.InstallDate})
	}

	return csvPath, nil
}
