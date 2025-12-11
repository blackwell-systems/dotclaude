# Shell to Go Migration Strategy

A case study of migrating a CLI tool from shell scripts to Go using the Strangler Fig pattern.

## Overview

This document describes the migration strategy used to move dotclaude from a shell-based implementation (~1000 lines of bash) to a Go implementation (~1500 lines), while maintaining zero downtime and full backward compatibility.

**Timeline:** 7 hours implementation + 2-day validation period
**Result:** Pure Go binary with no shell dependencies

## The Strangler Fig Pattern

Named after strangler fig trees that grow around a host tree and eventually replace it, this pattern allows gradual replacement of a legacy system while keeping it operational.

### Key Principles

1. **Coexistence**: New and old implementations run side-by-side
2. **Incremental Migration**: Replace one component at a time
3. **Easy Rollback**: Can revert to old system at any point
4. **Zero Downtime**: Users experience no disruption

## Architecture Evolution

### Phase 1: Dual Implementation

```
User → Wrapper Script → [Go Binary | Shell Scripts]
           ↓
    Routes based on:
    - Command implementation status
    - DOTCLAUDE_BACKEND env var
```

The wrapper script acted as a router:
- Commands implemented in Go → Go binary
- Commands not yet in Go → Shell fallback
- Manual override via environment variable

### Phase 2: Go Default

```
User → Wrapper Script → Go Binary (default)
                     → Shell (emergency fallback)
```

After all commands were implemented:
- Go became the default backend
- Shell available via `DOTCLAUDE_BACKEND=shell`

### Phase 3: Go Only (Final)

```
User → Go Binary (direct)
       Shell archived
```

After validation:
- Wrapper script removed
- Shell implementation archived
- Go binary installed directly to PATH

## Implementation Phases

### Phase 1: Foundation (1 hour)

**Goal:** Establish Go project structure and prove the pattern works

**Deliverables:**
- Go module setup (`go.mod`, `go.sum`)
- Cobra CLI framework integration
- First command: `version`
- Wrapper script with routing logic

**Validation:**
```bash
dotclaude version  # Routes to Go
dotclaude list     # Routes to Shell (not yet implemented)
```

### Phase 2: Read-Only Commands (1 hour)

**Goal:** Implement safe, read-only commands first

**Commands:**
- `list` - List profiles
- `show` - Show active profile

**Why read-only first:**
- No risk of data corruption
- Easy to verify correctness
- Builds confidence in the implementation

### Phase 3: Write Commands (1 hour)

**Goal:** Implement state-changing commands

**Commands:**
- `create` - Create new profile
- `delete` - Delete profile
- `edit` - Open editor

**Key consideration:** These modify the filesystem, so testing focused on:
- File permissions
- Directory creation
- Safe deletion

### Phase 4: Complex Commands (2 hours)

**Goal:** Implement the most critical command

**Commands:**
- `activate` - The core feature (merges base + profile)
- `restore` - Backup restoration

**Challenges:**
- Profile merging logic had to match shell exactly
- Backup file format compatibility
- State file (.current-profile) format

### Phase 5: Git Workflow Commands (2 hours)

**Goal:** Implement git-related commands

**Commands:**
- `sync` - Sync feature branches
- `check-branches` - Check branch status

**Note:** These were simpler since they mostly wrap git commands.

### Phase 6: Validation (1-2 weeks)

**Goal:** Prove stability in production

**Activities:**
- Run Go as default backend
- Monitor for any issues
- Compare outputs with shell version
- Tag beta releases (v1.0.0-beta.1, beta.2)

### Phase 7: Transition (2 hours)

**Goal:** Remove shell dependencies

**Steps:**
1. Archive shell scripts
2. Remove wrapper script
3. Update install process
4. Update documentation

## Testing Strategy

### Unit Tests

Each Go package had unit tests covering:
- Happy path scenarios
- Error handling
- Edge cases

### Parity Testing

Commands were tested for identical behavior:

```bash
# Shell version
DOTCLAUDE_BACKEND=shell dotclaude list > shell_output.txt

# Go version
DOTCLAUDE_BACKEND=go dotclaude list > go_output.txt

# Compare
diff shell_output.txt go_output.txt
```

### Integration Tests

End-to-end workflows tested:
1. Create profile → Activate → Verify CLAUDE.md content
2. Switch profiles → Verify backup created
3. Restore from backup → Verify state restored

## Rollback Strategy

### During Migration

At any point, users could force shell usage:
```bash
export DOTCLAUDE_BACKEND=shell
```

### After Migration

Shell scripts archived to `archive/` directory with README explaining:
- How to use archived version
- Emergency alias setup
- Older release download links

## Lessons Learned

### What Worked Well

1. **Strangler Fig pattern** - Zero downtime, gradual confidence building
2. **Starting with read-only commands** - Low risk entry point
3. **Keeping both implementations** - Easy debugging and comparison
4. **Environment variable overrides** - Quick rollback mechanism

### Challenges

1. **Behavioral parity** - Some shell quirks were hard to replicate
2. **File format compatibility** - Had to match exact backup formats
3. **Platform differences** - Go made cross-platform easier

### Recommendations

1. **Start small** - Implement the simplest command first
2. **Test parity early** - Catch behavioral differences quickly
3. **Keep the old system running** - Don't delete until confident
4. **Document everything** - Migration documents are invaluable
5. **Use feature flags** - Environment variables enable quick rollback

## Results

| Metric | Shell | Go |
|--------|-------|-----|
| Lines of code | ~1000 | ~1500 |
| Dependencies | bash, git | None (static binary) |
| Platforms | Unix only | Linux, macOS, Windows |
| Install size | N/A (scripts) | 6MB binary |
| Startup time | ~100ms | ~10ms |

## File Structure

### Before Migration

```
base/scripts/
├── dotclaude           # Main entry point
├── dotclaude-shell     # Shell implementation
├── shell-functions.sh  # Shared functions
└── lib/
    └── validation.sh   # Validation functions
```

### After Migration

```
bin/
└── dotclaude           # Go binary (direct)
archive/
├── README.md           # Rollback instructions
├── dotclaude-shell     # Archived shell implementation
└── ...
```

## Conclusion

The Strangler Fig pattern proved ideal for this migration because it:

1. **Eliminated risk** - Users were never affected by partial implementations
2. **Built confidence** - Each phase validated before moving on
3. **Preserved escape hatches** - Rollback always possible
4. **Simplified testing** - Could compare implementations side-by-side

The total migration took 7 hours of implementation plus a validation period, resulting in a more maintainable, cross-platform, and faster CLI tool.

---

*Document created: 2025-12-11*
*Based on dotclaude migration from v0.5.1 (shell) to v1.0.0-rc.1 (Go)*
