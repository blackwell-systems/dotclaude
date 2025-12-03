# Test Drive: Try Before You Trust

> **Don't trust random install scripts?** Smart. Test dotclaude in an isolated Alpine container before touching your system.

---

## Quick Start

```bash
# Clone and build lightweight test container
git clone https://github.com/blackwell-systems/dotclaude.git
cd dotclaude
docker build -f Dockerfile.lite -t dotclaude-lite .

# Run interactive container
docker run -it --rm dotclaude-lite

# You're now in a safe container - nothing touches your host
```

When you're done, type `exit` or press `Ctrl+D`. The container is destroyed instantly.

---

## What to Try

### 1. List Profiles

```bash
dotclaude list
```

Shows all available profiles (starts with the examples from `/dotclaude/profiles/`).

---

### 2. Create a Test Profile

```bash
dotclaude create test-project
```

This creates a new profile using a comprehensive 250+ line template with:
- Tech stack preferences (backend, frontend, testing)
- Coding standards (TypeScript, API design, error handling)
- Project workflows and best practices
- All customizable for your project

**What this shows:** Profile creation with scaffolding, not empty templates.

---

### 3. Switch Profiles

```bash
dotclaude list
dotclaude activate test-project
dotclaude show
```

Verify the profile is now active with `show`.

**What this shows:** Profile switching and `~/.claude/.current-profile` management.

---

### 4. Check Active Profile

```bash
dotclaude active
```

Returns just the profile name - designed for scripting and dotfiles integration.

**What this shows:** Machine-readable output for automation.

---

### 5. View Configuration

```bash
# Check what dotclaude created
cat ~/.claude/CLAUDE.md
cat ~/.claude/settings.json
cat ~/.claude/profiles.json

# See the profile directory structure
ls -la ~/.claude/
tree /dotclaude/profiles/  # if tree is installed
```

**What this shows:** How dotclaude organizes profiles and settings.

---

### 6. Test Backend Configuration

```bash
dotclaude create work-bedrock

# Check the generated settings
cat ~/.claude/settings.json
```

**What this shows:** How profiles configure Claude backends (Max, Bedrock, etc.).

---

### 7. Test profiles.json Generation

```bash
# Create multiple profiles
dotclaude create personal-max
dotclaude create work-bedrock
dotclaude activate personal-max

# Check auto-generated manifest
cat ~/.claude/profiles.json
```

You'll see JSON with:
- Active profile
- List of all profiles with metadata

**What this shows:** The profiles.json that syncs with dotfiles vault.

---

## Sample Workflows

### Workflow 1: Explore Without Installing

```bash
# Just browse the codebase
cd /dotclaude

# Read the main script
cat base/scripts/dotclaude

# Check base configuration
cat base/CLAUDE.md

# View command documentation
cat docs/COMMANDS.md

# Exit when done
exit
```

**Zero impact** - container destroyed, nothing persisted.

---

### Workflow 2: Test All Commands

```bash
dotclaude help
dotclaude list
dotclaude create test-1
dotclaude create test-2
dotclaude list --verbose
dotclaude show
dotclaude activate test-2
dotclaude active
dotclaude show --verbose

exit
```

**Learn the CLI** without touching your system.

---

### Workflow 3: Test Profile Management

```bash
# Create work and personal profiles
dotclaude create work-bedrock
dotclaude create personal-max

# Switch between them
dotclaude activate work-bedrock
dotclaude show

dotclaude switch personal-max
dotclaude show

# See the profiles.json state
cat ~/.claude/profiles.json

exit
```

**Practice profile workflows** in isolation.

---

### Workflow 4: Test with Your Profiles (Advanced)

Mount your actual profiles for testing:

```bash
# From host - mount your profiles directory
docker run -it --rm \
  -v ~/code/dotclaude/profiles:/root/code/dotclaude/profiles \
  dotclaude-lite

# Inside container
dotclaude list
dotclaude show

exit
```

**Your profiles are accessible** but changes to `~/.claude/` stay in container.

---

## What You Can't Test

Some features require Claude Code itself:

| Feature | Why Not in Container | Test How |
|---------|---------------------|----------|
| **Actually use Claude** | No Claude CLI in container | Install on host |
| **Test commands** | Commands live in profile dirs | Mount profile with commands |
| **Multi-backend switching** | Need actual API keys | Use on real system |
| **Git integration** | No git repos | Clone test repo in container |

---

## Container Specs

**Base:** Alpine Linux 3.19 (~29MB)

**Includes:**
- bash, git, jq
- coreutils, util-linux
- dotclaude CLI (pre-installed)
- Base configuration and examples directory with comprehensive template

**Does NOT include:**
- Claude CLI (`claude` command)
- Your actual profiles
- API keys/credentials
- AWS or Anthropic configuration

---

## Next Steps

### Ready to Install?

```bash
# One-line install (clones automatically)
curl -fsSL https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.sh | bash

# Or manual
git clone https://github.com/blackwell-systems/dotclaude.git ~/code/dotclaude
cd ~/code/dotclaude
./install.sh
```

### Want to Learn More?

- [Documentation](../README.md) - Full dotclaude guide
- [Commands Reference](COMMANDS.md) - All commands explained
- [Getting Started](GETTING-STARTED.md) - Setup and usage
- [Integration with dotfiles](https://github.com/blackwell-systems/dotfiles) - Portable sessions

---

## FAQ

**Q: Can I break my host system from the container?**
A: No. With `--rm` flag, nothing persists. Without volume mounts, the container can't touch your files.

**Q: Can I use Claude from inside the container?**
A: Not without mounting your Claude CLI and credentials. The container only tests dotclaude profile management.

**Q: Will my profiles be saved?**
A: No, they're created in the container and destroyed on exit. That's the point - safe exploration.

**Q: How do I test with my actual profiles?**
A: Mount your profiles directory: `-v ~/code/dotclaude/profiles:/dotclaude/profiles`

**Q: Can I test the install script?**
A: Yes, but you'd need a fresh container for that. The Dockerfile already has dotclaude installed.

---

## Troubleshooting

### "docker: command not found"

Install Docker:
- **macOS:** `brew install --cask docker` or [Docker Desktop](https://docker.com)
- **Linux:** `curl -fsSL https://get.docker.com | sh`
- **Windows:** [Docker Desktop for Windows](https://docker.com)

### Build fails

```bash
# Clean rebuild
docker build --no-cache -f Dockerfile.lite -t dotclaude-lite .
```

### Container won't start

```bash
# Check Docker is running
docker ps

# Try with explicit shell
docker run -it --rm dotclaude-lite bash
```

### dotclaude command not found

```bash
# Check it's installed
which dotclaude
dotclaude --version

# If missing, rebuild container
docker build -f Dockerfile.lite -t dotclaude-lite .
```

### Commands don't work

Make sure you're testing dotclaude commands (profile management), not Claude commands (AI usage). Claude CLI isn't installed in the test container.

---

**Ready to trust it?** Install for real: [Installation Guide](GETTING-STARTED.md)

**Still skeptical?** Read the code: [github.com/blackwell-systems/dotclaude](https://github.com/blackwell-systems/dotclaude)
