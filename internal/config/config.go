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
}

// ParseConfig parses command line flags and returns an AppConfig.
func ParseConfig() *AppConfig {
	caseName := flag.String("case", "Caso Padrao", "Nome do caso para o relatorio")
	investigator := flag.String("investigator", "Perito", "Nome do perito/investigador")
	extStr := flag.String("ext", "", "Extensoes para buscar (ex: exe,log,sqlite,db)")
	search := flag.String("search", "", "Termo de busca de software (ex: AnyDesk)")
	
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

	return &AppConfig{
		CaseName:     *caseName,
		Investigator: *investigator,
		Extensions:   extensions,
		SearchTerm:   *search,
	}
}
