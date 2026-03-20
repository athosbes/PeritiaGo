package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// AppConfig holds the configuration parameters for the execution.
type AppConfig struct {
	CaseName     string   `json:"case_name"`
	Investigator string   `json:"investigator"`
	Extensions   []string `json:"extensions"`
	SearchTerm   string   `json:"search_term"`
	Drives       []string `json:"drives"`
}

// ParseConfig parses command line flags and returns an AppConfig.
// It prioritizes command-line flags, then config file defaults, then hardcoded defaults.
func ParseConfig() *AppConfig {
	// Look for peritiago_config.json beside the executable
	exePath, _ := os.Executable()
	configPath := filepath.Join(filepath.Dir(exePath), "peritiago_config.json")

	defaultConfig := AppConfig{
		CaseName:     "Caso Padrao",
		Investigator: "Perito",
		Extensions:   []string{},
		SearchTerm:   "",
		Drives:       []string{"C:\\Users"},
	}

	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &defaultConfig); err != nil {
			log.Printf("[Warning] Failed to parse %s: %v\n", configPath, err)
		} else {
			log.Printf("[Info] Loaded configuration from %s\n", configPath)
		}
	} else if os.IsNotExist(err) {
		// Create default template for the user to edit later
		if outData, err := json.MarshalIndent(defaultConfig, "", "  "); err == nil {
			os.WriteFile(configPath, outData, 0644)
			log.Printf("[Info] Created default configuration template at: %s\n", configPath)
		}
	}

	extDefault := strings.Join(defaultConfig.Extensions, ",")
	drivesDefault := strings.Join(defaultConfig.Drives, ",")

	caseName := flag.String("case", defaultConfig.CaseName, "Nome do caso para o relatorio")
	investigator := flag.String("investigator", defaultConfig.Investigator, "Nome do perito/investigador")
	extStr := flag.String("ext", extDefault, "Extensoes para buscar (ex: exe,log,sqlite,db)")
	search := flag.String("search", defaultConfig.SearchTerm, "Termo de busca de software (ex: AnyDesk)")
	drivesStr := flag.String("drives", drivesDefault, "Diretorios ou drives alvo separados por virgula (ex: C:\\,D:\\)")

	flag.Parse()

	var extensions []string
	if *extStr != "" {
		parts := strings.Split(*extStr, ",")
		for _, p := range parts {
			ext := strings.TrimSpace(p)
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			extensions = append(extensions, ext)
		}
	}

	var drives []string
	if *drivesStr != "" {
		parts := strings.Split(*drivesStr, ",")
		for _, p := range parts {
			drives = append(drives, strings.TrimSpace(p))
		}
	}

	return &AppConfig{
		CaseName:     *caseName,
		Investigator: *investigator,
		Extensions:   extensions,
		SearchTerm:   *search,
		Drives:       drives,
	}
}
