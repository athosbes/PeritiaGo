package filesystem

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/athosbes/PeritiaGo/internal/executor"
	"github.com/athosbes/PeritiaGo/internal/hash"
	"github.com/athosbes/PeritiaGo/internal/models"
)

// SearchDrives recursively walks through given root paths looking for files
// that match either the target extensions or the explicit search term.
// It uses a worker pool for concurrent processing of matching files.
func SearchDrives(roots []string, targetExts []string, searchTerm string) []models.EvidenceFile {
	var evidence []models.EvidenceFile
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Use a bounded worker pool for processing matched files (hashing, metadata)
	// This prevents I/O exhaustion and unbounded goroutine growth.
	numWorkers := runtime.NumCPU() * 2
	jobs := make(chan string, 1000)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				ef, err := processFile(path, targetExts, searchTerm)
				if err == nil {
					mu.Lock()
					evidence = append(evidence, ef)
					mu.Unlock()
				}
			}
		}()
	}

	for _, root := range roots {
		log.Printf("Scanning root: %s ...\n", root)

		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.IsDir() {
				lowPath := strings.ToLower(path)
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

			// Fast check before sending to worker
			if MatchesExtension(d.Name(), targetExts) || MatchesSearch(path, searchTerm) {
				jobs <- path
			}
			return nil
		})

		if err != nil {
			log.Printf("Error walking %s: %v\n", root, err)
		}
	}

	close(jobs)
	wg.Wait()

	log.Printf("Found %d evidence files\n", len(evidence))
	return evidence
}

func processFile(path string, targetExts []string, searchTerm string) (models.EvidenceFile, error) {
	info, err := os.Stat(path)
	if err != nil {
		return models.EvidenceFile{}, err
	}

	h, _ := hash.FileSHA256(path)

	ef := models.EvidenceFile{
		Path:     path,
		Name:     filepath.Base(path),
		Size:     info.Size(),
		Created:  getWinCreationTime(info),
		Modified: info.ModTime(),
		SHA256:   h,
	}

	matchedSearch := MatchesSearch(path, searchTerm)
	if strings.HasSuffix(strings.ToLower(path), ".exe") || strings.HasSuffix(strings.ToLower(path), ".dll") {
		if matchedSearch && searchTerm != "" {
			ef.FileVersion, ef.CompanyName, ef.ProductName = getFileVersionInfo(path)
		}
	}

	return ef, nil
}

func getWinCreationTime(info fs.FileInfo) time.Time {
	if winAttr, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		return time.Unix(0, winAttr.CreationTime.Nanoseconds())
	}
	return info.ModTime()
}

func getFileVersionInfo(path string) (version, company, product string) {
	// Using PowerShell via secure executor to get VersionInfo on Windows without 3rd party libs
	// We pass the path as a single argument to avoid command injection
	psCmd := `(Get-Item -LiteralPath $args[0]).VersionInfo | Select-Object FileVersion, CompanyName, ProductName | ConvertTo-Csv -NoTypeInformation`
	out, err := executor.Execute("powershell", "-NoProfile", "-Command", psCmd, path)
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
