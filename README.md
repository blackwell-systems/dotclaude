# dotclaude

**The definitive profile management system for Claude Code**

Version-controlled profile system for managing `~/.claude/` configurations across different work contexts.

> **Disclaimer:** dotclaude is an independent, open-source tool and is not affiliated with or endorsed by Anthropic or the Claude product.

## What is dotclaude?

**dotclaude** manages your Claude Code configuration as layered, version-controlled profiles - similar to dotfiles but specifically for `~/.claude/`.

**The Problem:** You work in multiple contexts (OSS projects, client work, employer projects) that need different tech stacks, coding standards, and compliance requirements. Manually editing `~/.claude/CLAUDE.md` for each context is tedious and error-prone.

**The Solution:** Define universal practices once in a **base** configuration. Create **profiles** that add context-specific details. Switch between them with one command.

### How It Works: The Merge

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

**Common commands:**
```bash
dotclaude list              # List all profiles
dotclaude switch            # Interactive switcher
dotclaude diff <p1> <p2>    # Compare profiles
dotclaude help              # Full command reference
```

## Example Workflow

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

## Documentation

**📚 [View Documentation Site](https://blackwell-systems.github.io/dotclaude/)** (Recommended)

Or browse markdown files directly:
- **[docs/USAGE.md](docs/USAGE.md)** - Complete user guide and command reference
- **[docs/DOTCLAUDE-FILE.md](docs/DOTCLAUDE-FILE.md)** - `.dotclaude` file format and auto-detection
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - System architecture and technical design

### Local Documentation Site

View the documentation site locally:

```bash
# Install docsify-cli globally
npm install -g docsify-cli

# Serve the docs
cd ~/code/dotclaude
docsify serve .

# Open http://localhost:3000 in your browser
```

Or use Python's built-in server:

```bash
cd ~/code/dotclaude
python3 -m http.server 3000
# Open http://localhost:3000 in your browser
```

## Installation

**Quick install (recommended):**

```bash
curl -fsSL https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.sh | bash
```

**Or clone first:**

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

## FAQ

**Q: Do I edit `~/.claude/CLAUDE.md` directly?**
A: No. Edit `base/CLAUDE.md` or `profiles/*/CLAUDE.md`, then run `dotclaude activate` to merge and deploy. The `~/.claude/` files are generated.

**Q: What goes in base vs profiles?**
A: Base = universal practices (git, security, tools). Profiles = project-specific details (tech stack, team standards). If it applies to ALL your work, it's base.

**Q: Can I share profiles with my team?**
A: Yes. Commit your dotclaude repo to git and share. Everyone gets consistent Claude Code behavior.

**Q: What if I only have one project?**
A: You probably don't need dotclaude. Just edit `~/.claude/CLAUDE.md` directly. dotclaude is for managing multiple contexts.

**Q: Do profiles replace base or add to it?**
A: Profiles ADD to base. Both are merged together. This prevents duplication - you write universal practices once in base.

**Q: How do I see what Claude Code is actually reading?**
A: `cat ~/.claude/CLAUDE.md` - this is the merged result after activation.

**Q: Can I use different Claude models per profile?**
A: Yes. Add `settings.json` to your profile with `"model": "opus"` or `"model": "haiku"`.

**Q: How do I troubleshoot issues?**
A: Enable debug mode to see detailed operation logs:
```bash
# Method 1: Use --verbose flag (add to end of command)
dotclaude activate my-project --verbose
dotclaude list --verbose

# Method 2: Set DEBUG environment variable
DEBUG=1 dotclaude activate my-project

# Debug logs are also written to ~/.claude/.dotclaude-debug.log
```

## License

MIT

---

**Read next:**
- [docs/USAGE.md](docs/USAGE.md) - Complete usage guide
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) - How it works under the hood
