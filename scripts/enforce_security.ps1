# Security Enforcement Script for PeritiaGo
# This script ensures that no forbidden direct write calls are made outside the filesystem package.

$forbidden = @("os.WriteFile", "os.Create", "ioutil.WriteFile")
$allowedPackage = "internal\filesystem"
$configPackage = "internal\config"
$loggerPackage = "internal\logger"
$testFiles = "*_test.go"

$exitCode = 0

Write-Host "Checking for forbidden direct write calls..." -ForegroundColor Cyan

$files = Get-ChildItem -Recurse -Include *.go | Where-Object { 
    $_.FullName -notmatch [regex]::Escape($allowedPackage) -and 
    $_.FullName -notmatch [regex]::Escape($configPackage) -and 
    $_.FullName -notmatch [regex]::Escape($loggerPackage) -and 
    $_.Name -notlike $testFiles
}

foreach ($file in $files) {
    $content = Get-Content $file.FullName
    foreach ($call in $forbidden) {
        if ($content -match [regex]::Escape($call)) {
            Write-Host "Error: Forbidden call '$call' found in $($file.FullName)" -ForegroundColor Red
            $exitCode = 1
        }
    }
}

if ($exitCode -eq 0) {
    Write-Host "Success: No forbidden calls detected." -ForegroundColor Green
} else {
    Write-Host "Failure: Security policy violated. Please use 'internal/filesystem' wrappers." -ForegroundColor Red
}

exit $exitCode
