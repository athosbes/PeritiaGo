# PeritiaGo Code Signing Script - Sigstore Edition
# Usage: .\sign_binary.ps1 -BinaryPath ".\peritia.exe" [-IdToken "YOUR_OIDC_TOKEN"]
#
# Sigstore replaces Microsoft SignTool. Free of charge, no certificate purchase needed.
# Uses OIDC identity (GitHub, Google, Microsoft) for signing.
#
# In GitHub Actions CI/CD, the OIDC token is provided automatically via
# the SIGSTORE_ID_TOKEN environment variable.

param (
    [Parameter(Mandatory=$true)]
    [string]$BinaryPath,

    [Parameter(Mandatory=$false)]
    [string]$IdToken
)

# Check for existing Sigstore signature
$bundlePath = "$BinaryPath.sigstore.json"
if (Test-Path $bundlePath) {
    Write-Host "Verifying existing Sigstore signature..." -ForegroundColor Yellow
    & go run ./cmd/peritiatool verify $BinaryPath
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Binary already has a valid Sigstore signature." -ForegroundColor Green
        exit 0
    }
    Write-Host "Existing signature invalid, re-signing..." -ForegroundColor Yellow
}

# Resolve ID token from parameter or environment
if (-not $IdToken) {
    $IdToken = $env:SIGSTORE_ID_TOKEN
}

if (-not $IdToken) {
    Write-Error "No OIDC ID token provided. Set -IdToken parameter or SIGSTORE_ID_TOKEN env var."
    Write-Error "  GitHub Actions: uses automatic OIDC"
    Write-Error "  Local: use 'gcloud auth print-identity-token' or similar"
    exit 1
}

Write-Host "Signing binary with Sigstore: $BinaryPath"
$env:SIGSTORE_ID_TOKEN = $IdToken
& go run ./cmd/peritiatool sign $BinaryPath --id-token $IdToken

if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: Sigstore bundle saved to $bundlePath" -ForegroundColor Green

    # Verify the signature we just created
    Write-Host "Verifying signature..."
    & go run ./cmd/peritiatool verify $BinaryPath
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Verification passed." -ForegroundColor Green
    } else {
        Write-Error "Verification failed after signing!"
    }
} else {
    Write-Error "Signing failed!"
}
