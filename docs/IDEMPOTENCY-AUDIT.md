# Idempotency Audit - dotclaude

Audit performed: 2024-11-29

## Summary

Found **5 idempotency issues** across install.sh, activate-profile.sh, and dotclaude CLI.

**Severity:**
- 🔴 Critical: 0
- 🟡 Medium: 3 (unlimited backup files)
- 🟢 Low: 2 (unnecessary overwrites)

---

## Issues Found

### 1. 🟡 Unlimited Backup Files (activate-profile.sh)

**Location:** `base/scripts/activate-profile.sh:57-60, 84-86`

**Problem:**
Every profile activation creates timestamped backups:
```bash
cp "$CLAUDE_DIR/CLAUDE.md" "$CLAUDE_DIR/CLAUDE.md.backup.20241129-120000"
```

Multiple runs = many backups:
```
~/.claude/
  CLAUDE.md.backup.20241129-120000
  CLAUDE.md.backup.20241129-120100
  CLAUDE.md.backup.20241129-120200
  ...
```

**Impact:**
- Disk space waste
- Clutter in ~/.claude/
- No cleanup mechanism

**Recommendation:**
- Check if activating the same profile → skip backup
- OR keep only N most recent backups
- OR don't backup if source matches destination

### 2. 🟡 Unlimited Backup Files (dotclaude activate)

**Location:** `base/scripts/dotclaude:166-167, 188-189`

**Problem:**
Same issue as #1 - cmd_activate() creates unlimited backups.

**Recommendation:**
Same as #1

### 3. 🟢 Unconditional CLI Overwrite (install.sh)

**Location:** `install.sh:32-33`

**Problem:**
```bash
cp "$BASE_DIR/scripts/dotclaude" "$HOME/.local/bin/dotclaude"
```

Always overwrites without checking:
- If file exists
- If it's the same version
- If user has local modifications

**Impact:**
- Loses any local patches
- No version tracking
- Silent overwrite

**Recommendation:**
- Check if file exists and is different before copying
- Compare checksums or versions
- Prompt user if overwriting newer version

### 4. 🟢 Unconditional Scripts Overwrite (install.sh)

**Location:** `install.sh:50`

**Problem:**
```bash
cp -r "$BASE_DIR/scripts/"* "$CLAUDE_DIR/scripts/"
```

Overwrites ALL scripts every time without checking.

**Impact:**
- Loses local modifications to scripts
- No way to preserve custom tweaks

**Recommendation:**
- Selective copy (only missing files)
- OR prompt before overwrite
- OR use rsync with --update flag

### 5. 🟢 Interactive Prompts in install.sh

**Location:** `install.sh:68, 101`

**Problem:**
Prompts for user input:
```bash
read -p "  Overwrite? (y/N): " -n 1 -r
read -p "Which profile would you like to activate? (or 'skip' to skip): " PROFILE_NAME
```

**Impact:**
- Not automatable (hangs in CI/scripts)
- Not truly idempotent (requires user interaction)

**Recommendation:**
- Add `--non-interactive` flag
- Default to safe behavior when stdin is not a TTY
- Environment variable overrides (e.g., `DOTCLAUDE_AUTO_CONFIRM=yes`)

---

## What's Working Well ✅

1. **mkdir -p usage** - All directory creation uses `-p` flag (idempotent)
2. **Profile marker** - `.current-profile` file is overwritten cleanly
3. **CLAUDE.md merging** - Complete regeneration each time (idempotent)
4. **Profile existence checks** - `dotclaude create` checks before creating

---

## Recommended Fixes (Priority Order)

### High Priority

1. **Limit backups to 5 most recent**
   ```bash
   # Keep only 5 most recent backups
   ls -t "$CLAUDE_DIR"/CLAUDE.md.backup.* 2>/dev/null | tail -n +6 | xargs rm -f
   ```

2. **Skip backup if activating same profile**
   ```bash
   CURRENT=$(cat "$CLAUDE_DIR/.current-profile" 2>/dev/null || echo "")
   if [ "$CURRENT" = "$PROFILE_NAME" ]; then
       echo "Already on profile $PROFILE_NAME, skipping backup"
   fi
   ```

### Medium Priority

3. **Add --force flag to install.sh**
   ```bash
   # Only overwrite if --force or file doesn't exist
   if [ ! -f "$HOME/.local/bin/dotclaude" ] || [ "$FORCE" = "true" ]; then
       cp ...
   fi
   ```

4. **Add --non-interactive flag**
   ```bash
   if [ -t 0 ] && [ "$NON_INTERACTIVE" != "true" ]; then
       read -p "..."
   else
       # Auto-confirm or use defaults
   fi
   ```

### Low Priority

5. **Version checking for CLI**
   - Add VERSION to dotclaude script
   - Check installed version vs repo version
   - Only update if repo is newer

---

## Testing Idempotency

To test, run these commands multiple times and check results:

```bash
# Should be safe to run multiple times
./install.sh
./install.sh
./install.sh

# Should produce same result each time
dotclaude activate my-project
dotclaude activate my-project
dotclaude activate my-project

# Check for backup accumulation
ls -la ~/.claude/*.backup.*
```

---

## Notes

- These are **not critical bugs** - the system works correctly
- They're **quality of life** improvements for production use
- Most important: **backup file accumulation** in activate scripts
- All issues have straightforward fixes
