package filesystem

import (
	"fmt"
	"io/fs"

	"log"
	"os/exec"

	"path/filepath"
	"strings"
	"syscall"
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
				return nil
			}

			if d.IsDir() {
				lowPath := strings.ToLower(path)

				// Standard exclusions to prevent hanging and infinite recursion
				excludedDirs := []string{
					"\\windows",
					"\\$recycle.bin",
					"\\system volume information",
				}

				for _, excluded := range excludedDirs {
					if strings.Contains(lowPath, excluded) {
						return filepath.SkipDir
					}
				}
				return nil
			}

			matchedExt := MatchesExtension(d.Name(), targetExts)
			matchedSearch := MatchesSearch(path, searchTerm)

			// We only process files that match our criteria
			if matchedExt || matchedSearch {
				info, err := d.Info()
				if err != nil {
					return nil
				}

				h, _ := hash.FileSHA256(path)

				// Populate EvidenceFile with all requested fields
				ef := models.EvidenceFile{
					Path:     path,
					Name:     d.Name(),
					Size:     info.Size(),
					Created:  getWinCreationTime(info),
					Modified: info.ModTime(),
					SHA256:   h,
				}

				// If it's an executable, get more metadata
				if strings.HasSuffix(strings.ToLower(path), ".exe") || strings.HasSuffix(strings.ToLower(path), ".dll") {
					// CRITICAL PERFORMANCE FIX: Extracting metadata via powershell takes ~0.5s per file.
					// We only do this if it explicitly matched a search term, otherwise a broad ".exe" search on C:\ will hang for hours.
					if matchedSearch && searchTerm != "" {
						ef.FileVersion, ef.CompanyName, ef.ProductName = getFileVersionInfo(path)
					}
				}

				evidence = append(evidence, ef)
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

func getWinCreationTime(info fs.FileInfo) time.Time {
	if winAttr, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		return time.Unix(0, winAttr.CreationTime.Nanoseconds())
	}
	return info.ModTime()
}

func getFileVersionInfo(path string) (version, company, product string) {
	// Using PowerShell as a reliable way to get VersionInfo on Windows without 3rd party libs
	psCmd := fmt.Sprintf(`(Get-Item -Path "%s").VersionInfo | Select-Object FileVersion, CompanyName, ProductName | ConvertTo-Csv -NoTypeInformation`, path)
	out, err := exec.Command("powershell", "-Command", psCmd).Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) >= 2 {
			fields := strings.Split(lines[1], ",")
			if len(fields) >= 3 {
				version = strings.Trim(fields[0], "\" \r\n")
				company = strings.Trim(fields[1], "\" \r\n")
				product = strings.Trim(fields[2], "\" \r\n")
			}
		}
	}
	return
}
