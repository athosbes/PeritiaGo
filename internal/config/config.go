package config

import (
	"flag"
	"strings"
)

// AppConfig holds the configuration parameters for the execution.
type AppConfig struct {
	CaseName     string
	Investigator string
	Extensions   []string
	SearchTerm   string
	Drives       []string
}

// ParseConfig parses command line flags and returns an AppConfig.
func ParseConfig() *AppConfig {
	caseName := flag.String("case", "Caso Padrao", "Nome do caso para o relatorio")
	investigator := flag.String("investigator", "Perito", "Nome do perito/investigador")
	extStr := flag.String("ext", "", "Extensoes para buscar (ex: exe,log,sqlite,db)")
	search := flag.String("search", "", "Termo de busca de software (ex: AnyDesk)")
	drivesStr := flag.String("drives", "C:\\Users", "Diretorios ou drives alvo separados por virgula (ex: C:\\,D:\\)")
	
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
