package capture

import (
	"log"
	"strings"

	"github.com/athosbes/PeritiaGo/internal/models"
	"golang.org/x/sys/windows/registry"
)

// GetInstalledSoftware looks through common registry locations for installed programs.
func GetInstalledSoftware() []models.Software {
	var softwares []models.Software

	keysToSearch := []struct {
		Key  registry.Key
		Path string
		Arch string
	}{
		{registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\Uninstall`, "x64"},
		{registry.LOCAL_MACHINE, `Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, "x86"},
		{registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Uninstall`, "User"},
	}

	for _, ks := range keysToSearch {
		k, err := registry.OpenKey(ks.Key, ks.Path, registry.READ|registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			log.Printf("Failed to open key %s: %v", ks.Path, err)
			continue
		}

		subkeys, err := k.ReadSubKeyNames(-1)
		k.Close()
		if err != nil {
			continue
		}

		for _, subkeyName := range subkeys {
			subPath := ks.Path + `\` + subkeyName
			subK, err := registry.OpenKey(ks.Key, subPath, registry.QUERY_VALUE)
			if err != nil {
				continue
			}

			displayName, _, err := subK.GetStringValue("DisplayName")
			if err != nil || displayName == "" {
				subK.Close()
				continue
			}

			displayVersion, _, _ := subK.GetStringValue("DisplayVersion")
			publisher, _, _ := subK.GetStringValue("Publisher")
			installDate, _, _ := subK.GetStringValue("InstallDate")
			installLocation, _, _ := subK.GetStringValue("InstallLocation")
			uninstallString, _, _ := subK.GetStringValue("UninstallString")
			productID, _, _ := subK.GetStringValue("ProductID")

			msiGUID := ""
			if strings.HasPrefix(subkeyName, "{") && strings.HasSuffix(subkeyName, "}") {
				msiGUID = subkeyName
			}

			softwares = append(softwares, models.Software{
				DisplayName:     displayName,
				DisplayVersion:  displayVersion,
				Publisher:       publisher,
				InstallDate:     installDate,
				InstallLocation: installLocation,
				UninstallString: uninstallString,
				Architecture:    ks.Arch,
				MSIGUID:         msiGUID,
				ProductID:       productID,
				Source:          "Registry: " + subPath,
			})
			subK.Close()
		}
	}
	return softwares
}
