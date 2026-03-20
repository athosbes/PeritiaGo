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
        h1, h2, h3 { color: #2c3e50; margin-top: 0; }
        .header { background: #2980b9; color: white; padding: 20px; border-radius: 5px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .section { background: white; margin-top: 20px; padding: 20px; border-radius: 5px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
        table { width: 100%; border-collapse: collapse; margin-top: 15px; font-size: 14px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #34495e; color: white; position: sticky; top: 0; }
        tr:nth-child(even) { background-color: #f9f9f9; }
        tr:hover { background-color: #f1f1f1; }
        
        details { background: #fff; padding: 5px 10px; margin-bottom: 5px; }
        details summary { cursor: pointer; font-size: 20px; font-weight: bold; outline: none; list-style-type: none; border-bottom: 2px solid #ecf0f1; padding-bottom: 10px; }
        details summary::-webkit-details-marker { display: none; }
        details summary:before { content: "➕ "; font-size: 16px; }
        details[open] summary:before { content: "➖ "; }
        
        .controls { margin-top: 20px; display: flex; justify-content: space-between; align-items: center; background: white; padding: 15px; border-radius: 5px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
        #searchInput { padding: 10px; width: 400px; border: 1px solid #ccc; border-radius: 4px; font-size: 16px; outline: none; }
        #searchInput:focus { border-color: #2980b9; }
        .btn-print { padding: 10px 20px; background-color: #e67e22; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 16px; font-weight: bold; transition: background 0.3s; }
        .btn-print:hover { background-color: #d35400; }

        @media print {
            body { background: white; margin: 0; padding: 0; font-size: 12px; }
            .header { background: white; color: black; box-shadow: none; border-bottom: 2px solid black; }
            .section { box-shadow: none; border: 1px solid #ddd; page-break-inside: avoid; margin-top: 10px; }
            .controls { display: none; }
            details { display: block; }
            details summary:before { display: none; }
            table { page-break-inside: auto; }
            tr { page-break-inside: avoid; page-break-after: auto; }
            thead { display: table-header-group; }
            tfoot { display: table-footer-group; }
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>PeritiaGo - Forensic Report</h1>
        <p><strong>Case:</strong> {{.CaseName}}</p>
        <p><strong>Investigator:</strong> {{.Investigator}}</p>
        <p><strong>Date:</strong> {{.CaptureDate}}</p>
        <hr style="border-top: 1px solid #bdc3c7; border-bottom: none; margin: 15px 0;">
        <p><strong>Machine Name:</strong> {{.Machine.Hostname}}</p>
        <p><strong>OS:</strong> {{.Machine.OSName}} {{.Machine.OSVersion}} (Build {{.Machine.OSBuild}})</p>
        <p><strong>Machine GUID:</strong> {{.Machine.MachineGUID}}</p>
        <p><strong>IP Addresses:</strong> {{range .Machine.IPAddresses}}{{.}} {{end}}</p>
        <p><strong>MAC Addresses:</strong> {{range $key, $value := .Machine.MACAddresses}}{{$key}}: {{$value}} | {{end}}</p>
    </div>

    <div class="controls">
        <input type="text" id="searchInput" placeholder="Pesquisar em todas as tabelas e expandir..." onkeyup="filterTables()">
        <button class="btn-print" onclick="window.print()">🖨️ Imprimir Laudo / Exportar PDF</button>
    </div>

    <div class="section">
        <details open>
            <summary>1. Softwares Detectados ({{len .InstalledSofts}})</summary>
            <div style="overflow-x:auto;">
                <table class="searchable">
                    <thead><tr><th>Name</th><th>Version</th><th>Publisher</th><th>Install Date</th><th>Source</th></tr></thead>
                    <tbody>
                    {{range .InstalledSofts}}
                    <tr>
                        <td>{{.DisplayName}}</td>
                        <td>{{.DisplayVersion}}</td>
                        <td>{{.Publisher}}</td>
                        <td>{{.InstallDate}}</td>
                        <td>{{.Source}}</td>
                    </tr>
                    {{end}}
                    </tbody>
                </table>
            </div>
        </details>
    </div>

    <div class="section">
        <details>
            <summary>2. Artefatos de Execução e Sistema ({{len .Artifacts}})</summary>
            <div style="overflow-x:auto;">
                <table class="searchable">
                    <thead><tr><th>Name</th><th>Type</th><th>Description</th><th>Timestamp</th></tr></thead>
                    <tbody>
                    {{range .Artifacts}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{.Type}}</td>
                        <td>{{.Description}}</td>
                        <td>{{.Timestamp}}</td>
                    </tr>
                    {{end}}
                    </tbody>
                </table>
            </div>
        </details>
    </div>

    <div class="section">
        <details>
            <summary>3. Busca por Extensão / Residuais ({{len .EvidenceFiles}})</summary>
            <div style="overflow-x:auto;">
                <table class="searchable">
                    <thead><tr><th>Path</th><th>Size</th><th>Created</th><th>Modified</th><th>SHA256</th></tr></thead>
                    <tbody>
                    {{range .EvidenceFiles}}
                    <tr>
                        <td>{{.Path}}</td>
                        <td>{{.Size}}</td>
                        <td>{{.Created}}</td>
                        <td>{{.Modified}}</td>
                        <td>{{.SHA256}}</td>
                    </tr>
                    {{end}}
                    </tbody>
                </table>
            </div>
        </details>
    </div>

    <div class="section">
        <details>
            <summary>4. Timeline Forense Cronométrica ({{len .Timeline}})</summary>
            <div style="overflow-x:auto;">
                <table class="searchable">
                    <thead><tr><th>Time</th><th>Event</th><th>Source</th><th>Description</th></tr></thead>
                    <tbody>
                    {{range .Timeline}}
                    <tr>
                        <td>{{.Timestamp}}</td>
                        <td>{{.Event}}</td>
                        <td>{{.Source}}</td>
                        <td>{{.Description}}</td>
                    </tr>
                    {{end}}
                    </tbody>
                </table>
            </div>
        </details>
    </div>

    <div class="section">
        <details>
            <summary>5. Cadeia de Custódia e Arquivos Físicos ({{len .Evidences}})</summary>
            <p><strong>Master Hash SHA256 do Diretório:</strong> {{.MasterHash}}</p>
            <div style="overflow-x:auto;">
                <table class="searchable">
                    <thead><tr><th>Evidence File Name</th><th>Path</th><th>SHA256</th></tr></thead>
                    <tbody>
                    {{range .Evidences}}
                    <tr>
                        <td>{{.FileName}}</td>
                        <td>{{.Path}}</td>
                        <td>{{.Hash}}</td>
                    </tr>
                    {{end}}
                    </tbody>
                </table>
            </div>
        </details>
    </div>

    <script>
        function filterTables() {
            var input = document.getElementById("searchInput");
            var filter = input.value.toUpperCase();
            var tables = document.getElementsByClassName("searchable");
            
            // Se houver busca ativa, expanda todas as sections para mostrar os resultados
            var detailsList = document.querySelectorAll('details');
            if (filter.length > 0) {
                detailsList.forEach(function(detail) {
                    detail.setAttribute('open', '');
                });
            }
            
            for (var t = 0; t < tables.length; t++) {
                var tbody = tables[t].getElementsByTagName("tbody")[0];
                var tr = tbody.getElementsByTagName("tr");
                
                for (var i = 0; i < tr.length; i++) {
                    var tds = tr[i].getElementsByTagName("td");
                    var match = false;
                    for (var j = 0; j < tds.length; j++) {
                        if (tds[j]) {
                            var txtValue = tds[j].textContent || tds[j].innerText;
                            if (txtValue.toUpperCase().indexOf(filter) > -1) {
                                match = true;
                                break;
                            }
                        }
                    }
                    tr[i].style.display = match ? "" : "none";
                }
            }
        }
    </script>
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
