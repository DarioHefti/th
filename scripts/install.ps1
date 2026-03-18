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
$MaxRetries = 3
$RetryDelay = 2

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Cyan
}

function Write-ErrorMsg {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing -TimeoutSec 30
        $version = $response.tag_name -replace '^v', ''
        
        if ($version -notmatch '^\d+\.\d+\.\d+$') {
            Write-ErrorMsg "Invalid version format: $version"
            return $null
        }
        
        return $version
    } catch {
        Write-ErrorMsg "Failed to get latest version: $_"
        return $null
    }
}

function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        "x86"   { return "386" }
        default { 
            Write-ErrorMsg "Unsupported architecture: $arch"
            return $null
        }
    }
}

function Test-Url {
    param([string]$Url)
    
    try {
        $response = Invoke-WebRequest -Uri $Url -Method Head -UseBasicParsing -TimeoutSec 10 -ErrorAction SilentlyContinue
        return $response.StatusCode -eq 200
    } catch {
        return $false
    }
}

function Add-ToPath {
    param([string]$Path)
    
    try {
        $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($currentPath -notlike "*$Path*") {
            [Environment]::SetEnvironmentVariable("Path", "$currentPath;$Path", "User")
            Write-Info "Added $Path to user PATH"
            Write-Info "Please restart your terminal for changes to take effect"
        }
    } catch {
        Write-ErrorMsg "Failed to update PATH: $_"
    }
}

function Test-Binary {
    param([string]$Path)
    
    if (-not (Test-Path $Path)) {
        Write-ErrorMsg "Binary not found at $Path"
        return $false
    }
    
    $fileInfo = Get-Item $Path
    if ($fileInfo.Length -eq 0) {
        Write-ErrorMsg "Binary is empty"
        return $false
    }
    
    return $true
}

function Download-Binary {
    param(
        [string]$Url,
        [string]$OutputPath
    )
    
    if (-not (Test-Url $Url)) {
        Write-ErrorMsg "Release not found at: $Url"
        return $false
    }
    
    $attempt = 1
    $success = $false
    
    while ($attempt -le $MaxRetries) {
        Write-Info "Download attempt $attempt of $MaxRetries..."
        
        try {
            $ProgressPreference = 'SilentlyContinue'
            Invoke-WebRequest -Uri $Url -OutFile $OutputPath -UseBasicParsing -TimeoutSec 60 -ErrorAction Stop
            $ProgressPreference = 'Normal'
            
            if ((Test-Path $OutputPath) -and ((Get-Item $OutputPath).Length -gt 0)) {
                $success = $true
                break
            } else {
                Write-ErrorMsg "Downloaded file is empty or missing"
            }
        } catch {
            Write-ErrorMsg "Download attempt $attempt failed: $_"
        }
        
        if ($attempt -lt $MaxRetries) {
            Write-Info "Retrying in ${RetryDelay}s..."
            Start-Sleep -Seconds $RetryDelay
        }
        
        $attempt++
    }
    
    $ProgressPreference = 'Normal'
    
    if ($success) {
        Write-Info "Download complete"
        return $true
    }
    
    return $false
}

function Main {
    if (-not (Test-Path $InstallDir)) {
        try {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        } catch {
            Write-ErrorMsg "Failed to create installation directory: $_"
            exit 1
        }
    }
    
    if (-not (Test-Path $InstallDir -PathType Container)) {
        Write-ErrorMsg "Install directory is not accessible: $InstallDir"
        exit 1
    }
    
    $os = "windows"
    $arch = Get-Architecture
    
    if (-not $arch) {
        exit 1
    }
    
    Write-Info "Fetching latest version..."
    $version = Get-LatestVersion
    
    if (-not $version) {
        Write-ErrorMsg "Failed to get latest version"
        Write-ErrorMsg "Please check your network connection and try again"
        exit 1
    }
    
    $filename = "$BinaryName-$os-$arch.exe"
    $url = "https://github.com/$Repo/releases/download/v$version/$filename"
    $outputPath = Join-Path $InstallDir $filename
    $finalPath = Join-Path $InstallDir "$BinaryName.exe"
    
    if ((Test-Path $finalPath) -and -not $Force) {
        Write-Info "$BinaryName is already installed at $InstallDir"
        Write-Info "Use -Force to force reinstall"
        exit 0
    }
    
    Write-Info "Installing th v$version for $os/$arch..."
    Write-Info "Downloading $filename..."
    
    if (-not (Download-Binary -Url $url -OutputPath $outputPath)) {
        Write-ErrorMsg "Failed to download after $MaxRetries attempts"
        if (Test-Path $outputPath) {
            Remove-Item $outputPath -Force
        }
        exit 1
    }
    
    if (Test-Path $finalPath) {
        Remove-Item $finalPath -Force
    }
    
    try {
        Move-Item -Path $outputPath -Destination $finalPath -Force
    } catch {
        Write-ErrorMsg "Failed to install binary: $_"
        exit 1
    }
    
    if (-not (Test-Binary -Path $finalPath)) {
        Write-ErrorMsg "Binary verification failed"
        if (Test-Path $finalPath) {
            Remove-Item $finalPath -Force
        }
        exit 1
    }
    
    Write-Host ""
    Write-Host "✓ Installed to $finalPath" -ForegroundColor Green
    
    Add-ToPath -Path $InstallDir
    
    Write-Host ""
    Write-Host "Run 'th --config' to set up your Azure credentials" -ForegroundColor Cyan
}

Main
