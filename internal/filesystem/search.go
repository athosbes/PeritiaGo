package filesystem

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/athosbes/PeritiaGo/internal/hash"
	"github.com/athosbes/PeritiaGo/internal/models"
)

// SearchDrives recursively walks through given root paths looking for files
// that match either the target extensions or the explicit search term.
func SearchDrives(roots []string, targetExts []string, searchTerm string) []models.EvidenceFile {
	var evidence []models.EvidenceFile

	for _, root := range roots {
		log.Printf("Scanning root: %s ...\n", root)
		
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// We ignore access denied or missing files during traversal
				return nil
			}

			if d.IsDir() {
				// Skip very noisy or locked directories if needed
				if strings.Contains(strings.ToLower(path), "\\windows\\") {
					return filepath.SkipDir
				}
				return nil
			}

			matchedExt := MatchesExtension(d.Name(), targetExts)
			matchedSearch := MatchesSearch(path, searchTerm)

			if matchedExt || matchedSearch {
				info, err := d.Info()
				if err != nil {
					return nil
				}

				h, _ := hash.FileSHA256(path)
				
				// Create time on Windows requires syscalls, ModTime is acceptable as fallback 
				// or we can just use ModTime for both if direct syscalls are too verbose.
				// For the sake of simplicity, we'll use ModTime here as creation time in Go standard lib is obscured.
				ctime := info.ModTime() 

				evidence = append(evidence, models.EvidenceFile{
					Path:     path,
					Size:     info.Size(),
					Created:  ctime,
					Modified: info.ModTime(),
					SHA256:   h,
				})
			}
			return nil
		})

		if err != nil {
			log.Printf("Error walking %s: %v\n", root, err)
		}
	}

	log.Printf("Found %d evidence files\n", len(evidence))
	return evidence
}
