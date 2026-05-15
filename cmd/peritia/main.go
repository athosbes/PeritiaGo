package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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
	"github.com/athosbes/PeritiaGo/internal/logger"
	"github.com/athosbes/PeritiaGo/internal/models"
	"github.com/athosbes/PeritiaGo/internal/timeline"
	"github.com/athosbes/PeritiaGo/internal/ui"
	"github.com/athosbes/PeritiaGo/internal/verify"
	"github.com/google/uuid"
)

// Version is set at build time via -ldflags "-X main.Version=vX.Y.Z"
var Version = "1.2.0-dev"

func CalculateSelfHash() (string, string, int64, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", "", 0, err
	}

	file, err := os.Open(exePath)
	if err != nil {
		return exePath, "", 0, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return exePath, "", 0, err
	}

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return exePath, "", info.Size(), err
	}

	return exePath, hex.EncodeToString(h.Sum(nil)), info.Size(), nil
}

func main() {
	fmt.Println("=== PeritiaGo Digital Forensics ===")
	config.CurrentConfig = config.ParseConfig()
	cfg := config.CurrentConfig

	// Visual GUI for all parameters if Investigator is default
	if cfg.Investigator == "Perito" {
		ui.AskAllParameters(cfg)
	}

	// Dynamic Output Directory following pattern
	fullIdentity := identity.GetFullIdentity()
	timestamp := time.Now().Format("20060102_150405")
	outDir := fmt.Sprintf("software_inventory_%s_%s_%s", fullIdentity.MachineGUID, identity.GetMACAddress(), timestamp)

	// Normalize and store in config
	absOutDir, err := filepath.Abs(outDir)
	if err != nil {
		log.Fatalf("Failed to resolve output directory: %v", err)
	}
	cfg.OutputDir = absOutDir
	cfg.ExecutionID = uuid.New().String()

	if err := filesystem.MkdirAll(cfg.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Initialize Logger
	if err := logger.Init(cfg.OutputDir); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Self-Hash
	exePath, selfHash, exeSize, err := CalculateSelfHash()
	if err != nil {
		logger.Log.Error("failed to calculate self-hash", "error", err)
		os.Exit(1)
	}

	hostname, _ := os.Hostname()

	// Initial Forensic Log
	logger.Log.Info("execution started",
		"version", Version,
		"forensic_mode", cfg.ForensicMode,
		"execution_id", cfg.ExecutionID,
		"user", os.Getenv("USERNAME"),
		"hostname", hostname,
		"pid", os.Getpid(),
		"output_dir", cfg.OutputDir,
	)

	logger.Log.Info("binary info",
		"path", exePath,
		"sha256", selfHash,
		"size_bytes", exeSize,
	)

	// Self-Integrity Check: Check if the running binary is signed
	// Tries Sigstore first (.sigstore.json bundle), falls back to Windows Authenticode
	if err := verify.IsSigned(exePath); err != nil {
		logger.Log.Warn("binary integrity warning", "status", "unsigned or invalid signature", "error", err)
	} else {
		// Determine which method succeeded for reporting
		sigstoreBundle := exePath + ".sigstore.json"
		if _, statErr := os.Stat(sigstoreBundle); statErr == nil {
			logger.Log.Info("binary integrity verified", "status", "sigstore signature valid", "bundle", sigstoreBundle)
		} else {
			logger.Log.Info("binary integrity verified", "status", "windows authenticode signature valid")
		}
	}

	logger.Log.Info("starting collection", "step", 1, "description", "Capturing Installed Software")
	softwares := capture.GetInstalledSoftware()

	// Capture Installed Software via WMI
	wmiPath, err := capture.CaptureInstalledSoftwareWMI(cfg.OutputDir)
	if err == nil {
		logger.Log.Info("Installed software (WMI) saved", "path", wmiPath)
	}

	// Capture Winget
	capture.CaptureWinget(cfg.OutputDir)

	logger.Log.Info("opening control panels & capturing screenshots", "step", 2)
	screenshotPath, err := capture.OpenAppWizAndCapture(cfg.OutputDir)
	if err != nil {
		logger.Log.Warn("programs screenshot failed", "error", err)
	}

	systemInfoPath, err := capture.OpenSystemInfoAndCapture(cfg.OutputDir)
	if err != nil {
		logger.Log.Warn("system info screenshot failed", "error", err)
	}

	logger.Log.Info("parsing execution artifacts & event logs", "step", 3)
	var arts []models.Artifact
	arts = append(arts, artifacts.ParsePrefetch()...)
	arts = append(arts, artifacts.ParseAmcache(cfg.OutputDir)...)
	arts = append(arts, artifacts.ParseShimCache()...)
	arts = append(arts, artifacts.ParseUserAssist()...)
	arts = append(arts, capture.GetSystemStatus()...)
	arts = append(arts, capture.GetRunningProcesses()...)

	// Capture Event Logs for installs/uninstalls
	softEvents := capture.GetSoftwareEvents()

	// Capture residual traces
	arts = append(arts, capture.SearchResidualTraces()...)

	// Capture License Data
	licenses := capture.GetLicenseData()

	// Append search terms to check for residuals
	var searchTerms []string
	if cfg.SearchTerm != "" {
		searchTerms = append(searchTerms, cfg.SearchTerm)
	}
	arts = append(arts, capture.SearchResidualsByTerms(searchTerms)...)

	logger.Log.Info("filesystem search", "step", 4)
	var evidences []models.EvidenceFile
	if len(cfg.Extensions) > 0 || cfg.SearchTerm != "" {
		evidences = filesystem.SearchDrives(cfg.Drives, cfg.Extensions, cfg.SearchTerm)
	}

	logger.Log.Info("generating forensic timeline", "step", 5)
	tl := timeline.Generate(softwares, arts, evidences)
	// Add software events to timeline
	for _, se := range softEvents {
		tl = append(tl, models.TimelineEvent{
			Timestamp:   se.Timestamp,
			Event:       se.Event,
			Source:      se.Source,
			Description: se.Description,
		})
	}

	logger.Log.Info("generating final report & exports", "step", 6)
	date := time.Now()

	report := models.FinalReport{
		CaseName:     cfg.CaseName,
		Investigator: cfg.Investigator,
		Machine:      fullIdentity,
		Metadata: models.ForensicMetadata{
			CollectionDate: date,
			ToolName:       "PeritiaGo",
			ToolVersion:    Version,
			Executor:       cfg.Investigator,
			MachineName:    fullIdentity.Hostname,
			ForensicMode:   cfg.ForensicMode,
			WriteProtection: models.WriteProtectionInfo{
				Enabled:           cfg.ForensicMode,
				SymlinkProtection: true, // Always active if guard is used
			},
		},
		CaptureDate:    date,
		InstalledSofts: softwares,
		Artifacts:      arts,
		EvidenceFiles:  evidences,
		Timeline:       tl,
		Licenses:       licenses,
	}

	if screenshotPath != "" {
		report.Screenshots = append(report.Screenshots, screenshotPath)
	}
	if systemInfoPath != "" {
		report.Screenshots = append(report.Screenshots, systemInfoPath)
	}

	// Export CSV & HTML first (to be included in the manifest)
	var finalEvidences []models.Evidence
	csvOut, err := export.ToCSV(filepath.Join(cfg.OutputDir, "timeline.csv"), tl)
	if err == nil {
		finalEvidences = append(finalEvidences, csvOut)
	}

	// Scan the entire output directory recursively to find ALL generated files (CSVs from AMCache, etc.)
	filepath.WalkDir(cfg.OutputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(cfg.OutputDir, path)
		// Skip files that don't exist yet or shouldn't be in the manifest before its creation
		if relPath == "manifesto.txt" || relPath == "report.json" {
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

	// Finalize evidences list
	report.Evidences = finalEvidences

	// Build manifesto
	logger.Log.Info("creating master manifest", "step", 7)
	manifestPath := filepath.Join(cfg.OutputDir, "manifesto.txt")
	var manifestLines []string

	seen := make(map[string]bool)
	for _, e := range report.Evidences {
		if seen[e.Path] {
			continue
		}
		seen[e.Path] = true
		line := fmt.Sprintf("%s | %s", filepath.Base(e.Path), e.Hash)
		manifestLines = append(manifestLines, line)
	}

	manifestContent := strings.Join(manifestLines, "\n")
	filesystem.WriteFile(manifestPath, []byte(manifestContent), 0644)

	// Calculate Master Hash
	masterHash := hash.StringSHA256(manifestContent)
	report.MasterHash = masterHash
	logger.Log.Info("master chain of custody hash", "hash", masterHash)

	// Export JSON/HTML NOW with the master hash filled
	jsonOut, err := export.ToJSON(filepath.Join(cfg.OutputDir, "report.json"), report)
	if err == nil {
		report.Evidences = append(report.Evidences, jsonOut)
	}

	htmlOut, err := export.ToHTML(filepath.Join(cfg.OutputDir, "report.html"), report)
	if err == nil {
		report.Evidences = append(report.Evidences, htmlOut)
	}

	logger.Log.Info("forensic collection complete", "output_dir", cfg.OutputDir)
}
