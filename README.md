# dotclaude

[![Blackwell Systems™](https://raw.githubusercontent.com/blackwell-systems/blackwell-docs-theme/main/badge-trademark.svg)](https://github.com/blackwell-systems)
[![Claude Code](https://img.shields.io/badge/Built_for-Claude_Code-8A2BE2?logo=anthropic)](https://claude.ai/claude-code)
[![dotfiles](https://img.shields.io/badge/Integrates-dotfiles-2c5282)](https://blackwell-systems.github.io/dotfiles/#/)
[![Shell](https://img.shields.io/badge/Shell-Bash-4EAA25?logo=gnu-bash&logoColor=white)](https://www.gnu.org/software/bash/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20WSL2-blue)](https://github.com/blackwell-systems/dotclaude)
[![Version](https://img.shields.io/badge/Version-0.3.0-informational)](https://github.com/blackwell-systems/dotclaude/releases)

[![Tests](https://github.com/blackwell-systems/dotclaude/actions/workflows/test.yml/badge.svg)](https://github.com/blackwell-systems/dotclaude/actions/workflows/test.yml)
[![Tests](https://img.shields.io/badge/Tests-122_passing-success)](https://github.com/blackwell-systems/dotclaude/tree/main/tests)
[![Sponsor](https://img.shields.io/badge/Sponsor-Buy%20Me%20a%20Coffee-yellow?logo=buy-me-a-coffee&logoColor=white)](https://buymeacoffee.com/blackwellsystems)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**The definitive profile management system for Claude Code**

Manage your Claude Code configuration as layered, version-controlled profiles. Switch between work contexts (OSS, client, employer) with one command.

> **Disclaimer:** dotclaude is an independent, open-source tool and is not affiliated with or endorsed by Anthropic or the Claude product.

## What is dotclaude?

Stop manually editing `~/.claude/CLAUDE.md` every time you switch projects.

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
- **Preview mode** (dry-run before applying)
- **Backup & restore** with automatic versioning
- **Version controlled** (sync across machines via git)
- **Git workflow tools** for long-lived feature branches
- **Multi-provider** (AWS Bedrock + Claude Max)

## Platform Support

- **Linux** - Full support
- **macOS** - Full support
- **Windows (WSL2)** - Full support via Windows Subsystem for Linux 2

> **Windows Users:** dotclaude requires a Unix environment. Install [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install) to run dotclaude on Windows. Native Windows (CMD/PowerShell) is not supported.

## Try Before Installing

Don't trust random scripts from the internet? Smart. Test dotclaude in an isolated Docker container first:

```bash
git clone https://github.com/blackwell-systems/dotclaude.git
cd dotclaude
docker build -t dotclaude-test .
docker run -it --rm dotclaude-test

# Inside container - explore safely:
dotclaude create my-project
dotclaude activate my-project
dotclaude show
dotclaude active
exit  # Nothing persists
```

**→ [Full Test Drive Guide](docs/TESTDRIVE.md)** - Sample workflows, all commands explained, FAQ

Nothing touches your system. When you're ready, install for real.

## Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.sh | bash
```

Then create your first profile:

```bash
cp -r examples/sample-profile profiles/my-project
dotclaude edit my-project
dotclaude activate my-project
```

## Documentation

**📚 [Complete Documentation Site](https://blackwell-systems.github.io/dotclaude/)**

- **[Getting Started](https://blackwell-systems.github.io/dotclaude/#/GETTING-STARTED)** - Installation, concepts, and first steps
- **[Commands Reference](https://blackwell-systems.github.io/dotclaude/#/COMMANDS)** - Complete command guide
- **[dotfiles Integration](https://blackwell-systems.github.io/dotclaude/#/DOTFILES-INTEGRATION)** - Use with dotfiles for complete environment
- **[FAQ](https://blackwell-systems.github.io/dotclaude/#/FAQ)** - Common questions answered
- **[Usage Guide](https://blackwell-systems.github.io/dotclaude/#/USAGE)** - Advanced workflows and features
- **[Architecture](https://blackwell-systems.github.io/dotclaude/#/ARCHITECTURE)** - How it works under the hood

## Trademarks

**Blackwell Systems™** and the **Blackwell Systems logo** are trademarks of Dayna Blackwell. You may use the name "Blackwell Systems" to refer to this project, but you may not use the name or logo in a way that suggests endorsement or official affiliation without prior written permission. See [BRAND.md](docs/BRAND.md) for usage guidelines.

## License

MIT

---

**Questions?** Check the [FAQ](https://blackwell-systems.github.io/dotclaude/#/FAQ) or [open an issue](https://github.com/blackwell-systems/dotclaude/issues)
