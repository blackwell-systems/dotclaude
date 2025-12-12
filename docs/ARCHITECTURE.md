# dotclaude Architecture

Technical overview of how dotclaude works internally.

> **Note**: dotclaude v1.0.0-rc.1+ is a pure Go implementation that runs natively on Linux, macOS, and Windows. The legacy shell implementation has been archived to `archive/` for reference only.

## System Architecture

```mermaid
flowchart TB
    repo["<b>dotclaude Repository</b><br/>(Version-Controlled Git Repo)"]
    base["base/<br/>• CLAUDE.md<br/>• settings.json<br/>[Shared across ALL profiles]"]
    profiles["profiles/<br/>• client-work-oss/CLAUDE.md<br/>• client-work/CLAUDE.md<br/>• work-project/CLAUDE.md<br/>[Context-specific additions]"]

    go_binary["<b>~/.local/bin/dotclaude</b><br/>Go Binary<br/>(Cross-platform)"]

    claude_dir["<b>~/.claude/</b><br/>(Deployed Configuration)"]
    merged_claude["CLAUDE.md<br/>(base + profile merged)"]
    deployed_settings["settings.json<br/>(base or profile-specific)"]
    current_profile[".current-profile"]

    session["<b>Claude Code Session</b><br/>• Loads CLAUDE.md<br/>• Applies settings.json hooks<br/>• Executes hooks"]

    repo --> base
    repo --> profiles
    repo -->|"make install"| go_binary
    go_binary -->|"Commands"| claude_dir
    claude_dir --> merged_claude
    claude_dir --> deployed_settings
    claude_dir --> current_profile
    claude_dir -->|"Claude Code<br/>reads on startup"| session

    style repo fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style base fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style profiles fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style go_binary fill:#22543d,stroke:#2f855a,color:#e2e8f0
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
flowchart TB
    start["User runs: dotclaude activate my-profile"]

    subgraph validation["Phase 1: Validation"]
        val1["Validate Profile Name<br/>• Alphanumeric + hyphens only<br/>• No path traversal"]
        val2["Check Profile Exists"]
        val3["Check Current Profile"]

        val1 --> val2 --> val3
    end

    subgraph backup["Phase 2: Backup"]
        check{"Same profile<br/>as current?"}
        skip["Skip backup"]
        bkup["Backup Existing<br/>• CLAUDE.md<br/>• settings.json<br/>• Keep 5 recent"]

        check -->|Yes| skip
        check -->|No| bkup
    end

    subgraph deploy["Phase 3: Deployment"]
        direction LR
        merge["Merge CLAUDE.md<br/>base + profile"]
        settings["Apply Settings<br/>profile or base"]
        mark["Mark Active"]

        merge --> settings --> mark
    end

    complete["✓ Profile Activated"]
    error["✗ Error: Invalid/Not Found"]

    start --> validation
    validation -->|Valid| backup
    validation -.->|Invalid| error
    backup --> deploy
    skip --> deploy
    deploy --> complete

    style start fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style validation fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style backup fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style deploy fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style complete fill:#22543d,stroke:#2f855a,color:#e2e8f0
    style error fill:#742a2a,stroke:#c53030,color:#e2e8f0
    style check fill:#2c5282,stroke:#4299e1,color:#e2e8f0
```

## CLI Command Flow

```mermaid
flowchart LR
    start["User runs:<br/>dotclaude &lt;command&gt; [args]"]

    subgraph cobra_layer["Cobra Framework"]
        binary["Go Binary<br/>~/.local/bin/dotclaude"]
        dispatch["Dispatch<br/>• Find subcommand<br/>• Parse flags<br/>• Validate args"]

        binary --> dispatch
    end

    subgraph execution["Command Execution"]
        manager["Create<br/>Profile Manager"]
        execute["Execute<br/>• Operations<br/>• Error handling<br/>• Output"]

        manager --> execute
    end

    done["Return to Shell"]

    start --> cobra_layer
    cobra_layer --> execution
    execution --> done

    style start fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style cobra_layer fill:#22543d,stroke:#2f855a,color:#e2e8f0
    style execution fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style done fill:#2d3748,stroke:#4a5568,color:#e2e8f0
```

## Security Architecture

```mermaid
flowchart TB
    user["User Input<br/>dotclaude &lt;command&gt; &lt;profile-name&gt;"]

    subgraph input_layer["Input Validation Layer"]
        direction LR
        validate["ValidateProfileName()<br/>• a-zA-Z0-9_- only<br/>• Reject: .., /, spaces"]
        sanitize["Path Safety<br/>• filepath.Join()<br/>• No traversal"]

        validate --> sanitize
    end

    subgraph file_layer["File System Layer"]
        direction LR
        read["Read Operations<br/>• os.ReadFile<br/>• os.ReadDir filtering<br/>• No symlinks"]
        write["Write Operations<br/>• os.WriteFile<br/>• Secure perms<br/>• 0600/0644/0755"]
        delete["Delete Operations<br/>• Check not active<br/>• Validate exists<br/>• os.RemoveAll"]

        read ~~~ write ~~~ delete
    end

    subgraph protection["Protection Layer"]
        no_shell["✓ No shell injection<br/>Pure Go stdlib"]
        no_exec["✓ No command execution<br/>Direct file I/O only"]
        perms["✓ Secure permissions<br/>Owner-only backups"]

        no_shell ~~~ no_exec ~~~ perms
    end

    safe["✓ Safe Operations on<br/>~/.claude/ and profiles/"]

    user --> input_layer
    input_layer --> file_layer
    file_layer --> protection
    protection --> safe

    style user fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style input_layer fill:#742a2a,stroke:#c53030,color:#e2e8f0
    style file_layer fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style protection fill:#22543d,stroke:#2f855a,color:#e2e8f0
    style safe fill:#22543d,stroke:#2f855a,color:#e2e8f0,stroke-width:3px
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

## Implementation Notes

As of v1.0.0-rc.1, dotclaude is a pure Go implementation with no shell dependencies.

**Historical:** The migration from shell to Go used the Strangler Fig pattern. See [SHELL-TO-GO-MIGRATION.md](SHELL-TO-GO-MIGRATION.md) for details on this migration strategy.

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
│   ├── hooks/                          # Hook system
│   └── profile/                        # Profile business logic
├── bin/
│   └── dotclaude                       # Compiled Go binary
├── base/                               # Shared base configuration
│   ├── CLAUDE.md                       # Base development standards
│   ├── settings.json                   # Base hooks & settings
│   ├── hooks/                          # Hook scripts
│   └── agents/                         # Shared agents
├── archive/                            # Archived shell implementation
│   ├── dotclaude-shell                 # Legacy shell CLI
│   └── README.md                       # Rollback instructions
├── profiles/                           # Context-specific profiles
│   ├── my-project/
│   │   └── CLAUDE.md
│   └── work-project/
│       └── CLAUDE.md
├── examples/
│   └── sample-profile/                 # Template for new profiles
└── tests/
    ├── commands.bats                   # Legacy BATS tests
    ├── security.bats                   # Security tests
    └── integration.bats                # Integration tests

~/.local/bin/
└── dotclaude                           # Installed Go binary

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
make build    # Build bin/dotclaude
make test     # Run Go tests
make clean    # Remove bin/
make install  # Install to ~/.local/bin
```

### Dependencies

- **Go 1.23+** - Build requirement
- **github.com/spf13/cobra v1.10.2** - CLI framework

## Testing

### Go Tests (Primary)

```bash
go test ./...                  # All Go tests
go test ./... -cover           # With coverage
go test ./internal/profile/... # Specific package
```

### BATS Tests (Legacy)

The project includes 114+ legacy BATS tests for the archived shell implementation:

| Suite | Count | Description |
|-------|-------|-------------|
| `commands.bats` | 50 | All CLI commands and flags |
| `security.bats` | 40 | Input validation, path safety |
| `integration.bats` | 24 | End-to-end workflows |

```bash
bats tests/                    # All BATS tests
bats tests/commands.bats       # Command tests only
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DOTCLAUDE_REPO_DIR` | `~/code/dotclaude` | Repository location |
| `CLAUDE_DIR` | `~/.claude` | Claude config directory |
| `EDITOR` | `vim` | Editor for `edit` command |

---

**Back to:** [README.md](../README.md) | **See also:** [USAGE.md](USAGE.md) | [GO-MIGRATION.md](../GO-MIGRATION.md)
