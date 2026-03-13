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
	"github.com/athosbes/PeritiaGo/internal/identity"
	"github.com/athosbes/PeritiaGo/internal/models"
	"github.com/athosbes/PeritiaGo/internal/timeline"
	"github.com/athosbes/PeritiaGo/internal/ui"
)

func main() {
	log.Println("=== PeritiaGo Digital Forensics ===")
	cfg := config.ParseConfig()

	// Visual GUI for Investigator and Extensions if not provided via flags
	if cfg.Investigator == "Perito" {
		cfg.Investigator = ui.AskInvestigator()
	}
	if len(cfg.Extensions) == 0 {
		extStr := ui.AskExtensions()
		if extStr != "" {
			parts := strings.Split(extStr, ",")
			for _, p := range parts {
				ext := strings.TrimSpace(p)
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext
				}
				cfg.Extensions = append(cfg.Extensions, ext)
			}
		}
	}

	// Dynamic Output Directory following pattern
	machineUUID := identity.GetMachineUUID()
	macAddr := identity.GetMACAddress()
	timestamp := time.Now().Format("20060102_150405")
	outDir := fmt.Sprintf("software_inventory_%s_%s_%s", machineUUID, macAddr, timestamp)
	
	os.MkdirAll(outDir, 0755)
	log.Printf("Output directory: %s\n", outDir)

	log.Println("[1] Capturing Installed Software via Registry, WMIC & Winget")
	softwares := capture.GetInstalledSoftware()
	
	// Add WMIC and Winget captures
	capture.CaptureWMIC(outDir)
	capture.CaptureWinget(outDir)

	log.Println("[2] Opening Control Panels & Capturing Screenshots")
	screenshotPath, err := capture.OpenAppWizAndCapture(outDir)
	if err != nil {
		log.Printf("[Warning] Programs screenshot failed: %v\n", err)
	}
	
	systemInfoPath, err := capture.OpenSystemInfoAndCapture(outDir)
	if err != nil {
		log.Printf("[Warning] System info screenshot failed: %v\n", err)
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
	var evidences []models.EvidenceFile
	if len(cfg.Extensions) > 0 || cfg.SearchTerm != "" {
		evidences = filesystem.SearchDrives(cfg.Drives, cfg.Extensions, cfg.SearchTerm)
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

	if screenshotPath != "" { report.Screenshots = append(report.Screenshots, screenshotPath) }
	if systemInfoPath != "" { report.Screenshots = append(report.Screenshots, systemInfoPath) }

	// Export JSON/CSV/HTML
	var finalEvidences []models.Evidence
	
	csvOut, err := export.ToCSV(filepath.Join(outDir, "timeline.csv"), tl)
	if err == nil { finalEvidences = append(finalEvidences, csvOut) }
	
	jsonOut, err := export.ToJSON(filepath.Join(outDir, "report.json"), report)
	if err == nil { finalEvidences = append(finalEvidences, jsonOut) }

	// ROBUST EVIDENCE TRACKING: 
	// Scan the entire output directory recursively to find ALL generated files (CSVs from AMCache, etc.)
	filepath.WalkDir(outDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() { return nil }
		
		// Skip files we already handled explicitly if needed, but safer to just hash everything found
		// and avoid duplicates in the manifest later.
		relPath, _ := filepath.Rel(outDir, path)
		if relPath == "manifesto.txt" || relPath == "report.json" || relPath == "timeline.csv" {
			return nil
		}

		h, _ := hash.FileSHA256(path)
		finalEvidences = append(finalEvidences, models.Evidence{
			FileName:  filepath.Base(path),
			Path:      path,
			Hash:      h,
			Timestamp: date,
		})
		return nil
	})

	// Compute missing hashes and finalize report
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
	
	// Deduplicate and build manifest
	seen := make(map[string]bool)
	for _, e := range report.Evidences {
		if seen[e.Path] { continue }
		seen[e.Path] = true
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
	}

	log.Printf("Forensic collection complete. Review '%s' directory.\n", outDir)
}