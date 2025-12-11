# dotclaude Windows Installer
# PowerShell script for installing dotclaude on Windows
#
# Usage:
#   irm https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.ps1 | iex
#
# Or download and run:
#   .\install.ps1

$ErrorActionPreference = "Stop"

# Configuration
$DOTCLAUDE_VERSION = "latest"
$INSTALL_DIR = "$env:USERPROFILE\.local\bin"
$REPO_DIR = "$env:USERPROFILE\code\dotclaude"
$CLAUDE_DIR = "$env:USERPROFILE\.claude"
$GITHUB_REPO = "blackwell-systems/dotclaude"

function Write-Header {
    Write-Host ""
    Write-Host "  +-------------------------------------------------------------+" -ForegroundColor Cyan
    Write-Host "  |  dotclaude - Profile Management for Claude Code             |" -ForegroundColor Cyan
    Write-Host "  +-------------------------------------------------------------+" -ForegroundColor Cyan
    Write-Host ""
}

function Write-Step {
    param([string]$Message)
    Write-Host "  [*] $Message" -ForegroundColor Green
}

function Write-Info {
    param([string]$Message)
    Write-Host "      $Message" -ForegroundColor Gray
}

function Write-Warning {
    param([string]$Message)
    Write-Host "  [!] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "  [X] $Message" -ForegroundColor Red
}

function Get-Architecture {
    if ([Environment]::Is64BitOperatingSystem) {
        return "amd64"
    }
    return "386"
}

function Get-LatestVersion {
    try {
        $releases = Invoke-RestMethod -Uri "https://api.github.com/repos/$GITHUB_REPO/releases/latest" -Headers @{Accept = "application/vnd.github.v3+json"}
        return $releases.tag_name
    }
    catch {
        Write-Warning "Could not fetch latest version, using 'latest'"
        return "latest"
    }
}

function Install-Binary {
    Write-Step "Installing dotclaude binary..."

    # Create install directory
    if (-not (Test-Path $INSTALL_DIR)) {
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
        Write-Info "Created directory: $INSTALL_DIR"
    }

    # Determine version and architecture
    $arch = Get-Architecture
    $version = if ($DOTCLAUDE_VERSION -eq "latest") { Get-LatestVersion } else { $DOTCLAUDE_VERSION }
    $versionNumber = $version -replace '^v', ''

    # Download URL
    $downloadUrl = "https://github.com/$GITHUB_REPO/releases/download/$version/dotclaude_${versionNumber}_windows_${arch}.zip"
    $tempZip = "$env:TEMP\dotclaude.zip"
    $tempDir = "$env:TEMP\dotclaude-extract"

    Write-Info "Downloading from: $downloadUrl"

    try {
        # Download
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempZip -UseBasicParsing

        # Extract
        if (Test-Path $tempDir) {
            Remove-Item -Recurse -Force $tempDir
        }
        Expand-Archive -Path $tempZip -DestinationPath $tempDir -Force

        # Copy binary
        $binary = Get-ChildItem -Path $tempDir -Filter "dotclaude.exe" -Recurse | Select-Object -First 1
        if ($binary) {
            Copy-Item -Path $binary.FullName -Destination "$INSTALL_DIR\dotclaude.exe" -Force
            Write-Info "Installed: $INSTALL_DIR\dotclaude.exe"
        }
        else {
            throw "Could not find dotclaude.exe in archive"
        }

        # Cleanup
        Remove-Item -Path $tempZip -Force -ErrorAction SilentlyContinue
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    catch {
        Write-Warning "Binary download failed: $_"
        Write-Warning "Trying to build from source..."
        Install-FromSource
    }
}

function Install-FromSource {
    Write-Step "Building from source..."

    # Check for Go
    $go = Get-Command go -ErrorAction SilentlyContinue
    if (-not $go) {
        Write-Error "Go is not installed. Please install Go from https://go.dev/dl/"
        Write-Error "Or download pre-built binary from: https://github.com/$GITHUB_REPO/releases"
        exit 1
    }

    # Clone or update repository
    if (Test-Path $REPO_DIR) {
        Write-Info "Updating existing repository..."
        Push-Location $REPO_DIR
        git pull origin main
        Pop-Location
    }
    else {
        Write-Info "Cloning repository..."
        git clone "https://github.com/$GITHUB_REPO.git" $REPO_DIR
    }

    # Build
    Push-Location $REPO_DIR
    go build -o "$INSTALL_DIR\dotclaude.exe" .\cmd\dotclaude
    Pop-Location

    Write-Info "Built: $INSTALL_DIR\dotclaude.exe"
}

function Setup-Repository {
    Write-Step "Setting up dotclaude repository..."

    if (Test-Path $REPO_DIR) {
        Write-Info "Repository already exists: $REPO_DIR"
        return
    }

    # Clone repository for base configs and profiles
    git clone "https://github.com/$GITHUB_REPO.git" $REPO_DIR
    Write-Info "Cloned to: $REPO_DIR"
}

function Setup-ClaudeDir {
    Write-Step "Setting up Claude directory..."

    if (-not (Test-Path $CLAUDE_DIR)) {
        New-Item -ItemType Directory -Path $CLAUDE_DIR -Force | Out-Null
        Write-Info "Created: $CLAUDE_DIR"
    }

    # Initialize hooks directory
    $hooksDir = "$CLAUDE_DIR\hooks"
    $hookTypes = @("session-start", "post-tool-bash", "post-tool-edit", "pre-tool-bash", "pre-tool-edit")

    foreach ($hookType in $hookTypes) {
        $dir = "$hooksDir\$hookType"
        if (-not (Test-Path $dir)) {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
        }
    }
    Write-Info "Initialized hooks directory"
}

function Add-ToPath {
    Write-Step "Adding to PATH..."

    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")

    if ($currentPath -notlike "*$INSTALL_DIR*") {
        $newPath = "$currentPath;$INSTALL_DIR"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Info "Added $INSTALL_DIR to user PATH"
        Write-Warning "Please restart your terminal for PATH changes to take effect"
    }
    else {
        Write-Info "$INSTALL_DIR already in PATH"
    }
}

function Set-EnvironmentVariable {
    Write-Step "Setting environment variables..."

    [Environment]::SetEnvironmentVariable("DOTCLAUDE_REPO_DIR", $REPO_DIR, "User")
    Write-Info "Set DOTCLAUDE_REPO_DIR=$REPO_DIR"
}

function Test-Installation {
    Write-Step "Verifying installation..."

    # Refresh PATH for current session
    $env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [Environment]::GetEnvironmentVariable("Path", "User")
    $env:DOTCLAUDE_REPO_DIR = $REPO_DIR

    $dotclaude = "$INSTALL_DIR\dotclaude.exe"
    if (Test-Path $dotclaude) {
        $version = & $dotclaude version 2>&1
        Write-Info "dotclaude installed successfully"
        Write-Info "Version: $version"
        return $true
    }

    Write-Error "Installation verification failed"
    return $false
}

function Show-NextSteps {
    Write-Host ""
    Write-Host "  +-------------------------------------------------------------+" -ForegroundColor Green
    Write-Host "  |  Installation Complete!                                     |" -ForegroundColor Green
    Write-Host "  +-------------------------------------------------------------+" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Next Steps:" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  1. Restart your terminal (for PATH changes)" -ForegroundColor White
    Write-Host ""
    Write-Host "  2. Create your first profile:" -ForegroundColor White
    Write-Host "     dotclaude create my-project" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  3. Edit the profile:" -ForegroundColor White
    Write-Host "     dotclaude edit my-project" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  4. Activate it:" -ForegroundColor White
    Write-Host "     dotclaude activate my-project" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  Documentation: https://blackwell-systems.github.io/dotclaude/" -ForegroundColor Gray
    Write-Host ""
}

# Main installation flow
Write-Header

try {
    Install-Binary
    Setup-Repository
    Setup-ClaudeDir
    Add-ToPath
    Set-EnvironmentVariable

    if (Test-Installation) {
        Show-NextSteps
    }
}
catch {
    Write-Error "Installation failed: $_"
    exit 1
}
