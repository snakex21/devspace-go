# Build DevSpace for all platforms
# Run from project root: .\scripts\windows\build.ps1

$ErrorActionPreference = "Stop"
Set-Location (Split-Path (Split-Path $PSScriptRoot -Parent) -Parent)

$targets = @(
    @{OS="windows"; Arch="amd64";   Ext=".exe";   Dir="windows"},
    @{OS="linux";   Arch="amd64";   Ext="";        Dir="linux"},
    @{OS="darwin";  Arch="amd64";   Ext="";        Dir="macos-intel"},
    @{OS="darwin";  Arch="arm64";   Ext="";        Dir="macos-mchip"}
)

Write-Output "=== DevSpace Build All ==="
Write-Output ""

# Clean build dir
if (Test-Path build) { Remove-Item build -Recurse -Force }
New-Item -ItemType Directory -Path build -Force | Out-Null

foreach ($t in $targets) {
    $platformDir = "build\$($t.Dir)"
    New-Item -ItemType Directory -Path $platformDir -Force | Out-Null

    $env:GOOS = $t.OS
    $env:GOARCH = $t.Arch

    # Main server (lightweight, no GUI)
    Write-Output "  [$($t.Dir)] mcp-webcoder (serwer)..."
    go build -o "$platformDir\mcp-webcoder$($t.Ext)" ./cmd/devspace/

    # GUI configurator (Fyne — może nie działać cross-platform)
    Write-Output "  [$($t.Dir)] mcp-webcoder-gui (konfigurator)..."
    go build -o "$platformDir\mcp-webcoder-gui$($t.Ext)" ./cmd/devspace-gui/ 2>$null
    if ($LASTEXITCODE -ne 0) {
        Write-Output "         (GUI niedostępne na tej platformie przy cross-kompilacji)"
    }

    Write-Output ""
}

# Portable tools copied next to release binaries.
if (Test-Path -LiteralPath "tools") {
    foreach ($t in $targets) {
        $platformTools = "build\$($t.Dir)\tools"
        New-Item -ItemType Directory -Path $platformTools -Force | Out-Null

        if ($t.OS -eq "windows" -and (Test-Path -LiteralPath "tools\cloudflared.exe")) {
            Copy-Item -LiteralPath "tools\cloudflared.exe" -Destination "$platformTools\cloudflared.exe" -Force
        }
        elseif (Test-Path -LiteralPath "tools\cloudflared") {
            Copy-Item -LiteralPath "tools\cloudflared" -Destination "$platformTools\cloudflared" -Force
        }
    }
}

Write-Output "=== Done! ==="
Write-Output ""

foreach ($t in $targets) {
    $dir = "build\$($t.Dir)"
    Write-Output "--- $($t.Dir) ---"
    Get-ChildItem $dir -File | ForEach-Object {
        Write-Output "  $($_.Name) ($([math]::Round($_.Length/1MB, 1)) MB)"
    }
    Write-Output ""
}

Write-Output "Razem: serwer (4 platformy) + GUI (bieżąca platforma)"
Write-Output ""
Write-Output "UWAGA: GUI (Fyne) działa tylko na platformie na której kompilujesz."
Write-Output "Na Linux/macOS skompiluj GUI natywnie na tych systemach."
