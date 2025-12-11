# Archived Shell Implementation

This directory contains the original shell-based implementation of dotclaude, preserved for reference.

**Status:** Archived (superseded by Go implementation)

## Files

| File | Purpose |
|------|---------|
| `dotclaude-shell` | Main shell CLI (38KB, 1000+ lines) |
| `shell-functions.sh` | Shared utility functions |
| `sync-feature-branch.sh` | Git branch sync functionality |
| `activate-profile.sh` | Profile activation logic |
| `check-dotclaude.sh` | Auto-detection checks |
| `profile-management.sh` | Profile CRUD operations |
| `lib/validation.sh` | Input validation functions |

## History

- **v0.1.0 - v0.5.1**: Shell-only implementation
- **v1.0.0-beta.1**: Go implementation introduced (strangler fig pattern)
- **v1.0.0-beta.2**: Go as default backend
- **v1.0.0-rc.1**: Shell version archived, Go-only

## Emergency Rollback

If you need to use the shell version temporarily:

```bash
# Copy shell version back
cp archive/dotclaude-shell base/scripts/dotclaude-shell
cp archive/shell-functions.sh base/scripts/
cp -r archive/lib base/scripts/

# Update wrapper to use shell
export DOTCLAUDE_BACKEND=shell
```

## Why Go?

- Cross-platform binaries (Linux, macOS, Windows)
- No shell dependencies
- Better error handling
- Easier testing
- Faster execution

---

*Archived: 2025-12-11*
