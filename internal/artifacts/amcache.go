package artifacts

import (
	"encoding/csv"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// ParseAmcache executes Eric Zimmerman's AmcacheParser.exe if available.
// It parses the Amcache.hve file for evidence of execution.
// ParseAmcache executes Eric Zimmerman's AmcacheParser.exe if available.
// It parses the Amcache.hve file for evidence of execution.
func ParseAmcache(outputsDir string) []models.Artifact {
	var artifacts []models.Artifact

	// We output its CSV to outputsDir/amcache/
	outPath, _ := filepath.Abs(filepath.Join(outputsDir, "amcache"))
	os.MkdirAll(outPath, 0755)

	// Step 1: Find HVE files
	hveFiles := findHVEFiles()
	if len(hveFiles) == 0 {
		log.Println("No Amcache.hve files found")
		return artifacts
	}

	// Step 2: Find Parsers
	parsers := findAmcacheParsers()
	if len(parsers) == 0 {
		log.Println("AmcacheParser.exe not found anywhere. Skipping detailed parse.")
		for _, hve := range hveFiles {
			artifacts = append(artifacts, models.Artifact{
				Name:        filepath.Base(hve),
				Type:        "Amcache",
				Path:        hve,
				Description: "Amcache.hve file found but parser missing",
				Timestamp:   time.Now().Format(time.RFC3339),
			})
		}
		return artifacts
	}

	// Step 3: Run Parser for each HVE file
	for _, hve := range hveFiles {
		var success bool
		var lastErr error
		var usedParser string

		for _, parser := range parsers {
			log.Printf("Trying Amcache Parser: %s\n", parser)
			cmd := exec.Command(parser, "-f", hve, "-i", "--csv", outPath, "--dt", "yyyy-MM-ddTHH:mm:ss")
			cmd.Dir = filepath.Dir(parser)

			if err := cmd.Run(); err == nil {
				success = true
				usedParser = parser
				break
			} else {
				lastErr = err
				log.Printf("[Warning] Parser %s failed: %v\n", parser, err)
			}
		}

		if success {
			// Ingest generated CSVs into the artifact slice
			csvFiles, _ := filepath.Glob(filepath.Join(outPath, "*.csv"))
			entriesAdded := 0
			for _, csvFile := range csvFiles {
				entriesAdded += parseAmcacheCSV(csvFile, &artifacts)
			}

			if entriesAdded == 0 {
				artifacts = append(artifacts, models.Artifact{
					Name:        filepath.Base(hve),
					Type:        "Amcache",
					Path:        hve,
					Description: "Amcache parsed successfully to " + outPath + " using " + filepath.Base(filepath.Dir(usedParser)) + ", but no row entries were loaded into HTML.",
					Timestamp:   time.Now().Format(time.RFC3339),
				})
			}
		} else {
			errMsg := "Unknown error"
			if lastErr != nil {
				errMsg = lastErr.Error()
			}
			artifacts = append(artifacts, models.Artifact{
				Name:        filepath.Base(hve),
				Type:        "Amcache",
				Path:        hve,
				Description: "All AmcacheParsers failed. Last error: " + errMsg,
				Timestamp:   time.Now().Format(time.RFC3339),
			})
		}
	}

	return artifacts
}

func parseAmcacheCSV(csvPath string, artifacts *[]models.Artifact) int {
	f, err := os.Open(csvPath)
	if err != nil {
		return 0
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		return 0
	}

	header := records[0]
	entries := 0

	for _, row := range records[1:] {
		desc := ""
		limit := len(row)
		if limit > 6 {
			limit = 6 // Only take the first 6 columns to avoid massive text bloat
		}
		for i := 0; i < limit; i++ {
			if i < len(header) {
				desc += header[i] + ": " + row[i] + " | "
			}
		}

		// Use the first column as timestamp if it looks vaguely like a timestamp
		ts := time.Now().Format(time.RFC3339)
		if len(row) > 0 && len(row[0]) >= 10 {
			ts = row[0]
		}

		*artifacts = append(*artifacts, models.Artifact{
			Name:        filepath.Base(csvPath),
			Type:        "AmcacheEntry",
			Path:        csvPath,
			Description: desc,
			Timestamp:   ts,
		})
		entries++
	}

	return entries
}

func findHVEFiles() []string {
	var files []string
	commonPaths := []string{
		`C:\Windows\AppCompat\Programs\Amcache.hve`,
		`C:\Windows\AppCompat\Programs\Amcache.hve.tmp`,
	}

	for _, p := range commonPaths {
		if _, err := os.Stat(p); err == nil {
			files = append(files, p)
		}
	}

	// Also look for other .hve files in the same directory just in case
	dir := `C:\Windows\AppCompat\Programs`
	matches, _ := filepath.Glob(filepath.Join(dir, "*.hve*"))
	for _, m := range matches {
		alreadyAdded := false
		for _, f := range files {
			if f == m {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			files = append(files, m)
		}
	}

	return files
}

func findAmcacheParsers() []string {
	exePath, err := os.Executable()
	var baseDir string
	if err == nil {
		baseDir = filepath.Dir(exePath)
	}

	parserPaths := []string{
		filepath.Join(baseDir, "AmcacheParsernet9", "AmcacheParser.exe"),
		filepath.Join(baseDir, "AmcacheParsernet4", "AmcacheParser.exe"),
		filepath.Join(baseDir, "AmcacheParser.exe"),
		filepath.Join("AmcacheParsernet9", "AmcacheParser.exe"),
		filepath.Join("AmcacheParsernet4", "AmcacheParser.exe"),
	}

	var validPaths []string
	seen := make(map[string]bool)
	for _, p := range parserPaths {
		if seen[p] {
			continue
		}
		seen[p] = true
		if _, err := os.Stat(p); err == nil {
			validPaths = append(validPaths, p)
		}
	}

	if _, err := exec.LookPath("AmcacheParser.exe"); err == nil {
		validPaths = append(validPaths, "AmcacheParser.exe")
	}

	return validPaths
}
