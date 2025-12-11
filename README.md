# dotclaude

[![Blackwell Systems™](https://raw.githubusercontent.com/blackwell-systems/blackwell-docs-theme/main/badge-trademark.svg)](https://github.com/blackwell-systems)
[![Claude Code](https://img.shields.io/badge/Built_for-Claude_Code-8A2BE2?logo=anthropic)](https://claude.ai/claude-code)
[![blackdot](https://img.shields.io/badge/Integrates-blackdot-2c5282)](https://blackwell-systems.github.io/blackdot/#/)
[![Go Reference](https://pkg.go.dev/badge/github.com/blackwell-systems/dotclaude.svg)](https://pkg.go.dev/github.com/blackwell-systems/dotclaude)
[![Go Report Card](https://goreportcard.com/badge/github.com/blackwell-systems/dotclaude)](https://goreportcard.com/report/github.com/blackwell-systems/dotclaude)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/blackwell-systems/dotclaude)
[![GitHub Release](https://img.shields.io/github/v/release/blackwell-systems/dotclaude?include_prereleases)](https://github.com/blackwell-systems/dotclaude/releases)

[![Tests](https://github.com/blackwell-systems/dotclaude/actions/workflows/test.yml/badge.svg)](https://github.com/blackwell-systems/dotclaude/actions/workflows/test.yml)
[![Tests](https://img.shields.io/badge/Tests-122_passing-success)](https://github.com/blackwell-systems/dotclaude/tree/main/tests)
[![Sponsor](https://img.shields.io/badge/Sponsor-Buy%20Me%20a%20Coffee-yellow?logo=buy-me-a-coffee&logoColor=white)](https://buymeacoffee.com/blackwellsystems)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**The definitive profile management system for Claude Code**

Manage your Claude Code configuration as layered, version-controlled profiles. Switch between work contexts (OSS, client, employer) with one command.

> **Disclaimer:** dotclaude is an independent, open-source tool and is not affiliated with or endorsed by Anthropic or the Claude product.

## What is dotclaude?

Stop manually editing `~/.claude/CLAUDE.md` (or `%USERPROFILE%\.claude\CLAUDE.md` on Windows) every time you switch projects.

**dotclaude** lets you define universal practices once in a **base** configuration, then create **profiles** that add context-specific details (tech stack, coding standards, compliance requirements). Switch between them instantly.

```bash
# Morning: OSS work
dotclaude activate my-oss-project

# Afternoon: Client work
dotclaude activate client-work
```

Each profile merges your base configuration with project-specific additions - no duplication across profiles.

## Features

- **One-command switching** between work contexts
- **Layered profiles** (base + context-specific overlays)
- **Auto-detection** via `.dotclaude` file
- **Hook system** for automation (session start, post-tool events)
- **Preview mode** (dry-run before applying)
- **Backup & restore** with automatic versioning
- **Version controlled** (sync across machines via git)
- **Git workflow tools** for long-lived feature branches
- **Cross-platform** (Linux, macOS, Windows native)

## Platform Support

- **Linux** (amd64, arm64) - Full support
- **macOS** (Intel & Apple Silicon) - Full support
- **Windows** (amd64) - Full support (native binary)

## Try Before Installing

Don't trust random scripts from the internet? Smart. Test dotclaude in an isolated Docker container first:

```bash
# Quick test (uses pre-built lite image)
docker run -it --rm ghcr.io/blackwell-systems/dotclaude-lite

# Or build locally
git clone https://github.com/blackwell-systems/dotclaude.git
cd dotclaude
docker build -f Dockerfile.lite -t dotclaude-lite .
docker run -it --rm dotclaude-lite

# Inside container - explore safely:
dotclaude create my-project
dotclaude edit my-project
dotclaude activate my-project
dotclaude show
exit  # Nothing persists
```

**→ [Full Test Drive Guide](docs/TESTDRIVE.md)** - Sample workflows, all commands explained, FAQ

Nothing touches your system. When you're ready, install for real.

## Quick Install

### Option 1: Download Pre-built Binary (No Go Required)

Download the latest release for your platform:

```bash
# macOS (Apple Silicon)
curl -sL https://github.com/blackwell-systems/dotclaude/releases/latest/download/dotclaude_darwin_arm64.tar.gz | tar xz
sudo mv dotclaude /usr/local/bin/

# macOS (Intel)
curl -sL https://github.com/blackwell-systems/dotclaude/releases/latest/download/dotclaude_darwin_amd64.tar.gz | tar xz
sudo mv dotclaude /usr/local/bin/

# Linux (x86_64)
curl -sL https://github.com/blackwell-systems/dotclaude/releases/latest/download/dotclaude_linux_amd64.tar.gz | tar xz
sudo mv dotclaude /usr/local/bin/

# Linux (ARM64)
curl -sL https://github.com/blackwell-systems/dotclaude/releases/latest/download/dotclaude_linux_arm64.tar.gz | tar xz
sudo mv dotclaude /usr/local/bin/

# Windows (PowerShell - as Administrator)
Invoke-WebRequest -Uri "https://github.com/blackwell-systems/dotclaude/releases/latest/download/dotclaude_windows_amd64.zip" -OutFile dotclaude.zip
Expand-Archive dotclaude.zip -DestinationPath "$env:ProgramFiles\dotclaude"
# Add $env:ProgramFiles\dotclaude to your PATH
```

Or browse all releases at the [releases page](https://github.com/blackwell-systems/dotclaude/releases).

### Option 2: Go Install

If you have Go installed:

```bash
go install github.com/blackwell-systems/dotclaude/cmd/dotclaude@latest
```

### Option 3: Install Script

One-line install (clones repository with base configs):

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
iex (iwr -Uri "https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.ps1").Content
```

Then follow the guided setup:

```bash
dotclaude create my-project    # Create your first profile
dotclaude edit my-project      # Add your project context
dotclaude activate my-project  # Activate it
dotclaude show                 # Verify it's active
```

## Documentation

**📚 [Complete Documentation Site](https://blackwell-systems.github.io/dotclaude/)**

- **[Getting Started](https://blackwell-systems.github.io/dotclaude/#/GETTING-STARTED)** - Installation, concepts, and first steps
- **[Commands Reference](https://blackwell-systems.github.io/dotclaude/#/COMMANDS)** - Complete command guide
- **[Hooks Guide](https://blackwell-systems.github.io/dotclaude/#/HOOKS)** - Automation and custom hooks
- **[Usage Guide](https://blackwell-systems.github.io/dotclaude/#/USAGE)** - Advanced workflows and features
- **[blackdot Integration](https://blackwell-systems.github.io/dotclaude/#/BLACKDOT-INTEGRATION)** - Use with blackdot for complete environment
- **[Architecture](https://blackwell-systems.github.io/dotclaude/#/ARCHITECTURE)** - How it works under the hood
- **[FAQ](https://blackwell-systems.github.io/dotclaude/#/FAQ)** - Common questions answered

## Trademarks

**Blackwell Systems™** and the **Blackwell Systems logo** are trademarks of Dayna Blackwell. You may use the name "Blackwell Systems" to refer to this project, but you may not use the name or logo in a way that suggests endorsement or official affiliation without prior written permission. See [BRAND.md](docs/BRAND.md) for usage guidelines.

## License

Apache 2.0 - See [LICENSE](LICENSE) for details.

---

**Questions?** Check the [FAQ](https://blackwell-systems.github.io/dotclaude/#/FAQ) or [open an issue](https://github.com/blackwell-systems/dotclaude/issues)
