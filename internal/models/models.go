package models

import "time"

// Software represents an installed or historically executed software on the system.
type Software struct {
	DisplayName     string `json:"display_name" csv:"Nome"`
	DisplayVersion  string `json:"display_version" csv:"Versao"`
	Publisher       string `json:"publisher" csv:"Editor"`
	InstallDate     string `json:"install_date" csv:"Data_Instalacao"`
	InstallLocation string `json:"install_location" csv:"Local_Instalacao"`
	UninstallString string `json:"uninstall_string" csv:"Comando_Desinstalacao"`
	Source          string `json:"source" csv:"Fonte_Evidencia"` // Registry, Amcache, Prefetch, etc.
}

// Artifact represents a generic forensic artifact (e.g. registry key, file, log entry).
type Artifact struct {
	Name        string `json:"name" csv:"Nome"`
	Type        string `json:"type" csv:"Tipo"` // e.g., "Prefetch", "Amcache", "ShimCache", "UserAssist", "Registry"
	Path        string `json:"path" csv:"Caminho"`
	Description string `json:"description" csv:"Descricao"`
	Value       string `json:"value" csv:"Valor"`
	Timestamp   string `json:"timestamp" csv:"Data_Hora"`
}

// TimelineEvent represents a single chronological event in the forensic timeline.
type TimelineEvent struct {
	Timestamp   time.Time `json:"timestamp" csv:"Data_Hora"`
	Event       string    `json:"event" csv:"Evento"`
	Source      string    `json:"source" csv:"Fonte"`           // e.g., "File Created", "Prefetch Execution", "Amcache"
	Description string    `json:"description" csv:"Descricao"`  // Contextual details about the event
}

// EvidenceFile represents a file discovered during the file search phase (e.g., by extension or name).
type EvidenceFile struct {
	Path     string    `json:"path" csv:"Caminho"`
	Size     int64     `json:"size" csv:"Tamanho"`
	Created  time.Time `json:"created" csv:"Criacao"`
	Modified time.Time `json:"modified" csv:"Modificacao"`
	SHA256   string    `json:"sha256" csv:"Hash_SHA256"`
}

// Evidence represents a generated file (image or csv) for chain of custody.
type Evidence struct {
	FileName  string    `json:"file_name"`
	Path      string    `json:"path"`
	Hash      string    `json:"sha256_hash"`
	Timestamp time.Time `json:"timestamp"`
}

// FinalReport encapsulates all gathered evidence for export to JSON, CSV, or HTML.
type FinalReport struct {
	CaseName         string          `json:"case_name"`
	Investigator     string          `json:"investigator"`
	MachineName      string          `json:"machine_name"`
	CaptureDate      time.Time       `json:"capture_date"`
	InstalledSofts   []Software      `json:"installed_software"`
	Artifacts        []Artifact      `json:"artifacts"`
	Timeline         []TimelineEvent `json:"timeline"`
	EvidenceFiles    []EvidenceFile  `json:"evidence_files"`
	Evidences        []Evidence      `json:"evidences"`
	Screenshots      []string        `json:"screenshots"`
	MasterManifest   []string        `json:"-"` // Not marshalled directly, just holds paths/hashes for manifest
	MasterHash       string          `json:"master_hash"`
}