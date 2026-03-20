package ui

import (
	"log"
	"os/exec"
	"strings"

	"github.com/athosbes/PeritiaGo/internal/config"
)

// AskAllParameters presents a Windows Form to the user to capture all execution parameters.
func AskAllParameters(cfg *config.AppConfig) {
	script := `
[void] [System.Reflection.Assembly]::LoadWithPartialName("System.Windows.Forms")
[void] [System.Reflection.Assembly]::LoadWithPartialName("System.Drawing")

$objForm = New-Object System.Windows.Forms.Form 
$objForm.Text = "PeritiaGo - Configuração de Auditoria"
$objForm.Size = New-Object System.Drawing.Size(400,350) 
$objForm.StartPosition = "CenterScreen"
$objForm.KeyPreview = $True

# Case Name
$lblCase = New-Object System.Windows.Forms.Label
$lblCase.Location = New-Object System.Drawing.Point(10,20) 
$lblCase.Size = New-Object System.Drawing.Size(280,20) 
$lblCase.Text = "Nome do Caso:"
$objForm.Controls.Add($lblCase) 

$txtCase = New-Object System.Windows.Forms.TextBox 
$txtCase.Location = New-Object System.Drawing.Point(10,40) 
$txtCase.Size = New-Object System.Drawing.Size(360,20) 
$txtCase.Text = '` + cfg.CaseName + `'
$objForm.Controls.Add($txtCase) 

# Investigator
$lblInv = New-Object System.Windows.Forms.Label
$lblInv.Location = New-Object System.Drawing.Point(10,70) 
$lblInv.Size = New-Object System.Drawing.Size(280,20) 
$lblInv.Text = "Investigador/Perito:"
$objForm.Controls.Add($lblInv) 

$txtInv = New-Object System.Windows.Forms.TextBox 
$txtInv.Location = New-Object System.Drawing.Point(10,90) 
$txtInv.Size = New-Object System.Drawing.Size(360,20) 
$txtInv.Text = '` + cfg.Investigator + `'
$objForm.Controls.Add($txtInv)

# Extensions
$lblExt = New-Object System.Windows.Forms.Label
$lblExt.Location = New-Object System.Drawing.Point(10,120) 
$lblExt.Size = New-Object System.Drawing.Size(280,20) 
$lblExt.Text = "Extensões (vírgula):"
$objForm.Controls.Add($lblExt) 

$txtExt = New-Object System.Windows.Forms.TextBox 
$txtExt.Location = New-Object System.Drawing.Point(10,140) 
$txtExt.Size = New-Object System.Drawing.Size(360,20) 
$txtExt.Text = '` + strings.Join(cfg.Extensions, ",") + `'
$objForm.Controls.Add($txtExt)

# SearchTerm
$lblSearch = New-Object System.Windows.Forms.Label
$lblSearch.Location = New-Object System.Drawing.Point(10,170) 
$lblSearch.Size = New-Object System.Drawing.Size(280,20) 
$lblSearch.Text = "Termo de Busca (ex: AnyDesk):"
$objForm.Controls.Add($lblSearch) 

$txtSearch = New-Object System.Windows.Forms.TextBox 
$txtSearch.Location = New-Object System.Drawing.Point(10,190) 
$txtSearch.Size = New-Object System.Drawing.Size(360,20) 
$txtSearch.Text = '` + cfg.SearchTerm + `'
$objForm.Controls.Add($txtSearch)

# Drives
$lblDrives = New-Object System.Windows.Forms.Label
$lblDrives.Location = New-Object System.Drawing.Point(10,220) 
$lblDrives.Size = New-Object System.Drawing.Size(280,20) 
$lblDrives.Text = "Drives/Pastas (vírgula):"
$objForm.Controls.Add($lblDrives) 

$txtDrives = New-Object System.Windows.Forms.TextBox 
$txtDrives.Location = New-Object System.Drawing.Point(10,240) 
$txtDrives.Size = New-Object System.Drawing.Size(360,20) 
$txtDrives.Text = '` + strings.Join(cfg.Drives, ",") + `'
$objForm.Controls.Add($txtDrives)

$OKButton = New-Object System.Windows.Forms.Button
$OKButton.Location = New-Object System.Drawing.Point(150,280)
$OKButton.Size = New-Object System.Drawing.Size(75,23)
$OKButton.Text = "OK"
$OKButton.DialogResult = [System.Windows.Forms.DialogResult]::OK
$objForm.AcceptButton = $OKButton
$objForm.Controls.Add($OKButton)

$objForm.Topmost = $True

$result = $objForm.ShowDialog()

if ($result -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output ($txtCase.Text + "|" + $txtInv.Text + "|" + $txtExt.Text + "|" + $txtSearch.Text + "|" + $txtDrives.Text)
}
`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	outputBytes, err := cmd.Output()
	if err != nil {
		log.Printf("[Warning] Failed to show config UI: %v", err)
		return
	}

	outStr := strings.TrimSpace(string(outputBytes))
	if outStr == "" {
		return // user cancelled or closed
	}

	parts := strings.Split(outStr, "|")
	if len(parts) == 5 {
		cfg.CaseName = strings.TrimSpace(parts[0])
		cfg.Investigator = strings.TrimSpace(parts[1])

		exts := strings.Split(parts[2], ",")
		cfg.Extensions = []string{} // reset
		for _, e := range exts {
			e = strings.TrimSpace(e)
			if e != "" {
				if !strings.HasPrefix(e, ".") {
					e = "." + e
				}
				cfg.Extensions = append(cfg.Extensions, e)
			}
		}

		cfg.SearchTerm = strings.TrimSpace(parts[3])

		dvs := strings.Split(parts[4], ",")
		cfg.Drives = []string{} // reset
		for _, d := range dvs {
			d = strings.TrimSpace(d)
			if d != "" {
				cfg.Drives = append(cfg.Drives, d)
			}
		}
	}
}
