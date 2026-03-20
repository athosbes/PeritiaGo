package capture

import (
	"log"

	"github.com/athosbes/PeritiaGo/internal/models"
	"golang.org/x/sys/windows/registry"
)

// GetLicenseData attempts to retrieve license information for the OS and common softwares.
func GetLicenseData() []models.LicenseData {
	var licenses []models.LicenseData

	// Windows License Info
	winLicense := models.LicenseData{SoftwareName: "Windows Operating System"}

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\SoftwareProtectionPlatform`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		pid, _, _ := k.GetStringValue("BackupProductKeyDefault")
		if pid != "" {
			winLicense.ProductKey = pid
			licenses = append(licenses, winLicense)
		}
	}

	// Office and others could be added here

	log.Printf("Captured %d license records\n", len(licenses))
	return licenses
}
