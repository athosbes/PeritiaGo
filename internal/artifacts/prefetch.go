package artifacts

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// ParsePrefetch reads .pf files in C:\Windows\Prefetch to gather execution artifacts.
// Full native parsing of Windows 10 prefetch requires MAM/LZXpress decompression.
// This handles basic file discovery and timestamps as a fallback.
func ParsePrefetch() []models.Artifact {
	var artifacts []models.Artifact
	prefetchDir := `C:\Windows\Prefetch`

	files, err := os.ReadDir(prefetchDir)
	if err != nil {
		log.Printf("Cannot read Prefetch directory (requires admin): %v\n", err)
		return artifacts
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".pf") {
			fullPath := filepath.Join(prefetchDir, f.Name())
			info, err := f.Info()
			if err != nil {
				continue
			}

			// ModTime usually represents last execution time in uncompressed prefetch files
			// In Win10, the OS updates the file mod time upon execution updates.
			artifacts = append(artifacts, models.Artifact{
				Name:        f.Name(),
				Type:        "Prefetch",
				Path:        fullPath,
				Description: "Application execution trace",
				Timestamp:   info.ModTime().Format("2006-01-02 15:04:05"),
			})
		}
	}
	log.Printf("Collected %d Prefetch artifacts\n", len(artifacts))
	return artifacts
}
