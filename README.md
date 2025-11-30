# dotclaude

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
- **[FAQ](https://blackwell-systems.github.io/dotclaude/#/FAQ)** - Common questions answered
- **[Usage Guide](https://blackwell-systems.github.io/dotclaude/#/USAGE)** - Advanced workflows and features
- **[Architecture](https://blackwell-systems.github.io/dotclaude/#/ARCHITECTURE)** - How it works under the hood

## License

MIT

---

**Questions?** Check the [FAQ](https://blackwell-systems.github.io/dotclaude/#/FAQ) or [open an issue](https://github.com/blackwell-systems/dotclaude/issues)
