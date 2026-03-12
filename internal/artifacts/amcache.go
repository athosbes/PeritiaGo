package artifacts

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// ParseAmcache executes Eric Zimmerman's AmcacheParser.exe if available.
// It parses the Amcache.hve file for evidence of execution.
func ParseAmcache(outputsDir string) []models.Artifact {
	var artifacts []models.Artifact
	
	amcacheHve := `C:\Windows\AppCompat\Programs\Amcache.hve`
	if _, err := os.Stat(amcacheHve); os.IsNotExist(err) {
		log.Println("Amcache.hve not found")
		return artifacts
	}

	// For a complete integration, we assume AmcacheParser.exe is in PATH or current dir.
	// We output its CSV to outputsDir/amcache/
	outPath := filepath.Join(outputsDir, "amcache")
	os.MkdirAll(outPath, 0755)

	log.Println("Running AmcacheParser.exe ...")
	cmd := exec.Command("AmcacheParser.exe", "-f", amcacheHve, "--csv", outPath)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to run AmcacheParser.exe (is it in your PATH?): %v\n", err)
		// We still record the artifact attempts
		artifacts = append(artifacts, models.Artifact{
			Name:        "Amcache.hve",
			Type:        "Amcache",
			Path:        amcacheHve,
			Description: "AmcacheParser execution failed: " + err.Error(),
			Timestamp:   time.Now().Format(time.RFC3339),
		})
		return artifacts
	}

	artifacts = append(artifacts, models.Artifact{
		Name:        "Amcache.hve",
		Type:        "Amcache",
		Path:        amcacheHve,
		Description: "Amcache parsed successfully to " + outPath,
		Timestamp:   time.Now().Format(time.RFC3339),
	})

	return artifacts
}
