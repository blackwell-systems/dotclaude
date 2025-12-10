# Go Migration Plan

**Status:** In Progress (46% complete)
**Strategy:** Strangler Fig Pattern
**Timeline:** 4 weeks estimated
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
│   │   └── ... (7 more)
│   └── profile/            # Profile management
│       ├── profile.go      # ✅ Core types and Manager
│       ├── create.go       # ✅ Profile creation
│       └── delete.go       # ✅ Profile deletion
└── go.mod
```

## Progress

### ✅ Completed Commands (6/13 - 46%)

| Command | Status | Lines | Commit |
|---------|--------|-------|--------|
| **version** | ✅ Complete | ~20 | 8db41e1 |
| **list** | ✅ Complete | ~60 | ba96f9c |
| **show** | ✅ Complete | ~50 | ba96f9c |
| **create** | ✅ Complete | ~180 | 2d32d5c |
| **delete** | ✅ Complete | ~80 | e17314c |
| **edit** | ✅ Complete | ~70 | e17314c |

### 🔲 Remaining Commands (7/13 - 54%)

| Command | Priority | Complexity | Estimate |
|---------|----------|------------|----------|
| **activate** | HIGH | Complex | 4-6 hours |
| **deactivate** | HIGH | Medium | 2-3 hours |
| **backup** | MEDIUM | Simple | 1-2 hours |
| **restore** | MEDIUM | Simple | 1-2 hours |
| **sync** | LOW | Medium | 2-3 hours |
| **check-branches** | LOW | Simple | 1 hour |
| **feature-branch** | LOW | Medium | 2-3 hours |

**Total Remaining:** ~13-20 hours

## Implementation Details

### Completed Components

**Profile Management (`internal/profile/`):**
- ✅ Profile struct with metadata (Name, Path, IsActive, LastModified)
- ✅ Manager with RepoDir, ProfilesDir, ClaudeDir, StateFile
- ✅ ListProfiles() - Read and sort profiles
- ✅ GetActiveProfile() - Read .dotclaude-active state
- ✅ GetActiveProfileName() - Return active name
- ✅ ProfileExists() - Check existence
- ✅ ValidateProfileName() - Validate format (alphanumeric + - _)
- ✅ Create() - Copy template, init git
- ✅ Delete() - Remove profile, safety checks
- ✅ copyDir(), copyFile() - Recursive copying with permissions
- ✅ initGitRepo() - Initialize git with initial commit

**CLI Commands (`internal/cli/`):**
- ✅ root.go - Cobra foundation, global flags, config
- ✅ version.go - Display version
- ✅ list.go - List all profiles with active indicator
- ✅ show.go - Show active profile info
- ✅ create.go - Create new profile from template
- ✅ delete.go - Delete profile with confirmation
- ✅ edit.go - Open CLAUDE.md or settings.json in $EDITOR

### Still Needed

**Profile Management:**
- 🔲 Activate() - Merge base + profile, symlink to .claude
- 🔲 Deactivate() - Restore backup, clean state
- 🔲 Backup() - Copy .claude to backup location
- 🔲 Restore() - Restore .claude from backup
- 🔲 Merge() - Combine base/CLAUDE.md + profile/CLAUDE.md
- 🔲 Git operations - sync, branch checking, feature branch

**CLI Commands:**
- 🔲 activate.go - Most complex command
- 🔲 deactivate.go
- 🔲 backup.go
- 🔲 restore.go
- 🔲 sync.go
- 🔲 check-branches.go
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

### Phase 4: Complex Commands 🔲 IN PROGRESS
- Implement activate (most critical)
- Add profile merging logic
- Implement deactivate
- Test full workflow

**Duration:** 4-6 hours estimated
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
| 2025-12-10 | 3h | Foundation + 6 commands |

### Estimated Remaining

| Phase | Hours | Status |
|-------|-------|--------|
| Complex Commands | 4-6h | Next |
| Git Workflow | 4-6h | Pending |
| Finalization | 2-4h | Pending |
| **Total Remaining** | **10-16h** | **2-3 weekends** |

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
