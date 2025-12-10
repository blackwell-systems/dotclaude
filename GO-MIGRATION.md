# Go Migration Plan

**Status:** In Progress (75% complete)
**Strategy:** Strangler Fig Pattern
**Timeline:** 1 weekend remaining
**Branch:** `go-migration`

## Overview

Gradual migration from shell-based implementation to Go CLI while maintaining full backward compatibility. Both implementations coexist during development.

## Architecture

### Strangler Fig Pattern

```
dotclaude (wrapper)
    ↓
    ├─→ dotclaude-go (new Go implementation)
    └─→ dotclaude-shell (existing shell implementation)
```

The wrapper intelligently routes commands:
- **Implemented in Go** → Execute Go binary
- **Not yet in Go** → Fall back to shell
- **Manual override** → `DOTCLAUDE_BACKEND=go|shell|auto`

### Directory Structure

```
dotclaude/
├── base/scripts/
│   ├── dotclaude           # Smart wrapper (router)
│   └── dotclaude-shell     # Shell implementation (reference)
├── bin/
│   └── dotclaude-go        # Go binary (compiled)
├── cmd/dotclaude/
│   └── main.go             # Go entry point
├── internal/
│   ├── cli/                # Cobra commands
│   │   ├── root.go
│   │   ├── version.go      # ✅
│   │   ├── list.go         # ✅
│   │   ├── show.go         # ✅
│   │   ├── create.go       # ✅
│   │   ├── delete.go       # ✅
│   │   ├── edit.go         # ✅
│   │   ├── activate.go        # ✅
│   │   ├── restore.go         # ✅
│   │   ├── check_branches.go  # ✅
│   │   └── ... (3 more)
│   └── profile/            # Profile management
│       ├── profile.go      # ✅ Core types and Manager
│       ├── create.go       # ✅ Profile creation
│       ├── delete.go       # ✅ Profile deletion
│       ├── activate.go     # ✅ Profile activation
│       └── restore.go      # ✅ Backup restoration
└── go.mod
```

## Progress

### ✅ Completed Commands (9/12 - 75%)

| Command | Status | Lines | Commit |
|---------|--------|-------|--------|
| **version** | ✅ Complete | ~20 | 8db41e1 |
| **list** | ✅ Complete | ~60 | ba96f9c |
| **show** | ✅ Complete | ~50 | ba96f9c |
| **create** | ✅ Complete | ~180 | 2d32d5c |
| **delete** | ✅ Complete | ~80 | e17314c |
| **edit** | ✅ Complete | ~70 | e17314c |
| **activate** | ✅ Complete | ~220 | 1c2afb3 |
| **restore** | ✅ Complete | ~170 | aa5c989 |
| **check-branches** | ✅ Complete | ~100 | TBD |

### 🔲 Remaining Commands (3/12 - 25%)

| Command | Priority | Complexity | Estimate |
|---------|----------|------------|----------|
| **deactivate** | HIGH | Medium | 2-3 hours |
| **sync** | LOW | Medium | 2-3 hours |
| **feature-branch** | LOW | Medium | 2-3 hours |

**Total Remaining:** ~6-9 hours

**Note:** The "backup" command was removed from the plan as backups are created automatically by the `activate` command. The shell version never implemented a separate backup command.

## Implementation Details

### Completed Components

**Profile Management (`internal/profile/`):**
- ✅ Profile struct with metadata (Name, Path, IsActive, LastModified)
- ✅ Manager with RepoDir, ProfilesDir, ClaudeDir, StateFile (uses .current-profile)
- ✅ ListProfiles() - Read and sort profiles
- ✅ GetActiveProfile() - Read .current-profile state
- ✅ GetActiveProfileName() - Return active name
- ✅ ProfileExists() - Check existence
- ✅ ValidateProfileName() - Validate format (alphanumeric + - _)
- ✅ Create() - Copy template, init git
- ✅ Delete() - Remove profile, safety checks
- ✅ Activate() - Merge base + profile, manage backups
- ✅ copyDir(), copyFile() - Recursive copying with permissions
- ✅ initGitRepo() - Initialize git with initial commit
- ✅ mergeCLAUDEmd() - Merge base + profile with separator
- ✅ applySettings() - Copy settings.json (profile or base fallback)
- ✅ backupFile() - Create timestamped backups (keeps 5 most recent)
- ✅ cleanupBackups() - Remove old backups beyond limit
- ✅ ListBackups() - Find and sort all backup files
- ✅ Restore() - Restore from backup with current file backup
- ✅ updateProfileFromCLAUDE() - Extract profile name from restored CLAUDE.md

**CLI Commands (`internal/cli/`):**
- ✅ root.go - Cobra foundation, global flags, config
- ✅ version.go - Display version
- ✅ list.go - List all profiles with active indicator
- ✅ show.go - Show active profile info
- ✅ create.go - Create new profile from template
- ✅ delete.go - Delete profile with confirmation
- ✅ edit.go - Open CLAUDE.md or settings.json in $EDITOR
- ✅ activate.go - Activate profile (merge base + profile)
- ✅ restore.go - Interactive backup restoration
- ✅ check_branches.go - Check which branches are behind main

### Still Needed

**Profile Management:**
- 🔲 Deactivate() - Restore backup, clean state
- 🔲 Git operations - sync, feature branch

**CLI Commands:**
- 🔲 deactivate.go
- 🔲 sync.go
- 🔲 feature-branch.go

## Testing

### Test Results

**create command:**
```bash
✓ Creates profile from template
✓ Initializes git repository
✓ Creates initial commit
✓ Validates profile name
✓ Prevents duplicate profiles
```

**delete command:**
```bash
✓ Deletes profile directory
✓ Prompts for confirmation
✓ --force skips confirmation
✓ Prevents deleting active profile
✓ Handles non-existent profiles
```

**edit command:**
```bash
✓ Opens CLAUDE.md in $EDITOR
✓ --settings opens settings.json
✓ Waits for editor to close
✓ Falls back to vim if EDITOR unset
```

**list command:**
```bash
✓ Lists all profiles sorted by name
✓ Shows active profile with indicator
✓ Handles empty profiles directory
```

**show command:**
```bash
✓ Displays active profile info
✓ Shows helpful message if none active
✓ Checks Claude directory existence
```

**activate command:**
```bash
✓ Merges base + profile CLAUDE.md with separator
✓ Applies settings.json (profile or base fallback)
✓ Creates timestamped backups on profile switch
✓ Detects re-activation (update in place)
✓ Prevents deleting active profile (delete command)
✓ Keeps only 5 most recent backups
✓ Updates .current-profile state file
✓ Creates Claude directory if missing
```

**restore command:**
```bash
✓ Lists all backups sorted by modification time
✓ Groups backups by type (CLAUDE.md vs settings.json)
✓ Interactive selection with cancel option (q)
✓ Confirms overwrite before restoring
✓ Creates backup of current file before restoring
✓ Updates .current-profile marker when restoring CLAUDE.md
✓ Handles missing backups gracefully
```

### Parity Testing

Comparison with shell version:
- ✅ Same profile structure
- ✅ Same file contents
- ⚠️ Go initializes git, shell doesn't (acceptable difference)
- ✅ Same user-facing behavior

## Build System

**Makefile targets:**
```bash
make build    # Build Go binary
make test     # Run all tests (Go + shell)
make clean    # Remove build artifacts
make install  # Install to ~/bin
```

**Dependencies:**
- Go 1.24+
- github.com/spf13/cobra v1.10.2

## Migration Phases

### Phase 1: Foundation ✅ COMPLETE
- Set up Go project structure
- Create wrapper script
- Implement first command (version)
- Validate build system

**Duration:** 1 hour
**Commits:** 1

### Phase 2: Read-Only Commands ✅ COMPLETE
- Implement list, show
- Create profile management foundation
- Test against shell version

**Duration:** 1 hour
**Commits:** 1

### Phase 3: Write Commands ✅ COMPLETE
- Implement create, delete, edit
- Add file operations (copy, remove)
- Test state changes

**Duration:** 1 hour
**Commits:** 1

### Phase 4: Complex Commands 🟡 IN PROGRESS
- ✅ Implement activate (most critical)
- ✅ Add profile merging logic
- 🔲 Implement deactivate
- 🔲 Test full workflow

**Duration:** 2 hours (activate complete)
**Commits:** TBD

### Phase 5: Git Workflow Commands 🔲 TODO
- Implement sync, check-branches, feature-branch
- Add git integration helpers
- Test git workflows

**Duration:** 4-6 hours estimated
**Commits:** TBD

### Phase 6: Finalization 🔲 TODO
- Run full parity tests
- Update documentation
- Switch default to Go
- Tag v1.0.0

**Duration:** 2-4 hours estimated
**Commits:** TBD

## Timeline

### Actual Progress

| Date | Hours | Work Completed |
|------|-------|----------------|
| 2025-12-10 AM | 3h | Foundation + 6 commands (version, list, show, create, delete, edit) |
| 2025-12-10 PM | 4h | activate + restore + check-branches + container + docs cleanup |

### Estimated Remaining

| Phase | Hours | Status |
|-------|-------|--------|
| Complex Commands | 2-3h | deactivate pending |
| Git Workflow | 4-6h | sync + feature-branch pending |
| Finalization | 2-4h | Pending |
| **Total Remaining** | **6-9h** | **1 weekend** |

## Rollback Strategy

### Per-Command Rollback
```bash
# If create command has issues
export DOTCLAUDE_CREATE_BACKEND=shell
```

### Global Rollback
```bash
# Revert entire system to shell
export DOTCLAUDE_BACKEND=shell
```

### Emergency Rollback
```bash
# Abandon Go migration
git checkout main
# Shell version still works
```

## Success Criteria

Migration is complete when:
- ✅ All 13 commands implemented in Go
- ✅ All tests passing (Go + shell parity)
- ✅ Full workflow test passes
- ✅ No regressions vs shell version
- ✅ Windows support validated
- ✅ Documentation updated

## Notes

- Shell version preserved as `dotclaude-shell` for reference
- Wrapper allows testing both implementations side-by-side
- No users affected during migration (greenfield development)
- Can abort migration at any time by reverting to main branch

---

**Last Updated:** 2025-12-10
**Current Version:** 1.0.0-alpha.1 (Go)
**Shell Version:** 0.5.1 (preserved)
