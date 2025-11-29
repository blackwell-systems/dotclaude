# dotclaude

**The definitive profile management system for Claude Code**

Version-controlled profile system for managing `~/.claude/` configurations across different work contexts.

> **Disclaimer:** dotclaude is an independent, open-source tool and is not affiliated with or endorsed by Anthropic or the Claude product.

## What is dotclaude?

**dotclaude** manages your Claude Code configuration as layered, version-controlled profiles - similar to dotfiles but specifically for `~/.claude/`.

**The Problem:** You work in multiple contexts (OSS projects, proprietary business, employer work) that need different standards, practices, and tooling.

**The Solution:** Base configuration + profile overlays that merge on activation.

## Features

- **Layered profiles**: Base configuration + context-specific overlays
- **One-command switching**: Switch between work contexts instantly
- **Auto-detection**: `.dotclaude` file detects profile mismatches automatically
- **Preview mode**: Dry-run activation to see changes before applying
- **Profile comparison**: Diff profiles to see differences
- **Backup & restore**: Automatic backups with interactive restoration
- **Version controlled**: Track all configs in git, sync across machines
- **Multi-provider**: Works with both AWS Bedrock and Claude Max
- **Git workflow automation**: Long-lived feature branch sync tools
- **Security hardened**: Input validation, symlink protection, file locking
- **Shell compatible**: Bash scripts with zsh-compatible sourced functions

## Quick Start

```bash
# Clone and install
git clone https://github.com/yourusername/dotclaude.git ~/code/dotclaude
cd ~/code/dotclaude
./install.sh

# Basic commands
dotclaude show              # Show current profile
dotclaude list              # List available profiles
dotclaude switch            # Interactive profile switcher
dotclaude activate <name>   # Activate specific profile
dotclaude help              # Show all commands
```

## Example Workflow

```bash
# Working on OSS project
dotclaude activate blackwell-systems-oss
# → Loads OSS licensing, public docs standards

# Switch to proprietary work
dotclaude activate blackwell-systems
# → Loads internal policies, private repo standards

# Switch to employer work
dotclaude activate best-western
# → Loads corporate compliance, employer guidelines
```

## Repository Structure

```
dotclaude/
├── base/                   # Shared configuration for all profiles
│   ├── CLAUDE.md          # Base development standards
│   ├── settings.json      # Base hooks & settings
│   ├── scripts/           # Management tools
│   └── agents/            # Shared agents
│
└── profiles/              # Context-specific additions
    ├── blackwell-systems-oss/
    ├── blackwell-systems/
    └── best-western/
```

When you activate a profile, base + profile merge into `~/.claude/`.

## Documentation

- **[docs/USAGE.md](docs/USAGE.md)** - Complete user guide and command reference
- **[docs/DOTCLAUDE-FILE.md](docs/DOTCLAUDE-FILE.md)** - `.dotclaude` file format and auto-detection
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - System architecture and technical design
- **[docs/IDEMPOTENCY-AUDIT.md](docs/IDEMPOTENCY-AUDIT.md)** - Idempotency analysis
- **[docs/SECURITY-AUDIT.md](docs/SECURITY-AUDIT.md)** - Security hardening details

## Installation

```bash
# Basic installation
./install.sh

# Force overwrite existing files
./install.sh --force

# Non-interactive mode (for automation)
./install.sh --non-interactive

# Add to shell (optional)
echo 'export DOTCLAUDE_REPO_DIR="$HOME/code/dotclaude"' >> ~/.zshrc
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
```

The installer:
1. Installs `dotclaude` CLI to `~/.local/bin/`
2. Copies base scripts and agents to `~/.claude/`
3. Prompts to select and activate a profile

## Core Commands

```bash
# Profile management
dotclaude show              # Show current profile
dotclaude list              # List all profiles
dotclaude activate <name>   # Activate a profile (add --dry-run to preview)
dotclaude switch            # Interactive switcher
dotclaude diff <p1> [p2]    # Compare profiles or current vs profile
dotclaude create <name>     # Create new profile
dotclaude edit [name]       # Edit profile
dotclaude restore           # Restore from backup

# Git workflow
dotclaude sync              # Sync feature branch
dotclaude branches          # Check branch status

# System
dotclaude version           # Show version
dotclaude help              # Show help
```

## Multi-Provider Support

Works with both:
- **AWS Bedrock**: `us.anthropic.claude-sonnet-4-5-20250929-v1:0`
- **Claude Max**: `claude-sonnet-4-5-20250929`

Global configs are provider-agnostic. Projects specify models via `.claude/settings.json`.

## Shell Compatibility

- **Main CLI**: Uses `#!/bin/bash` shebang, works in any shell environment
- **Sourced functions**: POSIX-compatible, tested with bash and zsh
- **Recommended**: Add to `~/.bashrc` or `~/.zshrc` for convenience functions

## Maintenance

1. Edit configs in this repo
2. Test changes in a project
3. Commit to version control
4. Run `./install.sh` to redeploy
5. Share repo across machines via git

## Security

dotclaude includes comprehensive defensive programming:
- Input validation (prevent path traversal)
- Symlink attack prevention
- File locking (prevent concurrent execution)
- Secure backup permissions (chmod 600)
- Command injection prevention

See [docs/SECURITY-AUDIT.md](docs/SECURITY-AUDIT.md) for details.

## License

MIT

---

**Read next:**
- [docs/USAGE.md](docs/USAGE.md) - Complete usage guide
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) - How it works under the hood
