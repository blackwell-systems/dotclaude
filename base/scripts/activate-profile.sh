#!/bin/bash

# Claude Code Profile Activation
# Activates a specific profile by merging base + profile configs to ~/.claude/

set -e

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CLAUDE_DIR="$HOME/.claude"
PROFILES_DIR="$REPO_DIR/profiles"
BASE_DIR="$REPO_DIR/base"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Load validation library
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/lib/validation.sh" ]; then
    source "$SCRIPT_DIR/lib/validation.sh"
else
    # Fallback inline validation
    validate_profile_name() {
        if [[ ! "$1" =~ ^[a-zA-Z0-9_-]+$ ]]; then
            echo -e "${RED}Error: Invalid profile name: $1${NC}" >&2
            return 1
        fi
    }
fi

# Trap handler for cleanup
cleanup() {
    if [ -n "$LOCKFILE" ] && [ -f "$LOCKFILE" ]; then
        rm -f "$LOCKFILE" 2>/dev/null || true
    fi
}
trap cleanup EXIT ERR INT TERM

usage() {
    echo "Usage: activate-profile <profile-name>"
    echo ""
    echo "Available profiles:"
    for profile in "$PROFILES_DIR"/*; do
        if [ -d "$profile" ]; then
            echo "  - $(basename "$profile")"
        fi
    done
    echo ""
    echo "Example:"
    echo "  activate-profile my-project"
    exit 1
}

if [ $# -eq 0 ]; then
    usage
fi

PROFILE_NAME=$1

# Validate profile name (prevent path traversal and injection)
if ! validate_profile_name "$PROFILE_NAME"; then
    exit 1
fi

PROFILE_DIR="$PROFILES_DIR/$PROFILE_NAME"

if [ ! -d "$PROFILE_DIR" ]; then
    echo -e "${RED}Error: Profile '$PROFILE_NAME' not found${NC}"
    echo ""
    usage
fi

# Acquire lock to prevent concurrent operations
LOCKFILE="$CLAUDE_DIR/.activation.lock"
exec 200>"$LOCKFILE"
if ! flock -w 10 200; then
    echo -e "${RED}Error: Another activation in progress${NC}"
    echo "  Wait for the other operation to complete"
    exit 1
fi

echo -e "${BLUE}=== Claude Code Profile Activation ===${NC}"
echo -e "Profile: ${YELLOW}$PROFILE_NAME${NC}"
echo -e "Target: ${YELLOW}$CLAUDE_DIR${NC}"
echo ""

# Create ~/.claude if it doesn't exist
mkdir -p "$CLAUDE_DIR"

# Check if we're re-activating the same profile (skip backup if so)
CURRENT_PROFILE=""
if [ -f "$CLAUDE_DIR/.current-profile" ]; then
    CURRENT_PROFILE=$(cat "$CLAUDE_DIR/.current-profile")
fi

# Backup existing CLAUDE.md if it exists and we're switching profiles
if [ -f "$CLAUDE_DIR/CLAUDE.md" ] && [ "$CURRENT_PROFILE" != "$PROFILE_NAME" ]; then
    BACKUP="$CLAUDE_DIR/CLAUDE.md.backup.$(date +%Y%m%d-%H%M%S)"
    echo -e "${YELLOW}Backing up existing CLAUDE.md${NC}"
    cp "$CLAUDE_DIR/CLAUDE.md" "$BACKUP"

    # Set secure permissions on backup (may contain sensitive data)
    chmod 600 "$BACKUP" 2>/dev/null || true

    # Keep only 5 most recent backups
    ls -t "$CLAUDE_DIR"/CLAUDE.md.backup.* 2>/dev/null | tail -n +6 | xargs rm -f 2>/dev/null || true
elif [ "$CURRENT_PROFILE" = "$PROFILE_NAME" ]; then
    echo -e "${GREEN}Already on profile '$PROFILE_NAME', updating in place${NC}"
fi

# Merge base CLAUDE.md + profile CLAUDE.md
echo "[1/3] Merging CLAUDE.md (base + profile)..."
{
    cat "$BASE_DIR/CLAUDE.md"
    echo ""
    echo "# ========================================="
    echo "# Profile-Specific Additions: $PROFILE_NAME"
    echo "# ========================================="
    echo ""
    cat "$PROFILE_DIR/CLAUDE.md"
} > "$CLAUDE_DIR/CLAUDE.md"
echo -e "  ${GREEN}✓${NC} CLAUDE.md merged"

# Copy or merge settings.json
echo ""
echo "[2/3] Handling settings.json..."
if [ -f "$PROFILE_DIR/settings.json" ]; then
    echo -e "  ${YELLOW}Note: Profile has custom settings.json${NC}"
    echo "  Base settings will be overridden by profile settings"

    # Only backup if switching profiles
    if [ -f "$CLAUDE_DIR/settings.json" ] && [ "$CURRENT_PROFILE" != "$PROFILE_NAME" ]; then
        BACKUP="$CLAUDE_DIR/settings.json.backup.$(date +%Y%m%d-%H%M%S)"
        echo -e "  ${YELLOW}Backing up existing settings.json${NC}"
        cp "$CLAUDE_DIR/settings.json" "$BACKUP"

        # Set secure permissions on backup (contains hooks/sensitive config)
        chmod 600 "$BACKUP" 2>/dev/null || true

        # Keep only 5 most recent backups
        ls -t "$CLAUDE_DIR"/settings.json.backup.* 2>/dev/null | tail -n +6 | xargs rm -f 2>/dev/null || true
    fi

    # For now, just use profile settings (could merge with jq in future)
    cp "$PROFILE_DIR/settings.json" "$CLAUDE_DIR/settings.json"
    echo -e "  ${GREEN}✓${NC} Profile settings.json applied"
else
    # Use base settings
    cp "$BASE_DIR/settings.json" "$CLAUDE_DIR/settings.json"
    echo -e "  ${GREEN}✓${NC} Base settings.json applied"
fi

# Set current profile marker
echo ""
echo "[3/3] Marking active profile..."
echo "$PROFILE_NAME" > "$CLAUDE_DIR/.current-profile"
echo -e "  ${GREEN}✓${NC} Profile marker set"

echo ""
echo -e "${GREEN}╭─────────────────────────────────────────────────────────────╮${NC}"
echo -e "${GREEN}│  ✓ Profile Activated                                        │${NC}"
echo -e "${GREEN}╰─────────────────────────────────────────────────────────────╯${NC}"
echo ""
echo "  Active profile: $PROFILE_NAME"
echo "  Configuration deployed to: $CLAUDE_DIR"
echo ""
echo "  Verify:"
echo "    • show-profile"
echo "    • cat ~/.claude/CLAUDE.md"
echo ""
echo "╭─────────────────────────────────────────────────────────────╮"
echo "│  🍃 Tip: Use 'show-profile' to see your current setup      │"
echo "╰─────────────────────────────────────────────────────────────╯"
