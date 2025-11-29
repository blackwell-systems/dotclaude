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
    echo "  activate-profile blackwell-systems-oss"
    exit 1
}

if [ $# -eq 0 ]; then
    usage
fi

PROFILE_NAME=$1
PROFILE_DIR="$PROFILES_DIR/$PROFILE_NAME"

if [ ! -d "$PROFILE_DIR" ]; then
    echo -e "${RED}Error: Profile '$PROFILE_NAME' not found${NC}"
    echo ""
    usage
fi

echo -e "${BLUE}=== Claude Code Profile Activation ===${NC}"
echo -e "Profile: ${YELLOW}$PROFILE_NAME${NC}"
echo -e "Target: ${YELLOW}$CLAUDE_DIR${NC}"
echo ""

# Create ~/.claude if it doesn't exist
mkdir -p "$CLAUDE_DIR"

# Backup existing CLAUDE.md if it exists
if [ -f "$CLAUDE_DIR/CLAUDE.md" ]; then
    BACKUP="$CLAUDE_DIR/CLAUDE.md.backup.$(date +%Y%m%d-%H%M%S)"
    echo -e "${YELLOW}Backing up existing CLAUDE.md to: $BACKUP${NC}"
    cp "$CLAUDE_DIR/CLAUDE.md" "$BACKUP"
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

    if [ -f "$CLAUDE_DIR/settings.json" ]; then
        BACKUP="$CLAUDE_DIR/settings.json.backup.$(date +%Y%m%d-%H%M%S)"
        echo -e "  ${YELLOW}Backing up existing settings.json to: $BACKUP${NC}"
        cp "$CLAUDE_DIR/settings.json" "$BACKUP"
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
echo -e "${GREEN}=== Profile Activated ===${NC}"
echo ""
echo "Active profile: $PROFILE_NAME"
echo "Configuration deployed to: $CLAUDE_DIR"
echo ""
echo "To verify:"
echo "  show-profile"
echo "  cat ~/.claude/CLAUDE.md"
echo ""
echo "To switch profiles:"
echo "  activate-profile <profile-name>"
