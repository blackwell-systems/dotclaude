# Getting Started with dotclaude

Complete guide to installing and using dotclaude for the first time.

## What is dotclaude?

**dotclaude** manages your Claude Code configuration as layered, version-controlled profiles - similar to [dotfiles](https://blackwell-systems.github.io/dotfiles/#/) but specifically for `~/.claude/`.

**The Problem:** You work in multiple contexts (OSS projects, client work, employer projects) that need different tech stacks, coding standards, and compliance requirements. Manually editing `~/.claude/CLAUDE.md` for each context is tedious and error-prone.

**The Solution:** Define universal practices once in a **base** configuration. Create **profiles** that add context-specific details. Switch between them with one command.

## How It Works: The Merge

dotclaude uses two types of configuration:

- **base/** - Universal practices that apply to ALL your work (git workflow, security, tool usage)
- **profiles/** - Context-specific additions per project (tech stack, team standards, compliance)

When you activate a profile, they merge:

```
base/CLAUDE.md              profiles/my-project/CLAUDE.md
(universal standards)   +   (project-specific additions)

• Git workflow                • Tech stack (Node.js, React)
• Security practices          • API design patterns
• Tool preferences            • Team coding standards
• Task management             • Deployment process
                        ↓
              ~/.claude/CLAUDE.md
              (merged configuration)
```

**Key insight:** Write universal practices once in base. Profiles stay small and focused on what makes each project different. No duplication across profiles.

## Quick Start

### Step-by-Step First-Time Setup

```bash
# 1. Install
git clone https://github.com/blackwell-systems/dotclaude.git ~/code/dotclaude
cd ~/code/dotclaude
./install.sh
# → Installs dotclaude CLI, copies base to ~/.claude/

# 2. Create your first profile from the example
cp -r examples/sample-profile profiles/my-project
dotclaude edit my-project
# → Customize tech stack, coding standards, etc.

# 3. Activate it
dotclaude activate my-project
# → Merges base + my-project → ~/.claude/CLAUDE.md

# 4. Verify
dotclaude show              # See what's active
cat ~/.claude/CLAUDE.md     # View merged config

# 5. Create more profiles as needed
cp -r examples/sample-profile profiles/client-work
dotclaude edit client-work
dotclaude activate client-work
```

### Common Commands

```bash
dotclaude list              # List all profiles
dotclaude switch            # Interactive switcher
dotclaude diff <p1> <p2>    # Compare profiles
dotclaude help              # Full command reference
```

## Example Workflow: Multiple Contexts

```bash
# Create your profiles based on the example
cp -r examples/sample-profile profiles/my-oss-project
cp -r examples/sample-profile profiles/client-work

# Customize for each context
vim profiles/my-oss-project/CLAUDE.md
vim profiles/client-work/CLAUDE.md

# Switch between contexts
dotclaude activate my-oss-project
# → Merges base + my-oss-project into ~/.claude/

dotclaude activate client-work
# → Merges base + client-work into ~/.claude/
```

## Repository Structure

```
dotclaude/
├── base/                      # Universal standards (applies to ALL profiles)
│   ├── CLAUDE.md             # Git, security, tools, task management
│   ├── settings.json         # Base hooks & settings
│   ├── scripts/              # Management tools
│   └── agents/               # Shared agents
│
├── examples/                  # Example profile templates
│   └── sample-profile/       # Detailed example showing the overlay pattern
│       ├── CLAUDE.md         # Project-specific additions (tech stack, etc.)
│       ├── settings.json     # Optional profile settings
│       └── README.md         # Explains merge concept with diagrams
│
└── profiles/                  # Your custom profiles (you create these)
    └── (empty - copy from examples/ to start)
```

**When you activate a profile:**
```
~/.claude/CLAUDE.md = base/CLAUDE.md + profiles/my-project/CLAUDE.md
                      (universal)      (project-specific)
```

## Try Before Installing (Docker)

Don't trust random install scripts? Good instinct. Test dotclaude in an isolated container first - nothing touches your system.

### Quick Docker Test

```bash
git clone https://github.com/blackwell-systems/dotclaude.git
cd dotclaude
docker build -t dotclaude-test .
docker run -it --rm dotclaude-test
```

You're now in a container with dotclaude ready to use:

```bash
# Try the commands
dotclaude help
dotclaude create test-project
dotclaude activate test-project
dotclaude show
dotclaude active

# Check generated files
cat ~/.claude/CLAUDE.md
cat ~/.claude/profiles.json
```

Exit with `exit` or Ctrl+D. The container is destroyed - nothing persisted.

### What the Container Includes

- Alpine Linux (minimal ~15MB base)
- Bash, jq, coreutils, util-linux
- dotclaude CLI pre-installed
- Base configuration ready
- Empty profiles directory to experiment with

### Mount Your Own Profiles

Test with your actual profile files without installing:

```bash
docker run -it --rm \
  -v $(pwd)/profiles:/dotclaude/profiles \
  dotclaude-test
```

Changes to profiles persist on your host, but `~/.claude/` stays in the container.

---

## Installation Options

### Platform Requirements

- **Linux** - Fully supported
- **macOS** - Fully supported
- **Windows** - Requires [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install) (Windows Subsystem for Linux 2)

> **Windows Users:** dotclaude requires a Unix/Linux environment with Bash. Native Windows (CMD/PowerShell) is not supported. Install WSL2 first, then run the installation commands from your WSL2 terminal.

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.sh | bash
```

### Manual Installation

```bash
git clone https://github.com/blackwell-systems/dotclaude.git ~/code/dotclaude
cd ~/code/dotclaude

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

### What the Installer Does

1. Installs `dotclaude` CLI to `~/.local/bin/`
2. Copies base scripts and agents to `~/.claude/`
3. Prompts to select and activate a profile

## Core Commands Reference

### Profile Management

```bash
dotclaude show              # Show current profile
dotclaude list              # List all profiles
dotclaude activate <name>   # Activate a profile (add --dry-run to preview)
dotclaude switch            # Interactive switcher
dotclaude diff <p1> [p2]    # Compare profiles or current vs profile
dotclaude create <name>     # Create new profile
dotclaude edit [name]       # Edit profile
dotclaude restore           # Restore from backup
```

### Git Workflow

```bash
dotclaude sync              # Sync feature branch
dotclaude branches          # Check branch status
```

### System

```bash
dotclaude version           # Show version
dotclaude help              # Show help
```

### Debug Mode

Add `--verbose` to any command for detailed debug output:
```bash
dotclaude activate my-project --verbose
dotclaude list --verbose
```

Or set environment variable:
```bash
DEBUG=1 dotclaude list
```

Debug log location: `~/.claude/.dotclaude-debug.log`

## Next Steps

- **[Commands Reference](COMMANDS.md)** - Complete guide to all 10 commands
- **[Complete Usage Guide](USAGE.md)** - Advanced workflows and features
- **[FAQ](FAQ.md)** - Common questions and answers
- **[Architecture](ARCHITECTURE.md)** - How dotclaude works under the hood
- **[Sample Profile](https://github.com/blackwell-systems/dotclaude/blob/main/examples/sample-profile/README.md)** - Learn how to create profiles

---

**Need help?** Visit [GitHub Issues](https://github.com/blackwell-systems/dotclaude/issues)
