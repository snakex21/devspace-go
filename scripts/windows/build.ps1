# Build Dev Space Go for all platforms
# Run from project root: .\scripts\windows\build.ps1

$ErrorActionPreference = "Stop"
Set-Location (Split-Path (Split-Path $PSScriptRoot -Parent) -Parent)

$targets = @(
    @{OS="windows"; Arch="amd64";   Ext=".exe";   Dir="windows"},
    @{OS="linux";   Arch="amd64";   Ext="";        Dir="linux"},
    @{OS="darwin";  Arch="amd64";   Ext="";        Dir="macos-intel"},
    @{OS="darwin";  Arch="arm64";   Ext="";        Dir="macos-mchip"}
)

Write-Output "=== Dev Space Go Build All ==="
Write-Output ""

# Clean build dir
if (Test-Path build) { Remove-Item build -Recurse -Force }
New-Item -ItemType Directory -Path build -Force | Out-Null

$originalGOOS = $env:GOOS
$originalGOARCH = $env:GOARCH
try {
    foreach ($t in $targets) {
        $platformDir = "build\$($t.Dir)"
        New-Item -ItemType Directory -Path $platformDir -Force | Out-Null

        $env:GOOS = $t.OS
        $env:GOARCH = $t.Arch

        # Main server (lightweight, no GUI)
        Write-Output "  [$($t.Dir)] devspace (serwer)..."
        go build -o "$platformDir\devspace$($t.Ext)" ./cmd/devspace/

        # Fyne GUI must be compiled natively. This Windows script builds only the
        # Windows GUI and leaves Linux/macOS GUI builds to those platforms.
        if ($t.OS -eq "windows") {
            Write-Output "  [$($t.Dir)] devspace-gui (konfigurator)..."
            go build -o "$platformDir\devspace-gui$($t.Ext)" ./cmd/devspace-gui/
        }
        else {
            Write-Output "  [$($t.Dir)] devspace-gui: skipped (native build required)"
        }

        Write-Output ""
    }
}
finally {
    $env:GOOS = $originalGOOS
    $env:GOARCH = $originalGOARCH
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

Write-Output "Summary: server (4 platforms) + GUI (current platform)"
Write-Output ""
Write-Output "NOTE: Fyne GUI builds only on the current platform."
Write-Output "Build the GUI natively on Linux and macOS."
