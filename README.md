# PeritiaGo 🕵️‍♂️ (DFIR Artifact Collector)
🌐 [English Version](#english) | 🇧🇷 [Versão em Português](#portugues)

⚠️ **LEGAL DISCLAIMER / AVISO LEGAL:**  
*PeritiaGo is exclusively an **Evidence Acquisition Tool**. It blindly collects raw machine state and historical execution data. The generated reports do not constitute a legal conclusion. Any findings must be interpreted and contextualized by a qualified Digital Forensics Expert (Perito) to draft a formal and court-admissible forensic report (Laudo).*

---
<a name="english"></a>
# English

A portable Digital Forensics and Incident Response (DFIR) tool for Windows written in Go. Built to be completely standalone, focusing heavily on **Immutability** (clean system reads with zero contamination), strict **Chain of Custody** hashing, and automated **Deep Artifact Triangulation**. 

It generates air-gapped, searchable HTML/PDF reports and fully structured `.json`/`.csv` manifests natively.

## 🚀 Forensic Capabilities & "The Why"

PeritiaGo doesn't just collect data blindly. It extracts specific, court-admissible artifacts critical for attributing actions to a suspect or mapping malware behavior.

1. **Machine Identity & Environment Traces**
   - **What it does:** Extracts exact *OS Build*, IP & MAC addresses, and the *Machine GUID*. Tests the registry and `C:\Windows.old` for recent OS wipes or "Reset this PC" commands.
   - **Why it matters:** IP/MAC/GUID explicitly ties the extracted forensic evidence to a physical suspect and a corporate subnet. The "Reset" checking prevents suspects from claiming "I know nothing, the PC was wiped" by proving *when* they tried to erase evidence.

2. **Unlicensed Software Auditing (Piracy Detection)**
   - **What it does:** Captures **Active Live Processes** (Volatile RAM) to catch portable crack engines or unlicensed CAD/Design software currently running without installation. Crosschecks the hard drive for `--ext ".lic,.dll,.crack"` or by exact software name `--search "Photoshop"`.
   - **Why it matters:** Even if an employee runs a pirated portable software off a USB stick, the volatile memory capture combined with Amcache triangulates the exact time and frequency of use, providing devastating proof of unlicensed software usage within the corporate network to aid in licensing audits.

3. **The Execution Triad (Prefetch, ShimCache, Amcache)**
   - **What it does:** Reads and extracts binaries from `C:\Windows\Prefetch`, `AppCompatCache` blobs, and automates `Amcache.hve` parsing through Eric Zimmerman's engines.
   - **Why it matters:** Even if a suspect double-clicks a malware/unauthorized remote tool (like `anydesk.exe`) and then securely deletes it from the hard drive, Windows irreversibly logs the execution time, path, and frequency in these three artifacts. They establish the undisputable timeline of what ran, and when.

4. **Software & Uninstall Residuals**
   - **What it does:** Uses `wmic`, `winget`, and deep registry parsing (`HKLM`, `HKCU`). Automatically generates a silent visual screenshot of "Add/Remove Programs".
   - **Why it matters:** Identifies persistence mechanisms, unauthorized VPNs, tools used for lateral movement, or recently applied hotfixes (which might indicate the machine was vulnerable moments prior).

5. **UserAssist & Filesystem Hashing**
   - **What it does:** Decodes GUI application launches (`ROT13`) bypassing obfuscation. Scans targeted storage volumes looking for suspicious custom extensions (`.vbs`, `.sqlite`).
   - **Why it matters:** UserAssist proves exactly what the active UI user explicitly clicked. The filesystem scanner isolates payloads or staged databases.

6. **The Master Manifest (Chain of Custody)**
   - **What it does:** Every file generated (CSV exports, screenshots, JSON databases) inside the dynamic output folder `software_inventory_UUID_MAC_TIMESTAMP` receives a SHA-256 hash. Everything is sealed inside `manifesto.txt` with a final Master Directory Hash.
   - **Why it matters:** ISO 27037 / RFC 3161 Compliance. It acts as legal cryptography, denying anti-forensic repudiation claims in court.

---

## 🔌 Integrating `AmcacheParser.exe` (Eric Zimmerman)
For PeritiaGo to automatically read and convert `Amcache.hve` entries into your native HTML report, you must provide Eric Zimmerman's parser next to the executable.

We implemented a **Robust Fallback Engine**. PeritiaGo will natively try to execute the modern `.NET 9` version, log and skip if it fails, and automatically fall back to the `.NET 4` (older machines) version. 

**Required Output Folder Structure (Deploy it exactly like this to your Pen-Drive/Drive):**
```text
/PeritiaGo
  ┣ PeritiaGo.exe
  ┣ AmcacheParsernet9/
  ┃ ┗ AmcacheParser.exe  <-- (Priority 1)
  ┗ AmcacheParsernet4/
    ┗ AmcacheParser.exe  <-- (Priority 2 / Safe Fallback)
```
*Note: The executable also checks its own root folder and the global system `$PATH` if the subfolders are missing.*

---

## 🛠️ Deployment & Execution

You must strictly execute PeritiaGo using `CMD` or `PowerShell` in **Administrator Mode** to access OS Root/Registry cores.

### 1. Zero-Touch JSON Config (Recommended for Networks)
Drop `PeritiaGo.exe` in a clean folder and run it. If `peritiago_config.json` is missing, it creates a template. Fill this JSON out! Any subsequent executions across machines will instantly import these parameters silently.

### 2. Standalone UI & CLI Flags
If you run it without arguments, a native Windows GUI Prompt will pop up asking for the Case ID, Investigator, Drives, etc.
Alternatively, wire it up to your EDR pipelines:
```powershell
.\PeritiaGo.exe --case "HR-04" --investigator "John Doe" --ext "exe,dll,vbs" --search "AnyDesk" --drives "C:\Users,D:\"
```

---
---

<a name="portugues"></a>
# Português

Uma ferramenta portátil de Perícia Digital e Resposta a Incidentes (DFIR) para Windows, escrita em Go. Desenvolvida para ser nativamente independente, com foco em **Imutabilidade** (preservação do alvo baseando-se em leituras limpas sem contaminação), forte documentação baseada em **Cadeia de Custódia**, e automação profunda de **Artefatos Cronométricos**.

## 🚀 Capacidades Forenses e "O Porquê"

O PeritiaGo não realiza varreduras às cegas. Ele foca cirurgicamente em artefatos válidos e admitidos em tribunais (para provar furtos de dados, pirataria, execução de malwares ou ocultação).

1. **Digital do Sistema e Ambientes (Identidade)**
   - **Como Faz:** Coleta nativamente a *Build exata do SO*, IP/MAC, e o *GUID da Máquina*. Captura evidências visuais silenciosas das especificações e vasculha o Registro ou `C:\Windows.old` atrás de formatações ("Reset this PC").
   - **Por Que Importa:** Sem IPs associados a um MAC Address, o advogado de defesa ou tribunal pode alegar que a máquina avaliada não é a do suspeito. A detecção de formatação destrói a alegação "eu não formatei para esconder provas". O sistema expõe agressivamente a data original de instalação em contraponto à formatação recente.

2. **Auditoria de Softwares Não-Licenciados (Pirataria)**
   - **Como Faz:** Captura todos os **Processos Voláteis (RAM)** rodando em tempo real na máquina, expondo geradores portáteis ("Portable") ou ativadores em execução sem estarem instalados. Cruza isso varrendo o sistema todo atrás de extensões focadas em quebra de licença (`--ext ".lic,.dll,.crack"`) ou pelo nome flagrado (`--search "Photoshop"`).
   - **Por Que Importa:** Ideal para avaliações cíveis ou auditorias corporativas. O PeritiaGo captura trilhas vivas de ram e passadas de registro para provar concretamente que aquele software sem nota fiscal/licença foi intencionalmente executado na máquina e por quantas vezes.

3. **A Tríade de Execução (Prefetch, ShimCache, Amcache)**
   - **Como Faz:** Extrai binários isolados dos rastros pré-alocados em `C:\Windows\Prefetch`, injeta e decodifica o blob contido no Regedit (`AppCompatCache`) e parseia profundamente a chave `Amcache.hve` de forma automatizada.
   - **Por Que Importa:** *Mesmo se o suspeito utilizar um programa pendrive malicioso (exemplo `Mimikatz` ou `AnyDesk`) e em seguida apagá-lo.* O Windows registra e acusa irreversivelmente a existência temporal e a quantidade de aberturas desses executáveis nesta Tríade de rastros invisíveis.

4. **Inventário e Ocultação de Software**
   - **Como Faz:** Junta chamadas em lote (`wmic`, `winget`, Registro) e gera instantaneamente o print visual oculto da aba "Adicionar e Remover Programas" do usuário.
   - **Por Que Importa:** Ajuda analistas a cruzarem ferramentas de by-pass instaladas, VPNs em shadow-IT, e lista precisamente as "Windows Updates (Hotfixes)" presentes; essencial para relatar se a máquina estava vulnerabilizada por falha de patches conhecidos durante o momento da intrusão.

5. **Busca Focada (Extensões) e UserAssist**
   - **Como Faz:** Descriptografa com `ROT13` quais aplicações dotadas de GUI foram clicadas nos atalhos e faz buscas brutas em disco em extensões avulsas suspeitas (`.vbs`, `.sqlite`).
   - **Por Que Importa:** O `UserAssist` garante a "Ação Comprovada e Consciente" em GUI. O scanner de arquivo expõe payloads adormecidos que caçadores de ameaças ou indicadores de compromisso determinarem na varredura.

6. **A Árvore Mestra (Manifesto e Cadeia de Custódia)**
   - **Como Faz:** Todos os extratos (JSONs, relatórios CSVs, e HTML offline) entram num cofre nominado dinamicamente (`software_inventory_UUID_MAC_TIMESTAMP`). O algoritmo mastiga pasta por pasta no fim e assina rigorosamente todos os arquivos isolados com hashes SHA-256 e aglutina isso blindado num `manifesto.txt`.
   - **Por Que Importa:** Cumprimento da RFC 3161 / Cadeias Legais (ISO 27037). O perito adquire força intocável de Não-Repúdio; provando que nenhuma linha do laudo digital logado ali sofreu adulteraçõs após a captura no pendrive.

---

## 🔌 Estruturando a Lógica Tripla do Amcache (Eric Zimmerman)
Para a ferramenta engolir nativamente todos os logs complexos do `Amcache.hve` e inseri-los perfeitamente no nosso `HTML` (offline e expansivo), o executor requer emparelhar o conversor do Zimmerman.

Nós enraizamos um **Motor de Fallback Robusto**. O `PeritiaGo.exe` tentará sempre ativar o motor mais ágil contemporâneo (`.NET 9`). Se o sistema for muito arcaico ou faltar pacotes e a extração "quebrar" ou falhar invisivelmente, nossa ferramenta intercepta a falha e invoca a versão garantida (`.NET 4`).

Para isso ser possível, disponha sua árvore de execução de Pendrive **exatamente dessa forma:**
```text
/QualquerPasta
  ┣ PeritiaGo.exe
  ┣ AmcacheParsernet9/
  ┃ ┗ AmcacheParser.exe  <-- (Prioridade Base de Desempenho)
  ┗ AmcacheParsernet4/
    ┗ AmcacheParser.exe  <-- (O Fallback / Caminho Universal Alternativo Seguro)
```

---

## 🛠️ Execução na Máquina Alvo

Certifique-se sempre que rodar o prompt em **Modo Administrador** (`CMD` ou `PowerShell`) para liberar as chaves restritas da base do Sistema Operacional (.hve).

### 1. Configuração Silenciosa JSON (Redes e Multiplicidade)
Jogue o executável num diretório solitário e chame-o vazio. Se o `peritiago_config.json` inexiste, ele cria o molde na hora. Altere esse arquivo. Se acoplado ao pen-drive, os próximos PC's que você abrir e plugar absorverão inteiramente seu arquétipo sem você tocar uma flag manual sequer.

### 2. Automação e Frentes CLI/GUI
Se o invocar cru na pasta isolada sem json sem pre-formatalo, o PeritiaGo salda a sua salvação abrindo um InputBox em Windows nativo. Preencha ali ou acione a robustez de pipeline:
```powershell
.\PeritiaGo.exe --case "Incidente HR-44" --investigator "Athos" --ext "exe,dll,vbs" --search "Trojan.Agent" --drives "C:\Users,D:\"
```
O laudo HTML super rico contido na finalização suporta buscas offlines em javascript nativo e possui rotinas customizadas acionadas ao clicar no botão embutido de "Exportar PDF / Imprimir" onde expele um Relatório de Investigação completo.
