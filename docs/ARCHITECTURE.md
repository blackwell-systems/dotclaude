# dotclaude Architecture

Technical overview of how dotclaude works internally.

## System Architecture

```mermaid
flowchart TB
    repo["<b>dotclaude Repository</b><br/>(Version-Controlled Git Repo)"]
    base["base/<br/>• CLAUDE.md<br/>• settings.json<br/>• scripts/<br/>• agents/<br/>[Shared across ALL profiles]"]
    profiles["profiles/<br/>• client-work-oss/CLAUDE.md<br/>• client-work/CLAUDE.md<br/>• work-project/CLAUDE.md<br/>[Context-specific additions]"]

    cli["<b>~/.local/bin/dotclaude</b><br/>Main CLI entry point"]

    claude_dir["<b>~/.claude/</b><br/>(Deployed Configuration)"]
    merged_claude["CLAUDE.md<br/>(base + profile merged)"]
    deployed_settings["settings.json<br/>(base or profile-specific)"]
    deployed_scripts["scripts/<br/>• activate-profile.sh<br/>• sync-feature-branch.sh"]
    current_profile[".current-profile"]

    session["<b>Claude Code Session</b><br/>• Loads CLAUDE.md<br/>• Applies settings.json hooks<br/>• Makes agents available<br/>• Executes hooks"]

    repo --> base
    repo --> profiles
    repo -->|"./install.sh<br/>dotclaude activate"| cli
    cli -->|"Commands"| claude_dir
    claude_dir --> merged_claude
    claude_dir --> deployed_settings
    claude_dir --> deployed_scripts
    claude_dir --> current_profile
    claude_dir -->|"Claude Code<br/>reads on startup"| session

    style repo fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style base fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style profiles fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style cli fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style claude_dir fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style merged_claude fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style deployed_settings fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style deployed_scripts fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style current_profile fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style session fill:#2d3748,stroke:#4a5568,color:#e2e8f0
```

## Profile Activation Flow

```mermaid
flowchart TD
    start["User runs: dotclaude activate my-profile"]

    step1["1. Validate Profile Name<br/>• Alphanumeric + hyphens only<br/>• No path traversal (.., /)"]
    step2["2. Acquire File Lock<br/>• ~/.claude/.lock<br/>• Timeout: 10 seconds<br/>• Prevents concurrent execution"]
    step3["3. Check Current Profile<br/>• Read ~/.claude/.current-profile<br/>• Skip backup if same profile"]
    step4["4. Backup Existing Config<br/>• CLAUDE.md → .backup.timestamp<br/>• settings.json → .backup...<br/>• chmod 600 (secure)<br/>• Keep only 5 recent backups"]
    step5["5. Merge CLAUDE.md<br/>base/CLAUDE.md<br/>+<br/>profiles/my-profile/CLAUDE.md<br/>↓<br/>~/.claude/CLAUDE.md"]
    step6["6. Apply Settings<br/>• Use profile settings.json if exists, else base<br/>• Copy to ~/.claude/"]
    step7["7. Mark Active Profile<br/>• Write profile name to<br/>~/.claude/.current-profile"]
    step8["8. Release Lock<br/>• Close file descriptor<br/>• Allow next operation"]
    complete["✓ Complete"]

    start --> step1 --> step2 --> step3 --> step4 --> step5 --> step6 --> step7 --> step8 --> complete

    style start fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style step1 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step2 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step3 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step4 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step5 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step6 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step7 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style step8 fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style complete fill:#22543d,stroke:#2f855a,color:#e2e8f0
```

## CLI Command Flow

```mermaid
flowchart TD
    start["User runs: dotclaude &lt;command&gt; [args]"]

    entry["dotclaude<br/>(bash script with shebang)"]
    load["Load Validation Library<br/>• source lib/validation.sh<br/>• Or use fallback inline"]
    validate["Validate Repository Structure<br/>• Check REPO_DIR exists<br/>• Verify base/ and profiles/ dirs"]
    trap["Set Trap Handler<br/>• cleanup() on EXIT/ERR/INT/TERM<br/>• Release locks on exit"]
    parse["Parse Command<br/>• show, list, activate, switch,<br/>create, edit, sync, branches,<br/>version, help"]

    show["Show<br/>Profile"]
    activate["Activate<br/>Profile"]
    sync["Sync<br/>Feature<br/>Branch"]

    execute["Execute Command<br/>• Display UI (forest theme)<br/>• Perform operations<br/>• Handle errors"]
    done["Return to Shell"]

    start --> entry --> load --> validate --> trap --> parse
    parse --> show
    parse --> activate
    parse --> sync
    show --> execute
    activate --> execute
    sync --> execute
    execute --> done

    style start fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style entry fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style load fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style validate fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style trap fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style parse fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style show fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style activate fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style sync fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style execute fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style done fill:#22543d,stroke:#2f855a,color:#e2e8f0
```

## Security Architecture

```mermaid
graph TB
    title["Security Layers"]

    subgraph layer1["Layer 1: Input Validation"]
        l1_func["validate_profile_name()"]
        l1_rules["• Regex: ^[a-zA-Z0-9_-]+$<br/>• No path traversal (.., /)<br/>• No special chars"]
    end

    subgraph layer2["Layer 2: Path Safety"]
        l2_func["validate_directory()"]
        l2_rules["• Check not symlink<br/>• Verify real directory<br/>• Prevent symlink attacks"]
    end

    subgraph layer3["Layer 3: Command Safety"]
        l3_func["Single-quoted heredocs"]
        l3_rules["• Prevent variable expansion<br/>• Use sed for replacement<br/>• No command injection"]
    end

    subgraph layer4["Layer 4: File Locking"]
        l4_func["acquire_lock()"]
        l4_rules["• flock with timeout<br/>• Prevent race conditions<br/>• Concurrent execution safe"]
    end

    subgraph layer5["Layer 5: Secure Permissions"]
        l5_func["Backup files: chmod 600"]
        l5_rules["• Only owner can read<br/>• Protect sensitive data<br/>• CLAUDE.md may have secrets"]
    end

    subgraph layer6["Layer 6: Safe Removal"]
        l6_func["safe_remove_directory()"]
        l6_rules["• Validate not symlink<br/>• Check canonical path<br/>• Must be in safe zones"]
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

## Hook Execution Flow

```mermaid
sequenceDiagram
    participant User
    participant Claude as Claude Code
    participant Settings as ~/.claude/settings.json
    participant Hook as Hook Command

    User->>Claude: Start Claude Code Session
    Claude->>Settings: Read settings.json
    Settings-->>Claude: Return hooks configuration<br/>{<br/>  "hooks": {<br/>    "SessionStart": [...],<br/>    "PostToolUse": [...]<br/>  }<br/>}

    Claude->>Claude: Parse hooks configuration
    Claude->>Claude: Event: SessionStart
    Claude->>Claude: Match event hooks<br/>• matcher: "*" (all directories)<br/>• or matcher: "/path/to/project"

    Claude->>Hook: Execute hook command<br/>• type: "command"<br/>• command: "bash script..."<br/>• Runs in bash subshell
    Hook-->>Claude: Return output

    Claude->>User: Display output<br/>• stdout shown in Claude UI<br/>• stderr shown as error
    Claude->>Claude: Continue normal operation

    Note over Claude,Hook: Example Hook: SessionStart Git Branch Check<br/>1. Check if in git repo<br/>2. Get current branch name<br/>3. Compare with main/master<br/>4. Calculate commits behind<br/>5. If behind > 0: Display warning<br/>6. Suggest: sync-feature-branch
```

## Data Flow: Profile Merge

```mermaid
flowchart TB
    subgraph inputs["Input Files"]
        base["base/CLAUDE.md<br/>────────────────<br/># Global Instructions<br/>- Development Standards<br/>- Code Quality<br/>- File Operations<br/>- Security<br/>- Git Practices<br/>- Tool Usage<br/>- Project Context<br/>- Communication<br/>..."]

        profile["profiles/my-profile/CLAUDE.md<br/>────────────────────────────<br/># Profile: My Profile<br/>- Context-specific standards<br/>- Tech stack preferences<br/>- Licensing (for OSS)<br/>- Compliance (for work)<br/>- Team practices<br/>..."]
    end

    merge["Merge Process<br/>────────────<br/>{<br/>  cat 'base/CLAUDE.md'<br/>  echo ''<br/>  echo '# ==============='<br/>  echo '# Profile: X'<br/>  echo '# ==============='<br/>  echo ''<br/>  cat 'profiles/X/CLAUDE.md'<br/>} > ~/.claude/CLAUDE.md"]

    output["Output File<br/>────────────<br/>~/.claude/CLAUDE.md<br/><br/>[Base content]<br/># Global Instructions<br/>...<br/><br/># ===============<br/># Profile: My Profile<br/># ===============<br/><br/>[Profile content]<br/># Profile-specific additions<br/>..."]

    base --> merge
    profile --> merge
    merge --> output

    style inputs fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style base fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style profile fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style merge fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style output fill:#22543d,stroke:#2f855a,color:#e2e8f0
```

## Installation Architecture

```mermaid
flowchart TD
    start["./install.sh"]

    flags["Parse Flags<br/>• --force<br/>• --non-interactive<br/>• --help"]
    tty["Check TTY (Interactive?)<br/>• if [ ! -t 0 ]; then<br/>    NON_INTERACTIVE=true"]
    dirs["Create Directories<br/>• ~/.claude/agents/<br/>• ~/.claude/scripts/<br/>• ~/.local/bin/"]
    cli["Install dotclaude CLI<br/>• Copy to ~/.local/bin/dotclaude<br/>• chmod +x<br/>• Check if ~/.local/bin in PATH"]
    scripts["Install Scripts<br/>• Copy base/scripts/* to<br/>  ~/.claude/scripts/<br/>• chmod +x *.sh"]
    agents["Install Agents<br/>• For each base/agents/*<br/>• Check if already exists<br/>• Validate not symlink<br/>• Prompt or auto-overwrite<br/>• Copy to ~/.claude/agents/"]
    select["Select Profile (if interactive)<br/>• List available profiles<br/>• Prompt user for selection<br/>• Or skip if non-interactive"]
    activate["Activate Selected Profile<br/>• bash activate-profile.sh &lt;name&gt;<br/>• Merges CLAUDE.md<br/>• Applies settings.json"]
    complete["✓ Complete"]

    start --> flags --> tty --> dirs --> cli --> scripts --> agents --> select --> activate --> complete

    style start fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style flags fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style tty fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style dirs fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style cli fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style scripts fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style agents fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style select fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style activate fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style complete fill:#22543d,stroke:#2f855a,color:#e2e8f0
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
│   ├── my-project/
│   │   └── CLAUDE.md
│   ├── client-work/
│   │   └── CLAUDE.md
│   └── work-project/
│       └── CLAUDE.md
└── docs/
    ├── USAGE.md                        # Complete user guide
    └── ARCHITECTURE.md                 # This file

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

```mermaid
sequenceDiagram
    participant P1 as Process 1:<br/>dotclaude activate profile-a
    participant Lock as ~/.claude/.lock
    participant P2 as Process 2:<br/>dotclaude activate profile-b

    Note over P1,P2: Concurrent Execution Prevention

    P1->>Lock: acquire_lock(timeout=10)<br/>exec 200>"~/.claude/.lock"<br/>flock -w 10 200
    Lock-->>P1: ✓ Lock acquired

    P2->>Lock: acquire_lock(timeout=10)<br/>exec 200>"~/.claude/.lock"<br/>flock -w 10 200
    Note over P2,Lock: ⏳ Waiting for lock...<br/>(blocked until Process 1 releases)

    P1->>P1: Perform activation...
    P1->>Lock: release_lock()<br/>exec 200>&- (close FD)
    Lock-->>P1: Lock released

    Lock-->>P2: ✗ After 10s timeout:<br/>Error: Another operation in progress

    Note over P1,P2: Key Points:<br/>• Uses flock for advisory file locking<br/>• Timeout prevents indefinite blocking<br/>• Trap handler ensures lock release on exit/error<br/>• Lock file: ~/.claude/.lock<br/>• Lock scope: All dotclaude operations that modify files
```

## Provider-Agnostic Design

```mermaid
graph TB
    global["Global Config (~/.claude/)<br/>────────────────────────────────<br/>• No hardcoded model IDs<br/>• Provider-neutral hooks<br/>• Universal agents (no model in definition)"]

    subgraph bedrock["AWS Bedrock Project"]
        bedrock_dir[".claude/"]
        bedrock_settings["settings.json<br/>{<br/>  'model':<br/>  'us.anthropic.claude-sonnet-...'<br/>}"]
    end

    subgraph claude_max["Claude Max Project"]
        max_dir[".claude/"]
        max_settings["settings.json<br/>{<br/>  'model':<br/>  'claude-sonnet-4.5-...'<br/>}"]
    end

    precedence["Settings Precedence:<br/>1. Enterprise policies (if applicable)<br/>2. CLI arguments<br/>3. Project .claude/settings.local.json (gitignored)<br/>4. Project .claude/settings.json (team-shared)<br/>5. Global ~/.claude/settings.json (from dotclaude)"]

    result["Result: Same global standards, project-specific providers"]

    global --> bedrock
    global --> claude_max
    bedrock --> precedence
    claude_max --> precedence
    precedence --> result

    style global fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style bedrock fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style bedrock_dir fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style bedrock_settings fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style claude_max fill:#1a365d,stroke:#2c5282,color:#e2e8f0
    style max_dir fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style max_settings fill:#2c5282,stroke:#4299e1,color:#e2e8f0
    style precedence fill:#2d3748,stroke:#4a5568,color:#e2e8f0
    style result fill:#22543d,stroke:#2f855a,color:#e2e8f0
```

---

**Back to:** [README.md](../README.md) | **See also:** [USAGE.md](USAGE.md)
