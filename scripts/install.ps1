#Requires -RunAsAdministrator
<#
.SYNOPSIS
    Install th (Terminal Help) CLI
    
.DESCRIPTION
    Installs the th CLI tool on Windows
    
.PARAMETER InstallDir
    Installation directory (default: $env:LOCALAPPDATA\th\bin)
    
.PARAMETER Force
    Force reinstall
    
.EXAMPLE
    .\install.ps1
    
.EXAMPLE
    .\install.ps1 -InstallDir "C:\Program Files\th" -Force
#>

param(
    [string]$InstallDir = "$env:LOCALAPPDATA\th\bin",
    [switch]$Force
)

$Repo = "DarioHefti/th"
$BinaryName = "th"

function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
        return $response.tag_name -replace '^v', ''
    } catch {
        Write-Error "Failed to get latest version: $_"
        exit 1
    }
}

function Get-Architecture {
    if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
        return "amd64"
    } elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
        return "arm64"
    } else {
        return "amd64"
    }
}

function Add-ToPath {
    param([string]$Path)
    
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$Path*") {
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$Path", "User")
        Write-Host "Added $Path to user PATH"
        Write-Host "Please restart your terminal for changes to take effect"
    }
}

function Main {
    $os = "windows"
    $arch = Get-Architecture
    $version = Get-LatestVersion
    
    if (-not $version) {
        Write-Error "Failed to get latest version"
        exit 1
    }
    
    $filename = "$BinaryName-$os-$arch.exe"
    $url = "https://github.com/$Repo/releases/download/v$version/$filename"
    
    if ((Test-Path "$InstallDir\$BinaryName.exe") -and -not $Force) {
        Write-Host "$BinaryName is already installed at $InstallDir" -ForegroundColor Yellow
        Write-Host "Use -Force to force reinstall" -ForegroundColor Yellow
        exit 0
    }
    
    Write-Host "Installing th v$version for $os/$arch..." -ForegroundColor Cyan
    
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    Write-Host "Downloading $url..."
    try {
        Invoke-WebRequest -Uri $url -OutFile "$InstallDir\$filename" -UseBasicParsing
    } catch {
        Write-Error "Failed to download: $_"
        exit 1
    }
    
    # Rename to th.exe
    Move-Item "$InstallDir\$filename" "$InstallDir\$BinaryName.exe" -Force
    
    Write-Host ""
    Write-Host "✓ Installed to $InstallDir\$BinaryName.exe" -ForegroundColor Green
    
    # Add to PATH if not already present
    Add-ToPath -Path $InstallDir
    
    Write-Host ""
    Write-Host "Run 'th --config' to set up your Azure credentials" -ForegroundColor Cyan
}

Main
