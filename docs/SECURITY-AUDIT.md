# Security & Defensive Programming Audit - dotclaude

Audit performed: 2024-11-29

## Summary

Found **8 safety issues** ranging from critical to low severity.

**Severity:**
- 🔴 Critical: 2 (path traversal, command injection)
- 🟡 Medium: 3 (symlink attacks, missing validation, error cleanup)
- 🟢 Low: 3 (file permissions, disk space, concurrent execution)

---

## Critical Issues 🔴

### 1. Path Traversal in Profile Names

**Location:** All scripts accepting profile names

**Problem:**
```bash
dotclaude activate "../../../etc/passwd"
dotclaude create "../../.ssh/authorized_keys"
```

Profile names are used in file paths without validation:
```bash
PROFILE_DIR="$PROFILES_DIR/$profile_name"
# If profile_name = "../../../tmp", this escapes profiles directory
```

**Attack Vector:**
- Create malicious profile with path traversal
- Read/write arbitrary files on system
- Overwrite critical system files

**Fix:**
Validate profile names to only allow alphanumeric + hyphens:
```bash
validate_profile_name() {
    if [[ ! "$1" =~ ^[a-zA-Z0-9_-]+$ ]]; then
        error_box "Invalid profile name: $1"
        echo "  Profile names must contain only letters, numbers, hyphens, and underscores"
        exit 1
    fi
}
```

### 2. Command Injection via Profile Names

**Location:** `base/scripts/dotclaude:296` (create command)

**Problem:**
```bash
cat > "$profile_dir/CLAUDE.md" <<EOF
# Profile: $profile_name
...
EOF
```

If profile_name contains special chars:
```bash
dotclaude create '$(rm -rf /)'
# Results in: # Profile: $(rm -rf /)
# Executed when file is sourced!
```

**Fix:**
Sanitize profile names AND quote in heredoc:
```bash
cat > "$profile_dir/CLAUDE.md" <<'EOF'
# Profile: PLACEHOLDER
...
EOF
sed -i "s/PLACEHOLDER/$profile_name/g" "$profile_dir/CLAUDE.md"
```

---

## Medium Issues 🟡

### 3. Symlink Attack in Agents Installation

**Location:** `install.sh:111, 119`

**Problem:**
```bash
rm -rf "$target_dir"
cp -r "$agent_dir" "$target_dir"
```

If `$CLAUDE_DIR/agents/agent_name` is a symlink to `/etc/`, this:
```bash
rm -rf "$CLAUDE_DIR/agents/important-dir"  # Follows symlink!
```

Could delete system directories.

**Fix:**
Check for symlinks before rm -rf:
```bash
if [ -L "$target_dir" ]; then
    error_box "Agent directory is a symlink (potential attack)"
    exit 1
fi
```

### 4. Missing Input Validation in Array Access

**Location:** `base/scripts/dotclaude:269`

**Problem:**
```bash
local selected_profile="${profiles[$((choice - 1))]}"
```

Validates numeric range but doesn't handle:
- Very large numbers (array bounds)
- Negative numbers after subtraction
- Integer overflow

**Fix:**
Add bounds checking:
```bash
if [ "$choice" -lt 1 ] || [ "$choice" -gt "${#profiles[@]}" ]; then
    error_box "Selection out of range"
    exit 1
fi
```

(Actually already done, but good)

### 5. No Cleanup on Error

**Location:** All scripts use `set -e`

**Problem:**
If script fails mid-operation:
- Partial files written
- No rollback
- Corrupted state

**Fix:**
Add trap handlers:
```bash
cleanup() {
    if [ -n "$TEMP_FILE" ] && [ -f "$TEMP_FILE" ]; then
        rm -f "$TEMP_FILE"
    fi
}
trap cleanup EXIT ERR
```

---

## Low Severity Issues 🟢

### 6. Sensitive Data in Backups

**Location:** All backup operations

**Problem:**
Backups in `~/.claude/` might contain:
- API keys from settings.json
- Sensitive project paths
- Readable by other processes

**Fix:**
Set restrictive permissions:
```bash
chmod 600 "$BACKUP"
```

### 7. No Disk Space Checks

**Location:** All file write operations

**Problem:**
Large CLAUDE.md files could fill disk:
```bash
cat base.md profile.md > ~/.claude/CLAUDE.md  # No space check
```

**Fix:**
Check available space before large operations:
```bash
available=$(df -k "$CLAUDE_DIR" | tail -1 | awk '{print $4}')
if [ "$available" -lt 1024 ]; then  # Less than 1MB
    error_box "Insufficient disk space"
    exit 1
fi
```

### 8. Concurrent Execution

**Location:** All scripts

**Problem:**
Two processes running simultaneously:
```bash
Terminal 1: dotclaude activate profile-a
Terminal 2: dotclaude activate profile-b
# Race condition on ~/.claude/CLAUDE.md
```

**Fix:**
Use flock for exclusive access:
```bash
exec 200>"$CLAUDE_DIR/.lock"
flock -n 200 || {
    error_box "Another dotclaude operation in progress"
    exit 1
}
```

---

## Additional Defensive Programming Issues

### 9. Spaces in Paths

**Status:** ✅ GOOD - All variables properly quoted

Example:
```bash
cp "$BASE_DIR/settings.json" "$CLAUDE_DIR/settings.json"  # ✓ Quoted
```

### 10. Empty Directory Checks

**Problem:** Some operations assume directories aren't empty:
```bash
cp -r "$BASE_DIR/scripts/"* "$CLAUDE_DIR/scripts/"
# Fails if scripts/ is empty
```

**Fix:**
Check before copying:
```bash
if compgen -G "$BASE_DIR/scripts/*" > /dev/null; then
    cp -r "$BASE_DIR/scripts/"* "$CLAUDE_DIR/scripts/"
fi
```

### 11. REPO_DIR Resolution

**Location:** Multiple scripts

**Problem:**
```bash
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
```

Relies on ${BASH_SOURCE[0]} - fails if:
- Script is sourced
- Run via sh instead of bash
- Symlinked

**Fix:**
Add fallback and validation:
```bash
if [ -z "$REPO_DIR" ] || [ ! -d "$REPO_DIR/profiles" ]; then
    error_box "Cannot determine repository location"
    exit 1
fi
```

---

## Priority Fixes

### Must Fix (Critical) 🔴
1. **Profile name validation** - Prevent path traversal
2. **Sanitize heredoc** - Prevent command injection

### Should Fix (Medium) 🟡
3. **Symlink checks** - Prevent symlink attacks in rm -rf
4. **Trap handlers** - Cleanup on errors
5. **Input validation** - Strengthen bounds checking

### Nice to Have (Low) 🟢
6. **Backup permissions** - chmod 600 on backups
7. **Disk space checks** - Prevent filling disk
8. **File locking** - Prevent concurrent execution

---

## Testing Attack Vectors

### Path Traversal Test
```bash
# Should FAIL safely
dotclaude activate "../../../etc"
dotclaude create "../../malicious"
dotclaude activate "profile/../../../tmp"
```

### Special Characters Test
```bash
# Should FAIL or sanitize
dotclaude create "profile; rm -rf /"
dotclaude create "profile\$(whoami)"
dotclaude create "profile with spaces"
```

### Symlink Attack Test
```bash
# Should detect and prevent
ln -s /etc/passwd ~/.claude/agents/malicious
./install.sh  # Should NOT rm -rf /etc/passwd
```

---

## Recommended Actions

1. Add input validation function used by all scripts
2. Add symlink detection before rm -rf
3. Add trap handlers for cleanup
4. Set secure permissions on sensitive files
5. Add flock for concurrent execution protection
6. Add disk space checks for large operations
7. Add comprehensive tests for attack vectors

## Notes

- Current code is reasonable for personal use
- Critical for public release
- Most issues require malicious local access to exploit
- Defense in depth: multiple layers of validation
