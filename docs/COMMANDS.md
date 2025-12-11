# Command Reference

Complete reference for all dotclaude commands.

## Command Overview

dotclaude provides commands organized into four categories:

| Category | Commands | Purpose |
|----------|----------|---------|
| **Profile Management** | show, active, list, activate, switch, create, edit, diff, restore | Manage and switch between profiles |
| **Git Workflow** | sync, branches | Keep feature branches in sync |
| **Hooks** | hook run, hook list, hook init | Automation and custom hooks |
| **System** | version, help | Version info and help |
| **Debug** | --verbose flag | Troubleshooting |

---

## Profile Management Commands

### `dotclaude show`

Show current active profile and configuration status.

**Usage:**
```bash
dotclaude show [--verbose]

# Flag aliases
--verbose  (or --debug)
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Active Profile: my-project

  Configuration:
    • CLAUDE.md: 245 lines
    • settings.json: configured

  Location:
    • ~/.claude/

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Run 'dotclaude switch' to change profiles         │
╰─────────────────────────────────────────────────────────────╯
```

**When to use:**
- Check which profile is currently active
- Verify configuration status
- See where files are deployed

---

### `dotclaude active`

Print the active profile name (machine-readable, for scripting).

**Usage:**
```bash
dotclaude active
```

**Output:**
```
my-project
```

Or if no profile is active:
```
none
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Profile is active |
| 1 | No profile active (outputs "none") |

**When to use:**
- Shell scripts that need to check/use the active profile
- Integration with other tools (e.g., blackdot vault sync)
- Conditional logic based on current profile
- Capturing profile name in a variable

**Examples:**
```bash
# Store in variable
current=$(dotclaude active)

# Conditional logic
if [ "$(dotclaude active)" = "work" ]; then
    echo "Work profile active"
fi

# Integration with other tools
profile=$(dotclaude active)
echo "Syncing profile: $profile"
```

**Note:** This differs from `dotclaude show`:
- `show` provides human-readable formatted output with headers
- `active` provides machine-readable single-line output for scripting

---

### `dotclaude list`

List all available profiles.

**Usage:**
```bash
dotclaude list [--verbose]

# Command aliases
dotclaude ls

# Flag aliases
--verbose  (or --debug)
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Available Profiles:

    ● my-project (active)
    ○ client-work
    ○ work-project

  Total: 3 profiles

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Use 'dotclaude activate <name>' or 'dotclaude switch'│
╰─────────────────────────────────────────────────────────────╯
```

**When to use:**
- Discover available profiles
- See which profile is active
- Check total number of profiles

---

### `dotclaude activate`

Activate a specific profile, merging base + profile into `~/.claude/`.

**Usage:**
```bash
dotclaude activate <profile-name> [--dry-run] [--verbose]

# Command aliases
dotclaude use <profile-name>

# Flag aliases
--dry-run  (or --preview)
--verbose  (or --debug)
```

**Examples:**
```bash
# Basic activation
dotclaude activate my-project

# Preview before activating (recommended for first time)
dotclaude activate my-project --dry-run

# Activate with debug output
dotclaude activate my-project --verbose

# Combine flags
dotclaude activate my-project --dry-run --verbose
```

**What it does:**
1. Backs up existing `~/.claude/CLAUDE.md`
2. Merges `base/CLAUDE.md` + `profiles/<name>/CLAUDE.md`
3. Writes merged result to `~/.claude/CLAUDE.md`
4. Applies profile-specific `settings.json` (if present)
5. Updates `~/.claude/.current-profile` marker

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Activating profile: my-project

  [1/3] Backed up existing CLAUDE.md
  [2/3] Merged base + profile configuration
  [3/3] Applied profile settings

╭─────────────────────────────────────────────────────────────╮
│  ✓ Profile 'my-project' activated               │
╰─────────────────────────────────────────────────────────────╯

  Configuration deployed to: /home/user/.claude

  Verify with:
    • dotclaude show
    • cat ~/.claude/CLAUDE.md

╭─────────────────────────────────────────────────────────────╮
│  🍃 Happy coding!                                       │
╰─────────────────────────────────────────────────────────────╯
```

**Dry-Run Mode:**

Preview changes without applying them:

```bash
dotclaude activate my-project --dry-run
```

Shows what would happen:
- Which files would be backed up
- What would be merged (with line counts)
- Which files would be modified
- No actual changes made

**When to use:**
- Switching between work contexts
- After editing base or profile
- Setting up a new project

---

### `dotclaude switch`

Interactive profile switcher with menu selection.

**Usage:**
```bash
dotclaude switch [--verbose]

# Command aliases
dotclaude select

# Flag aliases
--verbose  (or --debug)
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Select a profile to activate:

    [1] my-project (active)
    [2] client-work
    [3] work-project

  Enter number (or 'q' to quit): 2

  Activating profile: client-work
  ...
```

**When to use:**
- Quick switching without typing profile names
- When you can't remember exact profile names
- Visual confirmation of available profiles

---

### `dotclaude create`

Create a new empty profile.

**Usage:**
```bash
dotclaude create <profile-name> [--verbose]

# Command aliases
dotclaude new <profile-name>

# Flag aliases
--verbose  (or --debug)
```

**Example:**
```bash
dotclaude create my-new-project
```

**What it does:**
1. Creates `profiles/<name>/` directory
2. Creates basic `CLAUDE.md` template
3. Profile is ready to edit and activate

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Creating new profile: my-new-project

╭─────────────────────────────────────────────────────────────╮
│  ✓ Profile 'my-new-project' created                        │
╰─────────────────────────────────────────────────────────────╯

  Location: /home/user/code/dotclaude/profiles/my-new-project

  Next steps:
    • dotclaude edit my-new-project
    • dotclaude activate my-new-project

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Edit the profile to add your guidelines           │
╰─────────────────────────────────────────────────────────────╯
```

**Note:** Creating from the example is usually better:
```bash
cp -r examples/sample-profile profiles/my-new-project
```

---

### `dotclaude edit`

Edit a profile's CLAUDE.md in your configured editor.

**Usage:**
```bash
dotclaude edit [profile-name] [--verbose]

# Flag aliases
--verbose  (or --debug)
```

**Examples:**
```bash
# Edit current active profile
dotclaude edit

# Edit specific profile
dotclaude edit my-project
```

**What it does:**
- Opens `profiles/<name>/CLAUDE.md` in `$EDITOR` (or nano)
- If no profile specified, edits currently active profile
- After editing, run `dotclaude activate <name>` to apply changes

**When to use:**
- Updating project guidelines
- Adding new tech stack preferences
- Modifying coding standards

---

### `dotclaude diff`

Compare two profiles or compare current profile with another.

**Usage:**
```bash
dotclaude diff <profile1> <profile2> [--verbose]
dotclaude diff <profile> [--verbose]  # Compare current vs profile

# Flag aliases
--verbose  (or --debug)
```

**Examples:**
```bash
# Compare two specific profiles
dotclaude diff my-project client-work

# Compare current active profile with another
dotclaude diff work-project
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Comparing profiles: my-project vs client-work

  CLAUDE.md differences:

    Differences found:

      @@ -1,5 +1,5 @@
      -# Profile: my-project
      +# Profile: client-work

      -Open source best practices
      +Proprietary code handling

      ... (first 50 lines shown)

    Tip: See full diff with:
      diff -u profiles/my-project/CLAUDE.md profiles/client-work/CLAUDE.md

  settings.json:

    ✓ No differences

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Use 'dotclaude activate <profile> --dry-run' to preview│
╰─────────────────────────────────────────────────────────────╯
```

**When to use:**
- See differences before switching
- Understand how profiles differ
- Verify profile-specific settings

---

### `dotclaude restore`

Restore from backup interactively.

**Usage:**
```bash
dotclaude restore [--verbose]

# Flag aliases
--verbose  (or --debug)
```

**What it does:**
1. Lists all available backups (CLAUDE.md and settings.json)
2. Interactive selection
3. Backs up current file before restoring
4. Restores selected backup
5. Updates active profile marker

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Backup Restoration

  Available backups:

  CLAUDE.md backups:
    [1] 20241129-143022 (24K)
    [2] 20241129-120145 (22K)
    [3] 20241128-183045 (23K)

  settings.json backups:
    [4] 20241129-143022 (4.2K)
    [5] 20241129-120145 (3.8K)

  Select backup to restore (or 'q' to quit): 1

  ⚠  This will overwrite:
    /home/user/.claude/CLAUDE.md

  Continue? (y/N): y

  [BACKUP] Current file backed up
  [RESTORE] Restored from backup
  [UPDATE] Active profile marker updated

╭─────────────────────────────────────────────────────────────╮
│  ✓ Backup restored successfully                            │
╰─────────────────────────────────────────────────────────────╯
```

**When to use:**
- Undo a profile activation
- Recover from accidental changes
- Roll back to previous configuration
- After testing a new profile

**Safety:**
- Always backs up current file before restoring
- Keeps 5 most recent backups automatically
- Requires confirmation before overwriting

---

## Git Workflow Commands

### `dotclaude sync`

Run the interactive feature branch sync tool.

**Usage:**
```bash
dotclaude sync [--verbose]

# Flag aliases
--verbose  (or --debug)
```

**What it does:**
Helps keep long-lived feature branches in sync with main:
1. Checks if current branch is behind main
2. Offers choice: rebase or merge
3. Handles sync automatically
4. Provides conflict resolution guidance

**When to use:**
- After your PR is merged to main
- Before continuing work on feature branch
- When branch is behind main

**See also:** [Long-Lived Feature Branch Management](USAGE.md#long-lived-feature-branch-management)

---

### `dotclaude branches`

Check status of all branches against main.

**Usage:**
```bash
dotclaude branches [--verbose]

# Command aliases
dotclaude br

# Flag aliases
--verbose  (or --debug)
```

**Output:**
```
Checking branches against main...

  feature-add-auth              10 ahead, 5 behind
  feature-refactor              2 ahead, 12 behind
  feature-experimental          0 ahead, 3 behind
```

**When to use:**
- See which branches need syncing
- Check branch status across your repo
- Identify branches that are behind main

---

## Hook Commands

### `dotclaude hook run`

Execute all hooks of a given type.

**Usage:**
```bash
dotclaude hook run <hook-type>
```

**Hook Types:**
| Type | Trigger | Description |
|------|---------|-------------|
| `session-start` | Session starts | Session info, profile checks |
| `post-tool-bash` | After Bash tool | Git workflow tips |
| `post-tool-edit` | After Edit tool | Custom post-edit hooks |
| `pre-tool-bash` | Before Bash tool | Pre-command validation |
| `pre-tool-edit` | Before Edit tool | Pre-edit validation |

**Examples:**
```bash
# Run session start hooks
dotclaude hook run session-start

# Run post-bash hooks
dotclaude hook run post-tool-bash
```

**When to use:**
- Called automatically by Claude Code via settings.json hooks
- Manually for testing custom hooks

---

### `dotclaude hook list`

List all available hooks for a given type.

**Usage:**
```bash
dotclaude hook list [hook-type]
```

**Examples:**
```bash
# List all hooks
dotclaude hook list

# List session-start hooks only
dotclaude hook list session-start
```

**Output:**
```
session-start:
  [00] session-info (built-in, enabled)
  [10] check-dotclaude (built-in, enabled)
  [20] custom-greeting.sh (enabled)

post-tool-bash:
  [10] git-tips (built-in, enabled)
```

---

### `dotclaude hook init`

Initialize the hooks directory structure.

**Usage:**
```bash
dotclaude hook init
```

**What it does:**
Creates `~/.claude/hooks/` with subdirectories for each hook type.

**When to use:**
- First-time setup
- After fresh installation
- To add custom hooks

**See also:** [HOOKS.md](HOOKS.md) for complete hook system documentation.

---

## System Commands

### `dotclaude version`

Show dotclaude version information.

**Usage:**
```bash
dotclaude version

# Aliases
dotclaude -v
dotclaude --version
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  dotclaude version 0.2.0

  The definitive profile management system for Claude Code

  Repository: /home/user/code/dotclaude
  Configuration: /home/user/.claude
```

---

### `dotclaude help`

Show help information.

**Usage:**
```bash
# General help
dotclaude help
dotclaude -h
dotclaude --help

# Command-specific help (coming soon)
dotclaude help activate
```

**Output:**
Shows complete command reference with usage examples.

---

## Debug Mode

All commands support the `--verbose` (or `--debug`) flag for detailed debug output.

### Usage

**Add --verbose (or --debug) to any command:**
```bash
dotclaude show --verbose
dotclaude list --debug          # --debug is an alias for --verbose
dotclaude activate my-project --verbose
dotclaude activate my-project --dry-run --debug
dotclaude activate my-project --preview --verbose  # --preview is alias for --dry-run
dotclaude switch --verbose
dotclaude create new-profile --debug
dotclaude edit my-project --verbose
dotclaude diff profile1 profile2 --debug
dotclaude restore --verbose
dotclaude sync --debug
dotclaude branches --verbose
```

**Or use environment variable:**
```bash
# Single command
DEBUG=1 dotclaude list

# Entire session
export DEBUG=1
dotclaude activate my-project
dotclaude switch
```

### Debug Output Includes

- Command parsing and arguments
- File paths being used
- Validation checks
- Lock acquisition
- Profile discovery
- Merge operations
- All file operations

### Debug Log File

**Location:** `~/.claude/.dotclaude-debug.log`

**Viewing:**
```bash
# View entire log
cat ~/.claude/.dotclaude-debug.log

# Watch in real-time
tail -f ~/.claude/.dotclaude-debug.log

# View recent entries
tail -20 ~/.claude/.dotclaude-debug.log
```

**Benefits:**
- Persists across sessions
- Useful for troubleshooting
- Share when reporting issues
- See exact operation sequence

---

## Command Aliases

Many commands have shorter aliases for convenience:

| Full Command | Aliases |
|--------------|---------|
| `dotclaude list` | `dotclaude ls` |
| `dotclaude activate <name>` | `dotclaude use <name>` |
| `dotclaude switch` | `dotclaude select` |
| `dotclaude create <name>` | `dotclaude new <name>` |
| `dotclaude branches` | `dotclaude br` |
| `dotclaude version` | `dotclaude -v`, `dotclaude --version` |
| `dotclaude help` | `dotclaude -h`, `dotclaude --help` |

---

## Common Workflows

### Daily Context Switching

```bash
# Morning: Start OSS work
dotclaude activate my-oss-project

# Afternoon: Switch to client work
dotclaude activate client-work

# Check current context anytime
dotclaude show
```

### Creating and Testing New Profile

```bash
# Create from example
cp -r examples/sample-profile profiles/test-project

# Edit it
dotclaude edit test-project

# Preview what would change
dotclaude activate test-project --dry-run

# Apply it
dotclaude activate test-project

# Verify
dotclaude show
cat ~/.claude/CLAUDE.md
```

### Comparing Profiles

```bash
# Compare two profiles
dotclaude diff oss-project client-work

# Compare current with another
dotclaude diff new-project

# Use full diff command for details
diff -u profiles/oss-project/CLAUDE.md profiles/client-work/CLAUDE.md
```

### Backup and Recovery

```bash
# Profiles auto-backup on activation
dotclaude activate new-profile
# → Previous config backed up automatically

# Made a mistake? Restore previous
dotclaude restore
# → Interactive selection from backups

# Or just switch back
dotclaude activate previous-profile
```

---

## Advanced Usage

### Environment Variables

**DOTCLAUDE_REPO_DIR**
- Override default repository location
- **Automatically set** by installer in shell RC file
- Default: `$HOME/code/dotclaude`
- Only change if you moved the repo: `export DOTCLAUDE_REPO_DIR="/new/path"`

**DEBUG**
- Enable debug output
- Values: `0` (off) or `1` (on)
- Usage: `DEBUG=1 dotclaude list`

**EDITOR**
- Choose editor for `dotclaude edit`
- Default: `$EDITOR` or `nano`
- Usage: `EDITOR=vim dotclaude edit my-project`

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (invalid input, profile not found, etc.) |

---

## Next Steps

- **[Getting Started Guide](GETTING-STARTED.md)** - Complete setup tutorial
- **[Usage Guide](USAGE.md)** - Comprehensive workflows and patterns
- **[FAQ](FAQ.md)** - Common questions answered
- **[Architecture](ARCHITECTURE.md)** - How it works under the hood
