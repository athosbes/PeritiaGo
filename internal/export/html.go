package export

import (
	"html/template"
	"os"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// ToHTML builds a comprehensive HTML report from the FinalReport.
func ToHTML(path string, report models.FinalReport) (models.Evidence, error) {
	file, err := os.Create(path)
	if err != nil {
		return models.Evidence{}, err
	}
	defer file.Close()

	const tpl = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>PeritiaGo - Forensic Report</title>
    <style>
        body { font-family: 'Segoe UI', Arial, sans-serif; background: #f4f7f6; color: #333; margin: 40px; }
        h1, h2, h3 { color: #2c3e50; }
        .header { background: #2980b9; color: white; padding: 20px; border-radius: 5px; }
        .section { background: white; margin-top: 20px; padding: 20px; border-radius: 5px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
        table { width: 100%; border-collapse: collapse; margin-top: 15px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #34495e; color: white; }
        pre { background: #eee; padding: 10px; overflow-x: auto; }
        img { max-width: 100%; border: 1px solid #ddd; }
    </style>
</head>
<body>
    <div class="header">
        <h1>PeritiaGo - Forensic Report</h1>
        <p><strong>Case:</strong> {{.CaseName}}</p>
        <p><strong>Investigator:</strong> {{.Investigator}}</p>
        <p><strong>Date:</strong> {{.CaptureDate}}</p>
    </div>

    <div class="section">
        <h2>1. Softwares Detectados</h2>
        <table>
            <tr><th>Name</th><th>Version</th><th>Publisher</th><th>Install Date</th><th>Source</th></tr>
            {{range .InstalledSofts}}
            <tr>
                <td>{{.DisplayName}}</td>
                <td>{{.DisplayVersion}}</td>
                <td>{{.Publisher}}</td>
                <td>{{.InstallDate}}</td>
                <td>{{.Source}}</td>
            </tr>
            {{end}}
        </table>
    </div>

    <div class="section">
        <h2>2. Artefatos de Execução (Amcache, Prefetch, ShimCache, UserAssist)</h2>
        <table>
            <tr><th>Name</th><th>Type</th><th>Description</th><th>Timestamp</th></tr>
            {{range .Artifacts}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.Type}}</td>
                <td>{{.Description}}</td>
                <td>{{.Timestamp}}</td>
            </tr>
            {{end}}
        </table>
    </div>

    <div class="section">
        <h2>3. Busca por Extensão / Residuais</h2>
        <table>
            <tr><th>Path</th><th>Size</th><th>Created</th><th>Modified</th><th>SHA256</th></tr>
            {{range .EvidenceFiles}}
            <tr>
                <td>{{.Path}}</td>
                <td>{{.Size}}</td>
                <td>{{.Created}}</td>
                <td>{{.Modified}}</td>
                <td>{{.SHA256}}</td>
            </tr>
            {{end}}
        </table>
    </div>

    <div class="section">
        <h2>4. Timeline Forense</h2>
        <table>
            <tr><th>Time</th><th>Event</th><th>Source</th><th>Description</th></tr>
            {{range .Timeline}}
            <tr>
                <td>{{.Timestamp}}</td>
                <td>{{.Event}}</td>
                <td>{{.Source}}</td>
                <td>{{.Description}}</td>
            </tr>
            {{end}}
        </table>
    </div>

    <div class="section">
        <h2>5. Cadeia de Custódia (Master Hash: {{.MasterHash}})</h2>
        <table>
            <tr><th>Evidence File Name</th><th>Path</th><th>SHA256</th></tr>
            {{range .Evidences}}
            <tr>
                <td>{{.FileName}}</td>
                <td>{{.Path}}</td>
                <td>{{.Hash}}</td>
            </tr>
            {{end}}
        </table>
    </div>
</body>
</html>
`

	t, err := template.New("report").Parse(tpl)
	if err != nil {
		return models.Evidence{}, err
	}

	if err := t.Execute(file, report); err != nil {
		return models.Evidence{}, err
	}

	return models.Evidence{
		FileName: path,
		Path:     path,
	}, nil
}
