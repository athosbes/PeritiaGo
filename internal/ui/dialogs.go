package ui

import (
	"os/exec"
	"strings"
)

// AskInvestigator uses a PowerShell InputBox to ask for the investigator's name.
func AskInvestigator() string {
	script := `
[System.Reflection.Assembly]::LoadWithPartialName('Microsoft.VisualBasic') | Out-Null
[Microsoft.VisualBasic.Interaction]::InputBox("Nome do Perito/Investigador:", "PeritiaGo - Entrada", "Perito")
`
	cmd := exec.Command("powershell", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return "Perito"
	}
	return strings.TrimSpace(string(output))
}

// AskExtensions uses a PowerShell InputBox to ask for file extensions to search.
func AskExtensions() string {
	script := `
[System.Reflection.Assembly]::LoadWithPartialName('Microsoft.VisualBasic') | Out-Null
[Microsoft.VisualBasic.Interaction]::InputBox("Extensões para buscar (separadas por vírgula, ex: exe,log,sqlite):", "PeritiaGo - Entrada", "exe,log,sqlite")
`
	cmd := exec.Command("powershell", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
