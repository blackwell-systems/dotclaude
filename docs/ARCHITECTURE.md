# dotclaude Architecture

Technical overview of how dotclaude works internally.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         dotclaude Repository                         │
│                      (Version-Controlled Git Repo)                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────────────┐                  ┌────────────────────────┐   │
│  │  base/           │                  │  profiles/             │   │
│  │                  │                  │                        │   │
│  │  • CLAUDE.md     │                  │  • proprietary-project-  │   │
│  │  • settings.json │                  │    oss/CLAUDE.md       │   │
│  │  • scripts/      │                  │  • proprietary-project/  │   │
│  │  • agents/       │                  │    CLAUDE.md           │   │
│  │                  │                  │  • employer-project/       │   │
│  │  [Shared across  │                  │    CLAUDE.md           │   │
│  │   ALL profiles]  │                  │                        │   │
│  └──────────────────┘                  │  [Context-specific     │   │
│                                         │   additions]           │   │
│                                         └────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              │ ./install.sh
                              │ dotclaude activate <profile>
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         ~/.local/bin/                                │
├─────────────────────────────────────────────────────────────────────┤
│  dotclaude (CLI)  ← Main entry point                                │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              │ Commands (show, list, activate, etc.)
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         ~/.claude/                                   │
│                    (Deployed Configuration)                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  CLAUDE.md  ← base/CLAUDE.md + profiles/X/CLAUDE.md         │   │
│  │                [Merged on activation]                        │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  settings.json  ← base or profile-specific                  │   │
│  │                   [Copied on activation]                     │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  scripts/                                                    │   │
│  │    • dotclaude                                               │   │
│  │    • activate-profile.sh                                     │   │
│  │    • sync-feature-branch.sh                                  │   │
│  │    • shell-functions.sh                                      │   │
│  │    • lib/validation.sh                                       │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  agents/                                                     │   │
│  │    • best-in-class-gap-analysis/                             │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  .current-profile  ← "oss-project"                │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              │ Claude Code reads on startup
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Claude Code Session                           │
├─────────────────────────────────────────────────────────────────────┤
│  • Loads ~/.claude/CLAUDE.md                                        │
│  • Applies ~/.claude/settings.json hooks                            │
│  • Makes agents available                                           │
│  • Executes hooks (SessionStart, PostToolUse, etc.)                 │
└─────────────────────────────────────────────────────────────────────┘
```

## Profile Activation Flow

```
User runs: dotclaude activate my-profile

    ┌─────────────────────────────────────┐
    │  1. Validate Profile Name           │
    │     • Alphanumeric + hyphens only   │
    │     • No path traversal (.., /)     │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  2. Acquire File Lock               │
    │     • ~/.claude/.lock               │
    │     • Timeout: 10 seconds           │
    │     • Prevents concurrent execution │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  3. Check Current Profile           │
    │     • Read ~/.claude/.current-profile│
    │     • Skip backup if same profile   │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  4. Backup Existing Config          │
    │     • CLAUDE.md → .backup.timestamp │
    │     • settings.json → .backup...    │
    │     • chmod 600 (secure)            │
    │     • Keep only 5 recent backups    │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  5. Merge CLAUDE.md                 │
    │     ┌───────────────────────────┐   │
    │     │ base/CLAUDE.md            │   │
    │     │   +                       │   │
    │     │ profiles/my-profile/      │   │
    │     │   CLAUDE.md               │   │
    │     │   ↓                       │   │
    │     │ ~/.claude/CLAUDE.md       │   │
    │     └───────────────────────────┘   │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  6. Apply Settings                  │
    │     • Use profile settings.json if  │
    │       exists, else base             │
    │     • Copy to ~/.claude/            │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  7. Mark Active Profile             │
    │     • Write profile name to         │
    │       ~/.claude/.current-profile    │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  8. Release Lock                    │
    │     • Close file descriptor         │
    │     • Allow next operation          │
    └──────────────┬──────────────────────┘
                   │
                   ▼
              ✓ Complete
```

## CLI Command Flow

```
User runs: dotclaude <command> [args]

    ┌─────────────────────────────────────┐
    │  dotclaude                          │
    │  (bash script with shebang)         │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  Load Validation Library            │
    │  • source lib/validation.sh         │
    │  • Or use fallback inline           │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  Validate Repository Structure      │
    │  • Check REPO_DIR exists            │
    │  • Verify base/ and profiles/ dirs  │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  Set Trap Handler                   │
    │  • cleanup() on EXIT/ERR/INT/TERM   │
    │  • Release locks on exit            │
    └──────────────┬──────────────────────┘
                   │
    ┌──────────────▼──────────────────────┐
    │  Parse Command                      │
    │  • show, list, activate, switch,    │
    │    create, edit, sync, branches,    │
    │    version, help                    │
    └──────────────┬──────────────────────┘
                   │
         ┌─────────┴─────────┬─────────────────┐
         │                   │                  │
    ┌────▼────┐      ┌───────▼────────┐  ┌─────▼──────┐
    │  Show   │      │   Activate     │  │   Sync     │
    │ Profile │      │   Profile      │  │  Feature   │
    └────┬────┘      └───────┬────────┘  │  Branch    │
         │                   │            └─────┬──────┘
         │                   │                  │
         └───────────┬───────┴──────────────────┘
                     │
         ┌───────────▼──────────────┐
         │  Execute Command         │
         │  • Display UI (forest    │
         │    theme)                │
         │  • Perform operations    │
         │  • Handle errors         │
         └───────────┬──────────────┘
                     │
                     ▼
               Return to Shell
```

## Security Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                    Security Layers                                │
└──────────────────────────────────────────────────────────────────┘

Layer 1: Input Validation
─────────────────────────
┌──────────────────────────────┐
│  validate_profile_name()     │
│  • Regex: ^[a-zA-Z0-9_-]+$   │
│  • No path traversal (.., /) │
│  • No special chars          │
└──────────────────────────────┘

Layer 2: Path Safety
────────────────────
┌──────────────────────────────┐
│  validate_directory()        │
│  • Check not symlink         │
│  • Verify real directory     │
│  • Prevent symlink attacks   │
└──────────────────────────────┘

Layer 3: Command Safety
───────────────────────
┌──────────────────────────────┐
│  Single-quoted heredocs      │
│  • Prevent variable expansion│
│  • Use sed for replacement   │
│  • No command injection      │
└──────────────────────────────┘

Layer 4: File Locking
─────────────────────
┌──────────────────────────────┐
│  acquire_lock()              │
│  • flock with timeout        │
│  • Prevent race conditions   │
│  • Concurrent execution safe │
└──────────────────────────────┘

Layer 5: Secure Permissions
───────────────────────────
┌──────────────────────────────┐
│  Backup files: chmod 600     │
│  • Only owner can read       │
│  • Protect sensitive data    │
│  • CLAUDE.md may have secrets│
└──────────────────────────────┘

Layer 6: Safe Removal
─────────────────────
┌──────────────────────────────┐
│  safe_remove_directory()     │
│  • Validate not symlink      │
│  • Check canonical path      │
│  • Must be in safe zones     │
└──────────────────────────────┘
```

## Hook Execution Flow

```
Claude Code Session Starts
    │
    ┌─────────▼────────────────────────────────┐
    │  Claude Code reads ~/.claude/settings.json│
    └─────────┬────────────────────────────────┘
              │
    ┌─────────▼────────────────────────────────┐
    │  Parse hooks configuration               │
    │  {                                       │
    │    "hooks": {                            │
    │      "SessionStart": [...],              │
    │      "PostToolUse": [...]                │
    │    }                                     │
    │  }                                       │
    └─────────┬────────────────────────────────┘
              │
    ┌─────────▼────────────────────────────────┐
    │  Event: SessionStart                     │
    └─────────┬────────────────────────────────┘
              │
    ┌─────────▼────────────────────────────────┐
    │  Match event hooks                       │
    │  • matcher: "*" (all directories)        │
    │  • or matcher: "/path/to/project"        │
    └─────────┬────────────────────────────────┘
              │
    ┌─────────▼────────────────────────────────┐
    │  Execute hook command                    │
    │  • type: "command"                       │
    │  • command: "bash script..."             │
    │  • Runs in bash subshell                 │
    └─────────┬────────────────────────────────┘
              │
    ┌─────────▼────────────────────────────────┐
    │  Display output to user                  │
    │  • stdout shown in Claude UI             │
    │  • stderr shown as error                 │
    └─────────┬────────────────────────────────┘
              │
    ┌─────────▼────────────────────────────────┐
    │  Continue normal operation               │
    └──────────────────────────────────────────┘

Example Hook: SessionStart Git Branch Check
────────────────────────────────────────────
1. Check if in git repo
2. Get current branch name
3. Compare with main/master
4. Calculate commits behind
5. If behind > 0: Display warning
6. Suggest: sync-feature-branch
```

## Data Flow: Profile Merge

```
Merging CLAUDE.md
─────────────────

Input Files:
┌─────────────────────────────────┐
│  base/CLAUDE.md                 │
│  ────────────────               │
│  # Global Instructions          │
│  - Development Standards        │
│  - Code Quality                 │
│  - File Operations              │
│  - Security                     │
│  - Git Practices                │
│  - Tool Usage                   │
│  - Project Context              │
│  - Communication                │
│  ...                            │
└─────────────────────────────────┘
            +
┌─────────────────────────────────┐
│  profiles/my-profile/CLAUDE.md  │
│  ────────────────────────────   │
│  # Profile: My Profile          │
│  - Context-specific standards   │
│  - Tech stack preferences       │
│  - Licensing (for OSS)          │
│  - Compliance (for work)        │
│  - Team practices               │
│  ...                            │
└─────────────────────────────────┘

Merge Process:
┌─────────────────────────────────┐
│  {                              │
│    cat "base/CLAUDE.md"         │
│    echo ""                      │
│    echo "# ==============="     │
│    echo "# Profile: X"          │
│    echo "# ==============="     │
│    echo ""                      │
│    cat "profiles/X/CLAUDE.md"   │
│  } > ~/.claude/CLAUDE.md        │
└─────────────────────────────────┘

Output File:
┌─────────────────────────────────┐
│  ~/.claude/CLAUDE.md            │
│  ────────────────               │
│  [Base content]                 │
│  # Global Instructions          │
│  ...                            │
│                                 │
│  # ===============              │
│  # Profile: My Profile          │
│  # ===============              │
│                                 │
│  [Profile content]              │
│  # Profile-specific additions   │
│  ...                            │
└─────────────────────────────────┘
```

## Installation Architecture

```
./install.sh
    │
    ┌─────────▼───────────────────────────┐
    │  Parse Flags                        │
    │  • --force                          │
    │  • --non-interactive                │
    │  • --help                           │
    └─────────┬───────────────────────────┘
              │
    ┌─────────▼───────────────────────────┐
    │  Check TTY (Interactive?)           │
    │  • if [ ! -t 0 ]; then              │
    │      NON_INTERACTIVE=true           │
    └─────────┬───────────────────────────┘
              │
    ┌─────────▼───────────────────────────┐
    │  Create Directories                 │
    │  • ~/.claude/agents/                │
    │  • ~/.claude/scripts/               │
    │  • ~/.local/bin/                    │
    └─────────┬───────────────────────────┘
              │
    ┌─────────▼───────────────────────────┐
    │  Install dotclaude CLI              │
    │  • Copy to ~/.local/bin/dotclaude   │
    │  • chmod +x                         │
    │  • Check if ~/.local/bin in PATH    │
    └─────────┬───────────────────────────┘
              │
    ┌─────────▼───────────────────────────┐
    │  Install Scripts                    │
    │  • Copy base/scripts/* to           │
    │    ~/.claude/scripts/               │
    │  • chmod +x *.sh                    │
    └─────────┬───────────────────────────┘
              │
    ┌─────────▼───────────────────────────┐
    │  Install Agents                     │
    │  • For each base/agents/*           │
    │  • Check if already exists          │
    │  • Validate not symlink             │
    │  • Prompt or auto-overwrite         │
    │  • Copy to ~/.claude/agents/        │
    └─────────┬───────────────────────────┘
              │
    ┌─────────▼───────────────────────────┐
    │  Select Profile (if interactive)    │
    │  • List available profiles          │
    │  • Prompt user for selection        │
    │  • Or skip if non-interactive       │
    └─────────┬───────────────────────────┘
              │
    ┌─────────▼───────────────────────────┐
    │  Activate Selected Profile          │
    │  • bash activate-profile.sh <name>  │
    │  • Merges CLAUDE.md                 │
    │  • Applies settings.json            │
    └─────────┬───────────────────────────┘
              │
              ▼
        ✓ Complete
```

## Component Responsibilities

### `dotclaude` (Main CLI)
- **Purpose**: Unified command-line interface
- **Responsibilities**:
  - Parse and dispatch commands
  - Load validation library
  - Validate repository structure
  - Handle trap cleanup
  - Display forest-themed UI
- **Key Functions**: `cmd_show()`, `cmd_list()`, `cmd_activate()`, `cmd_switch()`, `cmd_create()`, `cmd_edit()`

### `lib/validation.sh`
- **Purpose**: Centralized security validation
- **Responsibilities**:
  - Profile name validation
  - Directory safety checks
  - Symlink attack prevention
  - File locking
  - Disk space checks
  - Command verification
- **Key Functions**: `validate_profile_name()`, `validate_directory()`, `safe_remove_directory()`, `acquire_lock()`, `release_lock()`

### `activate-profile.sh`
- **Purpose**: Profile activation backend
- **Responsibilities**:
  - Acquire exclusive lock
  - Backup existing config
  - Merge CLAUDE.md files
  - Apply settings.json
  - Mark active profile
  - Release lock on exit
- **Called by**: `dotclaude activate`, `dotclaude switch`

### `sync-feature-branch.sh`
- **Purpose**: Git branch sync automation
- **Responsibilities**:
  - Check git repo status
  - Detect uncommitted changes
  - Calculate ahead/behind commits
  - Guide user through rebase/merge
  - Handle conflicts
  - Push changes safely
- **Called by**: `dotclaude sync`, sourced function `sync-feature-branch`

### `shell-functions.sh`
- **Purpose**: Shell convenience functions
- **Responsibilities**:
  - Provide wrapper functions
  - Git workflow helpers
  - Branch status checking
  - Post-PR workflows
- **Sourced by**: User's `~/.bashrc` or `~/.zshrc`

### `profile-management.sh`
- **Purpose**: Legacy profile management functions
- **Responsibilities**:
  - Profile activation wrapper
  - Show profile status
  - List available profiles
  - Quick profile switches
- **Sourced by**: User's `~/.bashrc` or `~/.zshrc` (optional)

## File System Layout

```
dotclaude/                              # Repository (version controlled)
├── README.md                           # Quick start guide
├── install.sh                          # Installer with flags
├── base/                               # Shared base configuration
│   ├── CLAUDE.md                      # Base development standards
│   ├── settings.json                  # Base hooks & settings
│   ├── scripts/
│   │   ├── dotclaude                  # Main CLI (deployed to ~/.local/bin)
│   │   ├── activate-profile.sh        # Profile activation
│   │   ├── sync-feature-branch.sh     # Git branch sync
│   │   ├── shell-functions.sh         # Shell helpers
│   │   ├── profile-management.sh      # Profile wrappers
│   │   └── lib/
│   │       └── validation.sh          # Security validation
│   └── agents/
│       └── best-in-class-gap-analysis/
│           └── definition.json
├── profiles/                           # Context-specific profiles
│   ├── oss-project/
│   │   └── CLAUDE.md
│   ├── proprietary-project/
│   │   └── CLAUDE.md
│   └── employer-project/
│       └── CLAUDE.md
└── docs/
    ├── USAGE.md                        # Complete user guide
    ├── ARCHITECTURE.md                 # This file
    ├── IDEMPOTENCY-AUDIT.md           # Idempotency analysis
    └── SECURITY-AUDIT.md              # Security hardening details

~/.local/bin/                           # User binaries (in PATH)
└── dotclaude                          # Main CLI (copy from base/scripts)

~/.claude/                              # Deployed configuration
├── .current-profile                   # Active profile name
├── .lock                              # Concurrent execution lock
├── CLAUDE.md                          # Merged: base + profile
├── CLAUDE.md.backup.*                 # Up to 5 recent backups
├── settings.json                      # Active settings
├── settings.json.backup.*             # Up to 5 recent backups
├── scripts/                           # Management scripts (copied)
│   ├── dotclaude
│   ├── activate-profile.sh
│   ├── sync-feature-branch.sh
│   ├── shell-functions.sh
│   ├── profile-management.sh
│   └── lib/
│       └── validation.sh
└── agents/                            # Shared agents (copied)
    └── best-in-class-gap-analysis/
        └── definition.json
```

## Concurrency Model

```
Concurrent Execution Prevention
────────────────────────────────

Process 1: dotclaude activate profile-a
    │
    ├─> acquire_lock("~/.claude/.lock", timeout=10)
    │       exec 200>"~/.claude/.lock"
    │       flock -w 10 200
    │       └─> ✓ Lock acquired
    │
    ├─> Perform activation...
    │
    └─> release_lock()
            exec 200>&-  (close FD)

Process 2: dotclaude activate profile-b
    │
    ├─> acquire_lock("~/.claude/.lock", timeout=10)
    │       exec 200>"~/.claude/.lock"
    │       flock -w 10 200
    │       └─> ⏳ Waiting for lock...
    │           (blocked until Process 1 releases)
    │
    └─> After 10s timeout:
            └─> ✗ Error: Another operation in progress

Key Points:
• Uses flock for advisory file locking
• Timeout prevents indefinite blocking
• Trap handler ensures lock release on exit/error
• Lock file: ~/.claude/.lock
• Lock scope: All dotclaude operations that modify files
```

## Provider-Agnostic Design

```
Multi-Provider Support Architecture
────────────────────────────────────

┌──────────────────────────────────────────────────┐
│  Global Config (~/.claude/)                      │
│  ────────────────────────────────────           │
│  • No hardcoded model IDs                        │
│  • Provider-neutral hooks                        │
│  • Universal agents (no model in definition)    │
└──────────────────┬───────────────────────────────┘
                   │
         ┌─────────┴─────────┐
         │                   │
         ▼                   ▼
┌─────────────────┐  ┌──────────────────┐
│  AWS Bedrock    │  │  Claude Max      │
│  Project        │  │  Project         │
├─────────────────┤  ├──────────────────┤
│  .claude/       │  │  .claude/        │
│  settings.json  │  │  settings.json   │
│  {              │  │  {               │
│    "model":     │  │    "model":      │
│    "us.anthro  │  │    "claude-son   │
│     pic.claude │  │     net-4.5-..."│
│     -sonnet-..." │  │  }               │
│  }              │  │                  │
└─────────────────┘  └──────────────────┘

Settings Precedence:
1. Enterprise policies (if applicable)
2. CLI arguments
3. Project .claude/settings.local.json (gitignored)
4. Project .claude/settings.json (team-shared)
5. Global ~/.claude/settings.json (from dotclaude)

Result: Same global standards, project-specific providers
```

---

## Technical Decisions

### Why Bash?

**Pros:**
- Universal availability on Unix systems
- Simple process execution and file manipulation
- Native git integration
- No external dependencies

**Cons:**
- String handling can be error-prone
- Requires defensive programming
- Less readable than Python/Ruby

**Solution:**
- Use `#!/bin/bash` shebang for all scripts
- Comprehensive validation library
- POSIX-compatible where possible for zsh compatibility

### Why File Locking?

Prevents race conditions when:
- Multiple terminals run `dotclaude activate` simultaneously
- CI/CD and user activate different profiles concurrently
- Multiple processes read/write same files

### Why Merge Instead of Symlink?

**Merged files (current approach):**
- ✓ Self-contained `~/.claude/` directory
- ✓ Works if repository moves/deleted
- ✓ Simple for Claude Code to read
- ✓ Easy backup/restore

**Symlinks (alternative):**
- ✗ Breaks if repository moves
- ✗ Security concerns (symlink attacks)
- ✗ Harder to reason about state
- ✗ Not friendly for backups

### Why Keep Backups?

**Benefits:**
- Safety net for accidental overwrites
- Rollback to previous profile
- Recover from mistakes
- Debug configuration issues

**Limits:**
- Only 5 most recent (prevent disk fill)
- chmod 600 (secure permissions)
- Skip when re-activating same profile (idempotency)

---

**Back to:** [README.md](../README.md) | **See also:** [USAGE.md](USAGE.md)
