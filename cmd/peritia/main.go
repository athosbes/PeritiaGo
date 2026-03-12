package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/athosbes/PeritiaGo/internal/artifacts"
	"github.com/athosbes/PeritiaGo/internal/capture"
	"github.com/athosbes/PeritiaGo/internal/config"
	"github.com/athosbes/PeritiaGo/internal/export"
	"github.com/athosbes/PeritiaGo/internal/filesystem"
	"github.com/athosbes/PeritiaGo/internal/hash"
	"github.com/athosbes/PeritiaGo/internal/models"
	"github.com/athosbes/PeritiaGo/internal/timeline"
)

func main() {
	log.Println("=== PeritiaGo Digital Forensics ===")
	cfg := config.ParseConfig()
	outDir := "outputs"
	os.MkdirAll(outDir, 0755)

	log.Println("[1] Capturing Installed Software via Registry")
	softwares := capture.GetInstalledSoftware()

	log.Println("[2] Opening Control Panel & Capturing Screenshot")
	screenshotPath, err := capture.OpenAppWizAndCapture(outDir)
	if err != nil {
		log.Printf("[Warning] Screenshot failed: %v\n", err)
	}

	log.Println("[3] Parsing Execution Artifacts")
	var arts []models.Artifact
	arts = append(arts, artifacts.ParsePrefetch()...)
	arts = append(arts, artifacts.ParseAmcache(outDir)...)
	arts = append(arts, artifacts.ParseShimCache()...)
	arts = append(arts, artifacts.ParseUserAssist()...)

	// Append search terms to check for residuals
	var searchTerms []string
	if cfg.SearchTerm != "" {
		searchTerms = append(searchTerms, cfg.SearchTerm)
	}
	arts = append(arts, artifacts.SearchResiduals(searchTerms)...)

	log.Println("[4] Filesystem Search")
	// Search in C:\Users by default for speed, but ideally could search entire drive
	drives := []string{filepath.Join(os.Getenv("SystemDrive")+"\\", "Users")}
	var evidences []models.EvidenceFile
	if len(cfg.Extensions) > 0 || cfg.SearchTerm != "" {
		evidences = filesystem.SearchDrives(drives, cfg.Extensions, cfg.SearchTerm)
	}

	log.Println("[5] Generating Forensic Timeline")
	tl := timeline.Generate(softwares, arts, evidences)

	log.Println("[6] Generating Final Report & Exports")
	date := time.Now()
	machine, _ := os.Hostname()
	
	report := models.FinalReport{
		CaseName:       cfg.CaseName,
		Investigator:   cfg.Investigator,
		MachineName:    machine,
		CaptureDate:    date,
		InstalledSofts: softwares,
		Artifacts:      arts,
		EvidenceFiles:  evidences,
		Timeline:       tl,
	}

	if screenshotPath != "" {
		report.Screenshots = append(report.Screenshots, screenshotPath)
	}

	// Export JSON/CSV
	var finalEvidences []models.Evidence
	csvOut, err := export.ToCSV(filepath.Join(outDir, "timeline.csv"), tl)
	if err == nil { finalEvidences = append(finalEvidences, csvOut) }
	
	jsonOut, err := export.ToJSON(filepath.Join(outDir, "report.json"), report)
	if err == nil { finalEvidences = append(finalEvidences, jsonOut) }

	if screenshotPath != "" {
		h, _ := hash.FileSHA256(screenshotPath)
		finalEvidences = append(finalEvidences, models.Evidence{
			FileName:  "programas_instalados.png",
			Path:      screenshotPath,
			Hash:      h,
			Timestamp: date,
		})
	}

	// Compute hashes for evidence files to write into manifest
	for i := range finalEvidences {
		if finalEvidences[i].Hash == "" {
			h, _ := hash.FileSHA256(finalEvidences[i].Path)
			finalEvidences[i].Hash = h
			finalEvidences[i].Timestamp = date
		}
	}
	report.Evidences = finalEvidences

	// Build manifesto
	log.Println("[7] Creating Master Manifest")
	manifestPath := filepath.Join(outDir, "manifesto.txt")
	var manifestLines []string
	for _, e := range report.Evidences {
		line := fmt.Sprintf("%s | %s", filepath.Base(e.Path), e.Hash)
		manifestLines = append(manifestLines, line)
	}
	
	manifestContent := strings.Join(manifestLines, "\n")
	os.WriteFile(manifestPath, []byte(manifestContent), 0644)
	
	masterHash := hash.StringSHA256(manifestContent)
	report.MasterHash = masterHash
	log.Printf("MASTER CHAIN OF CUSTODY HASH: %s\n", masterHash)

	// Finally, generate HTML
	htmlOut, err := export.ToHTML(filepath.Join(outDir, "report.html"), report)
	if err == nil {
		hHTML, _ := hash.FileSHA256(htmlOut.Path)
		htmlOut.Hash = hHTML
		// We append the HTML hash manually to manifest later if needed, but it's generated after master calculation.
	}

	log.Println("Forensic collection complete. Review 'outputs' directory.")
}