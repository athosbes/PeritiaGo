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

	// We output its CSV to outputsDir/amcache/
	outPath := filepath.Join(outputsDir, "amcache")
	os.MkdirAll(outPath, 0755)

	// Determine executable directory for relative path resolution
	exePath, err := os.Executable()
	var baseDir string
	if err == nil {
		baseDir = filepath.Dir(exePath)
	}

	// Dual Version Support: Search in net9 and net4 folders relative to executable
	parserPaths := []string{
		filepath.Join(baseDir, "AmcacheParsernet9", "AmcacheParser.exe"),
		filepath.Join(baseDir, "AmcacheParsernet4", "AmcacheParser.exe"),
		filepath.Join("AmcacheParsernet9", "AmcacheParser.exe"), // Fallback to current dir
		filepath.Join("AmcacheParsernet4", "AmcacheParser.exe"),
		"AmcacheParser.exe", // Fallback to PATH
	}

	var selectedParser string
	for _, p := range parserPaths {
		if _, err := os.Stat(p); err == nil {
			selectedParser = p
			break
		}
	}

	if selectedParser == "" {
		selectedParser = "AmcacheParser.exe" // Final attempt assuming it's in PATH
	}

	log.Printf("Using Amcache Parser: %s\n", selectedParser)
	cmd := exec.Command(selectedParser, "-f", amcacheHve, "--csv", outPath)
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to run Amcache Parser (ensure it is in AmcacheParsernet9 or AmcacheParsernet4): %v\n", err)
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
