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
	Architecture    string `json:"architecture" csv:"Arquitetura"` // 32 or 64 bits
	MSIGUID         string `json:"msi_guid" csv:"GUID_MSI"`        // MSI package GUID
	ProductID       string `json:"product_id" csv:"ID_Produto"`    // Product ID if exists
	Source          string `json:"source" csv:"Fonte_Evidencia"`   // Registry, Amcache, Prefetch, etc.
}

// MachineIdentity carries all identification data for the target machine.
type MachineIdentity struct {
	Hostname     string            `json:"hostname"`
	CurrentUser  string            `json:"current_user"`
	Domain       string            `json:"domain"`
	OSName       string            `json:"os_name"`
	OSVersion    string            `json:"os_version"`
	OSBuild      string            `json:"os_build"`
	Manufacturer string            `json:"manufacturer"`
	Model        string            `json:"model"`
	SerialNumber string            `json:"serial_number"`
	BIOSUUID     string            `json:"bios_uuid"`
	MachineGUID  string            `json:"machine_guid"`
	IPAddresses  []string          `json:"ip_addresses"`
	MACAddresses map[string]string `json:"mac_addresses"` // Interface Name -> MAC
}

// ForensicMetadata contains information about the collection process.
type ForensicMetadata struct {
	CollectionDate time.Time `json:"collection_date"`
	ToolName       string    `json:"tool_name"`
	ToolVersion    string    `json:"tool_version"`
	Executor       string    `json:"executor"`
	MachineName    string    `json:"machine_name"`
	ReportHash     string    `json:"report_hash"`
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
	Source      string    `json:"source" csv:"Fonte"`          // e.g., "File Created", "Prefetch Execution", "Amcache"
	Description string    `json:"description" csv:"Descricao"` // Contextual details about the event
}

// EvidenceFile represents a file discovered during the file search phase (e.g., by extension or name).
type EvidenceFile struct {
	Path        string    `json:"path" csv:"Caminho"`
	Name        string    `json:"name" csv:"Nome"`
	Size        int64     `json:"size" csv:"Tamanho"`
	Created     time.Time `json:"created" csv:"Criacao"`
	Modified    time.Time `json:"modified" csv:"Modificacao"`
	SHA256      string    `json:"sha256" csv:"Hash_SHA256"`
	FileVersion string    `json:"file_version" csv:"Versao_Arquivo"`
	CompanyName string    `json:"company_name" csv:"Empresa"`
	ProductName string    `json:"product_name" csv:"Produto"`
}

// Evidence represents a generated file (image or csv) for chain of custody.
type Evidence struct {
	FileName  string    `json:"file_name"`
	Path      string    `json:"path"`
	Hash      string    `json:"sha256_hash"`
	Timestamp time.Time `json:"timestamp"`
}

// LicenseData captures license information.
type LicenseData struct {
	SoftwareName string `json:"software_name"`
	ProductKey   string `json:"product_key"`
	LicenseType  string `json:"license_type"` // OEM / Volume / Retail / Trial
	Status       string `json:"status"`
	LicenseID    string `json:"license_id"`
	FilePath     string `json:"file_path"`
}

// FinalReport encapsulates all gathered evidence for export to JSON, CSV, or HTML.
type FinalReport struct {
	CaseName       string           `json:"case_name"`
	Investigator   string           `json:"investigator"`
	Machine        MachineIdentity  `json:"machine_identity"`
	Metadata       ForensicMetadata `json:"metadata"`
	CaptureDate    time.Time        `json:"capture_date"`
	InstalledSofts []Software       `json:"installed_software"`
	Artifacts      []Artifact       `json:"artifacts"`
	Timeline       []TimelineEvent  `json:"timeline"`
	EvidenceFiles  []EvidenceFile   `json:"evidence_files"`
	Licenses       []LicenseData    `json:"licenses"`
	Evidences      []Evidence       `json:"evidences"`
	Screenshots    []string         `json:"screenshots"`
	MasterManifest []string         `json:"-"`
	MasterHash     string           `json:"master_hash"`
}
