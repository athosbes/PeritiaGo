# PeritiaGo

Uma ferramenta portátil de Perícia Digital e Resposta a Incidentes (DFIR) para Windows, escrita em Go. Desenvolvida para ser nativamente independente, com foco em imutabilidade (preservação do alvo baseando-se em leituras limpas), rastreamento residual de softwares e rigorosa cadeia de custódia utilizando hashes SHA256 no modelo Master Manifest.

---

## 🚀 Capacidades Forenses Implementadas

1. **Detecção de Softwares Instalados**: Varredura direta das chaves do Registro (`HKLM` e `HKCU` - `Uninstall`) extraindo metadados.
2. **Evidência Visual (Autônoma)**: Abre silenciosamente o Painel de Controle (`appwiz.cpl`) e captura um _screenshot_ provando o estado visual dos programas listados pelo Sistema Operacional.
3. **Múltiplos Artefatos de Execução**:
   - **Prefetch**: Extração de rastros de execução nas pastas `C:\Windows\Prefetch`.
   - **Amcache**: Acionamento compatível com a ferramenta parser `AmcacheParser.exe`.
   - **ShimCache**: Extração direta do BLOB do Registro em `AppCompatCache`.
   - **UserAssist**: Coleta da execução GUI contendo decodificação automatizada em `ROT13`.
4. **Resíduos de Desinstalação**: Busca restos nas pastas `AppData`, `Program Files` e `ProgramData`.
5. **Busca Filesystem com Hashes**: Capaz de varrer drivers inteiros buscando combinações de nome de pastas ou extensões (como `.exe`, `.dll`, `.sqlite`) e gerando o log metadado + hash automaticamente.
6. **Timeline Cronométrica e Geração do Report**: Cruza datas de instalação, execução e modificação de arquivos, emitindo `.json`, `.csv`, e um Report `.html` lindamente formatado - coroando o diretório de saída com o `manifesto.txt` que possui o **Master Hash** inalterável.

---

## 🛠️ Como Compilar

Certifique-se de possuir o [Golang](https://go.dev/dl/) instalado e devidamente embutido na nas variáveis de ambiente `PATH`.

1. Inicialize ou baixe os módulos dependentes:
   ```powershell
   go mod tidy
   ```
2. Compile designando arquitetura alvo Windows:
   ```powershell
   $env:GOOS="windows"
   $env:GOARCH="amd64"
   go build -ldflags "-s -w" -o PeritiaGo.exe ./cmd/peritia
   ```
   > 💡 **Dica**: `ldflags "-s -w"` ajuda a remover símbolos de depuração, reduzindo dramaticamente o tamanho final, o que é ideal se quiser transportar só o `.exe` em um Pen-Drive forense.

---

## 📖 Instruções de Uso e Parâmetros

Para a execução e varredura correta sobre rastros que tocam a Raiz do Sistema e Registro, você precisará abrir um terminal `CMD` ou `PowerShell` localizando no **Modo Administrador**.

Exemplo de execução máxima com todos os argumentos:
```powershell
.\PeritiaGo.exe --case "Investigacao-HR-04" --investigator "Athos" --ext "exe,dll,vbs,sqlite" --search "AnyDesk" --drives "C:\,D:\Downloads"
```

### Argumentos:
- `--case`: O nome ou ID do caso em andamento (Padrão: `"Caso Padrao"`).
- `--investigator`: O nome do perito (Padrão: `"Perito"`).
- `--ext`: Lista delimitada por vírgula das ramificações ou tipos suspeitos que o aplicativo isolará e fará _hash_ se encontrar (Ex: `exe,log`).
- `--search`: Qualquer fragamento de texto (ignorando capitalização) que deseje buscar para pastas residuais de software (Ex: `AnyDesk`, `TeamViewer`).
- `--drives`: Pastas ou volumes em que deseja realizar as pesquisas profundas. Dependendo do tamanho do cenário, você limitará a profundidade. (Padrão: `C:\Users`).

A conclusão dos procedimentos criará uma pasta local na raiz executiva chamada `/outputs` contendo todos as evidências da custódia.

---

## 🔌 Estrutura Intuitiva para o Desenvolvedor

Esta ferramenta foi pensada na modularidade. Cada domínio de funcionalidade funciona no seu isolamento:
- `internal/config/`: Lidagem com os comandos CLI de entraveder.
- `internal/models/`: Contém as estruturas (Structs) globais e suas representações em JSON e CSV via *Tags*.
- `internal/capture/`: Funções intrincadas ao ambiente visual e nativo de listagem sistêmica.
- **`internal/artifacts/`**: Núcleo duro dos Parsers das trilhas forenses do alvo.
- `internal/export/`: Sistema generalista (`interface{}`) reflexivo de exportação de dados para disco e HTMLs.

### Quero adicionar um novo módulo/parser
O sistema absorverá novos coletores de forma transparente.
1. Desenvolva sua nova lógica (ex:`internal/artifacts/sysmon.go`)
2. Garanta que retorne `[]models.Artifact` populando as informações desejáveis.
3. Chame sua lógica abrindo `cmd/peritia/main.go` e atribuindo:
 `arts = append(arts, artifacts.SuaNovaRotinaOuFerramenta()...)`

O framework central passará automaticamente o seu tipo de _Artifact_ em ordem cronológica pelas Timeline, renderizará em HTML e calculará a Hash de custódia final!
