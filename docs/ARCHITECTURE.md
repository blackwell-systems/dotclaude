# dotclaude Architecture

Technical overview of how dotclaude works internally.

> **Note**: dotclaude v1.0.0+ uses a cross-platform Go implementation as the primary (and recommended) backend on all platforms: Linux, macOS, and Windows. The legacy shell implementation is deprecated but available via `DOTCLAUDE_BACKEND=shell` for backwards compatibility on Unix systems.

## System Architecture

```mermaid
flowchart TB
    repo["<b>dotclaude Repository</b><br/>(Version-Controlled Git Repo)"]
    base["base/<br/>• CLAUDE.md<br/>• settings.json<br/>• scripts/<br/>[Shared across ALL profiles]"]
    profiles["profiles/<br/>• client-work-oss/CLAUDE.md<br/>• client-work/CLAUDE.md<br/>• work-project/CLAUDE.md<br/>[Context-specific additions]"]

    wrapper["<b>base/scripts/dotclaude</b><br/>Wrapper Script<br/>Routes to Go or Shell"]
    go_binary["<b>bin/dotclaude-go</b><br/>Go Implementation<br/>(Primary)"]
    shell["<b>base/scripts/dotclaude-shell</b><br/>Shell Implementation<br/>(Fallback/Reference)"]

    claude_dir["<b>~/.claude/</b><br/>(Deployed Configuration)"]
    merged_claude["CLAUDE.md<br/>(base + profile merged)"]
    deployed_settings["settings.json<br/>(base or profile-specific)"]
    current_profile[".current-profile"]

    session["<b>Claude Code Session</b><br/>• Loads CLAUDE.md<br/>• Applies settings.json hooks<br/>• Executes hooks"]

    repo --> base
    repo --> profiles
    repo -->|"make build"| go_binary
    wrapper -->|"DOTCLAUDE_BACKEND=go"| go_binary
    wrapper -->|"DOTCLAUDE_BACKEND=shell"| shell
    go_binary -->|"Commands"| claude_dir
    claude_dir --> merged_claude
    claude_dir --> deployed_settings
    claude_dir --> current_profile
    claude_dir -->|"Claude Code<br/>reads on startup"| session

    style repo fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style base fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style profiles fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style wrapper fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style go_binary fill:#22543d,stroke:#2f855a,color:#e2e8f0
    style shell fill:#4a5568,stroke:#718096,color:#e2e8f0
    style claude_dir fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style merged_claude fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style deployed_settings fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style current_profile fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style session fill:#2d3748,stroke:#4a5568,color:#e2e8f0
```

## Go Implementation Architecture

dotclaude v1.0.0 is implemented in Go using the [Cobra](https://github.com/spf13/cobra) CLI framework.

### Package Structure

```
dotclaude/
├── cmd/dotclaude/
│   └── main.go              # Entry point
├── internal/
│   ├── cli/                 # Cobra command implementations
│   │   ├── root.go          # Root command, global flags, config
│   │   ├── version.go       # version command
│   │   ├── list.go          # list/ls command
│   │   ├── show.go          # show command
│   │   ├── create.go        # create/new command
│   │   ├── delete.go        # delete/rm command
│   │   ├── edit.go          # edit command (cross-platform editor)
│   │   ├── activate.go      # activate/use command
│   │   ├── switch.go        # switch/select command
│   │   ├── restore.go       # restore command
│   │   ├── diff.go          # diff command
│   │   ├── check_branches.go # check-branches command
│   │   ├── sync.go          # sync command
│   │   ├── hook.go          # hook run/list/init commands
│   │   ├── terminal.go      # Cross-platform color support
│   │   ├── terminal_unix.go # Unix terminal handling
│   │   └── terminal_windows.go # Windows ANSI VT support
│   ├── hooks/               # Hook system
│   │   ├── hooks.go         # Hook runner, priority ordering
│   │   └── builtins.go      # Built-in hook implementations
│   └── profile/             # Business logic
│       ├── profile.go       # Manager, Profile types, validation
│       ├── create.go        # Profile creation with git init
│       ├── delete.go        # Safe profile deletion
│       ├── activate.go      # Profile activation with merge
│       └── restore.go       # Backup restoration
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
└── Makefile                 # Build targets
```

### Key Types

```go
// Profile represents a dotclaude profile
type Profile struct {
    Name         string
    Path         string
    IsActive     bool
    LastModified time.Time
}

// Manager handles profile operations
type Manager struct {
    RepoDir     string  // dotclaude repository location
    ProfilesDir string  // RepoDir/profiles
    ClaudeDir   string  // ~/.claude
    StateFile   string  // ~/.claude/.current-profile
}

// Backup represents a backup file
type Backup struct {
    Path      string
    Filename  string
    Timestamp string
    Size      int64
    Type      string  // "CLAUDE.md" or "settings.json"
}
```

## Profile Activation Flow

```mermaid
flowchart TD
    start["User runs: dotclaude activate my-profile"]

    step1["1. Validate Profile Name<br/>• profile.ValidateProfileName()<br/>• Alphanumeric + hyphens only<br/>• No path traversal (.., /)"]
    step2["2. Check Profile Exists<br/>• mgr.ProfileExists(name)<br/>• Returns error if not found"]
    step3["3. Check Current Profile<br/>• mgr.GetActiveProfileName()<br/>• Read ~/.claude/.current-profile<br/>• Skip backup if same profile"]
    step4["4. Backup Existing Config<br/>• mgr.backupFile('CLAUDE.md')<br/>• mgr.backupFile('settings.json')<br/>• Timestamped: *.backup.20241210-155544<br/>• Keep only 5 recent backups"]
    step5["5. Merge CLAUDE.md<br/>• mgr.mergeCLAUDEmd(name)<br/>• Read base/CLAUDE.md<br/>• Read profiles/X/CLAUDE.md<br/>• Concatenate with separator<br/>• Write to ~/.claude/CLAUDE.md"]
    step6["6. Apply Settings<br/>• mgr.applySettings(name)<br/>• Use profile settings.json if exists<br/>• Fall back to base settings.json<br/>• Copy to ~/.claude/"]
    step7["7. Mark Active Profile<br/>• Write profile name to<br/>~/.claude/.current-profile"]
    complete["Complete"]

    start --> step1 --> step2 --> step3 --> step4 --> step5 --> step6 --> step7 --> complete

    style start fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style step1 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step2 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step3 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step4 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step5 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step6 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step7 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style complete fill:#22543d,stroke:#2f855a,color:#e2e8f0
```

## CLI Command Flow

```mermaid
flowchart TD
    start["User runs: dotclaude <command> [args]"]

    wrapper["Wrapper Script<br/>base/scripts/dotclaude<br/>• Check DOTCLAUDE_BACKEND env var<br/>• Default: 'go'"]

    backend_check{"Backend<br/>Selection?"}

    go_path["Go Binary<br/>bin/dotclaude-go<br/>• Cobra framework<br/>• cobra.Command execution"]
    shell_path["Shell Script<br/>base/scripts/dotclaude-shell<br/>• Legacy fallback"]

    cobra["Cobra Dispatch<br/>• rootCmd.Execute()<br/>• Find matching subcommand<br/>• Parse flags<br/>• Validate args"]

    manager["Create Profile Manager<br/>profile.NewManager(RepoDir, ClaudeDir)"]

    execute["Execute Command<br/>• Perform operations<br/>• Handle errors<br/>• Display formatted output"]

    done["Return to Shell"]

    start --> wrapper --> backend_check
    backend_check -->|"go (default)"| go_path
    backend_check -->|"shell"| shell_path
    go_path --> cobra --> manager --> execute --> done
    shell_path --> done

    style start fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style wrapper fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style backend_check fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style go_path fill:#22543d,stroke:#2f855a,color:#e2e8f0
    style shell_path fill:#4a5568,stroke:#718096,color:#e2e8f0
    style cobra fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style manager fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style execute fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style done fill:#22543d,stroke:#2f855a,color:#e2e8f0
```

## Security Architecture

```mermaid
graph TB
    title["Security Layers (Go Implementation)"]

    subgraph layer1["Layer 1: Input Validation"]
        l1_func["profile.ValidateProfileName()"]
        l1_rules["• Check: a-zA-Z0-9_- only<br/>• Reject: empty, .., /, spaces<br/>• Returns error on invalid"]
    end

    subgraph layer2["Layer 2: Path Safety"]
        l2_func["filepath.Join() usage"]
        l2_rules["• OS-agnostic path construction<br/>• Profile path: ProfilesDir/name<br/>• Validated name prevents traversal"]
    end

    subgraph layer3["Layer 3: File Operations"]
        l3_func["Go standard library"]
        l3_rules["• os.ReadFile / os.WriteFile<br/>• No shell command injection<br/>• Direct file I/O"]
    end

    subgraph layer4["Layer 4: Directory Listing"]
        l4_func["os.ReadDir filtering"]
        l4_rules["• entry.IsDir() check<br/>• Symlinks filtered out<br/>• Only real directories listed"]
    end

    subgraph layer5["Layer 5: Secure Permissions"]
        l5_func["File mode settings"]
        l5_rules["• Backups: 0600 (owner only)<br/>• Config files: 0644<br/>• Directories: 0755"]
    end

    subgraph layer6["Layer 6: Safe Deletion"]
        l6_func["profile.Delete()"]
        l6_rules["• Check not active profile<br/>• Validate profile exists<br/>• os.RemoveAll on validated path"]
    end

    title --> layer1 --> layer2 --> layer3 --> layer4 --> layer5 --> layer6

    style title fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style layer1 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style layer2 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style layer3 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style layer4 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style layer5 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style layer6 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
```

## Data Flow: Profile Merge

```mermaid
flowchart TB
    subgraph inputs["Input Files"]
        base["base/CLAUDE.md<br/>────────────────<br/># Global Instructions<br/>- Development Standards<br/>- Code Quality<br/>- File Operations<br/>- Security<br/>- Git Practices<br/>..."]

        profile["profiles/my-profile/CLAUDE.md<br/>────────────────────────────<br/># Profile: My Profile<br/>- Context-specific standards<br/>- Tech stack preferences<br/>- Licensing (for OSS)<br/>- Compliance (for work)<br/>..."]
    end

    merge["mergeCLAUDEmd() in Go<br/>────────────────────<br/>baseContent, _ := os.ReadFile(basePath)<br/>profileContent, _ := os.ReadFile(profilePath)<br/><br/>separator := '# =========...<br/># Profile: X<br/># =========...'<br/><br/>merged := base + separator + profile<br/>os.WriteFile(outputPath, merged, 0644)"]

    output["Output File<br/>────────────<br/>~/.claude/CLAUDE.md<br/><br/>[Base content]<br/># Global Instructions<br/>...<br/><br/># ===============<br/># Profile: my-profile<br/># ===============<br/><br/>[Profile content]<br/>..."]

    base --> merge
    profile --> merge
    merge --> output

    style inputs fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style base fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style profile fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style merge fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style output fill:#22543d,stroke:#2f855a,color:#e2e8f0
```

## Backend Selection (Strangler Fig Pattern)

The wrapper script `base/scripts/dotclaude` routes commands based on `DOTCLAUDE_BACKEND`:

| Value | Behavior |
|-------|----------|
| `go` (default) | Execute Go binary directly |
| `shell` | Execute shell implementation |
| `auto` | Try Go first, fall back to shell for unknown commands |

```bash
# Force Go backend (default)
export DOTCLAUDE_BACKEND=go
dotclaude list

# Force shell backend
export DOTCLAUDE_BACKEND=shell
dotclaude list

# Smart routing (deprecated)
export DOTCLAUDE_BACKEND=auto
dotclaude list
```

## Commands Reference

| Command | Aliases | Description | Flags |
|---------|---------|-------------|-------|
| `version` | - | Show version | - |
| `list` | `ls` | List all profiles | `--verbose` |
| `show` | - | Show active profile | `--debug` |
| `create` | `new` | Create new profile | `--verbose` |
| `delete` | `rm`, `remove` | Delete profile | `--force` |
| `edit` | - | Edit profile in $EDITOR (uses active if no name) | `--settings` |
| `activate` | `use` | Activate profile | `--dry-run`, `--preview`, `--verbose`, `--debug` |
| `switch` | `select` | Interactive profile selector | - |
| `restore` | - | Restore from backup | - |
| `diff` | - | Compare profiles | `--verbose` |
| `check-branches` | `branches`, `br` | Check branch status | `--base` |
| `sync` | - | Sync with main | `--base` |
| `hook run` | - | Execute hooks of a type | - |
| `hook list` | - | List available hooks | - |
| `hook init` | - | Initialize hooks directory | - |

## File System Layout

```
dotclaude/                              # Repository (version controlled)
├── README.md                           # Quick start guide
├── install.sh                          # Installer
├── Makefile                            # Build targets
├── go.mod                              # Go module
├── go.sum                              # Dependencies
├── cmd/dotclaude/
│   └── main.go                         # Go entry point
├── internal/
│   ├── cli/                            # Command implementations
│   └── profile/                        # Profile business logic
├── bin/
│   └── dotclaude-go                    # Compiled Go binary
├── base/                               # Shared base configuration
│   ├── CLAUDE.md                       # Base development standards
│   ├── settings.json                   # Base hooks & settings
│   └── scripts/
│       ├── dotclaude                   # Wrapper script
│       └── dotclaude-shell             # Shell implementation (legacy)
├── profiles/                           # Context-specific profiles
│   ├── my-project/
│   │   └── CLAUDE.md
│   └── work-project/
│       └── CLAUDE.md
├── examples/
│   └── sample-profile/                 # Template for new profiles
└── tests/
    ├── commands.bats                   # Command tests
    ├── security.bats                   # Security tests
    └── integration.bats                # Integration tests

~/.claude/                              # Deployed configuration
├── .current-profile                    # Active profile name
├── CLAUDE.md                           # Merged: base + profile
├── CLAUDE.md.backup.*                  # Up to 5 recent backups
├── settings.json                       # Active settings
└── settings.json.backup.*              # Up to 5 recent backups
```

## Build System

### Makefile Targets

```bash
make build    # Build bin/dotclaude-go
make test     # Run all tests (bats)
make clean    # Remove bin/
make install  # Install to ~/bin
```

### Dependencies

- **Go 1.23+** - Build requirement
- **github.com/spf13/cobra v1.10.2** - CLI framework
- **bats** - Test framework (for tests only)

## Testing

The project includes 114+ automated tests:

| Suite | Count | Description |
|-------|-------|-------------|
| `commands.bats` | 50 | All CLI commands and flags |
| `security.bats` | 40 | Input validation, path safety |
| `integration.bats` | 24 | End-to-end workflows |

Run tests:
```bash
bats tests/                    # All tests
bats tests/commands.bats       # Command tests only
bats tests/security.bats       # Security tests only
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DOTCLAUDE_REPO_DIR` | `~/code/dotclaude` | Repository location |
| `CLAUDE_DIR` | `~/.claude` | Claude config directory |
| `DOTCLAUDE_BACKEND` | `go` | Backend selection |
| `EDITOR` | `vim` | Editor for `edit` command |

---

**Back to:** [README.md](../README.md) | **See also:** [USAGE.md](USAGE.md) | [GO-MIGRATION.md](../GO-MIGRATION.md)
