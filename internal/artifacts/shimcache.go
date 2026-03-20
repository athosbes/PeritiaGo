package artifacts

import (
	"fmt"
	"log"

	"github.com/athosbes/PeritiaGo/internal/models"
	"golang.org/x/sys/windows/registry"
)

// ParseShimCache reads the AppCompatCache from the Registry.
func ParseShimCache() []models.Artifact {
	var artifacts []models.Artifact
	keyPath := `SYSTEM\CurrentControlSet\Control\Session Manager\AppCompatCache`

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
	if err != nil {
		log.Printf("Failed to open ShimCache key: %v", err)
		return artifacts
	}
	defer k.Close()

	val, _, err := k.GetBinaryValue("AppCompatCache")
	if err != nil {
		return artifacts
	}

	artifacts = append(artifacts, models.Artifact{
		Name:        "AppCompatCache",
		Type:        "ShimCache",
		Path:        `HKLM\` + keyPath,
		Description: "ShimCache binary data containing execution traces.",
		Value:       fmt.Sprintf("%d bytes of execution cache data found", len(val)),
	})

	log.Printf("ShimCache evidence found: %d bytes\n", len(val))
	return artifacts
}
