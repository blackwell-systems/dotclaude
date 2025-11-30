#!/bin/bash

# Check for .dotclaude file and detect profile mismatches
# Called by SessionStart hook in settings.json

# Exit silently if no .dotclaude file
if [ ! -f .dotclaude ]; then
    exit 0
fi

# Exit silently if no dotclaude repo configured
if [ -z "$DOTCLAUDE_REPO_DIR" ] && [ ! -d "$HOME/code/dotclaude" ]; then
    exit 0
fi

REPO_DIR="${DOTCLAUDE_REPO_DIR:-$HOME/code/dotclaude}"
CLAUDE_DIR="${CLAUDE_DIR:-$HOME/.claude}"

# Parse .dotclaude file
DESIRED_PROFILE=""

# Support both YAML-style and simple key=value format
if grep -q '^profile:' .dotclaude 2>/dev/null; then
    # YAML format: profile: my-profile
    DESIRED_PROFILE=$(grep '^profile:' .dotclaude | head -1 | sed 's/^profile:[[:space:]]*//' | sed 's/[[:space:]]*$//' | sed 's/^["'"'"']//' | sed 's/["'"'"']$//')
elif grep -q '^profile=' .dotclaude 2>/dev/null; then
    # Shell format: profile=my-profile
    DESIRED_PROFILE=$(grep '^profile=' .dotclaude | head -1 | cut -d= -f2 | sed 's/^["'"'"']//' | sed 's/["'"'"']$//' | xargs)
fi

# Exit if no profile specified
if [ -z "$DESIRED_PROFILE" ]; then
    exit 0
fi

# Validate profile name (security: prevent path traversal)
if [[ ! "$DESIRED_PROFILE" =~ ^[a-zA-Z0-9_-]+$ ]]; then
    echo "⚠️  Invalid profile name in .dotclaude: $DESIRED_PROFILE"
    echo "   Profile names must contain only letters, numbers, hyphens, and underscores"
    exit 0
fi

# Check if profile exists
if [ ! -d "$REPO_DIR/profiles/$DESIRED_PROFILE" ]; then
    echo "⚠️  Profile '$DESIRED_PROFILE' specified in .dotclaude not found"
    echo "   Available profiles: dotclaude list"
    exit 0
fi

# Get current profile
CURRENT_PROFILE=""
if [ -f "$CLAUDE_DIR/.current-profile" ]; then
    CURRENT_PROFILE=$(cat "$CLAUDE_DIR/.current-profile")
fi

# Compare profiles
if [ "$DESIRED_PROFILE" != "$CURRENT_PROFILE" ]; then
    echo ""
    echo "╭─────────────────────────────────────────────────────────────╮"
    echo "│  🍃 Profile Mismatch Detected                               │"
    echo "╰─────────────────────────────────────────────────────────────╯"
    echo ""
    echo "  This project uses:    $DESIRED_PROFILE"
    echo "  Currently active:     ${CURRENT_PROFILE:-none}"
    echo ""
    echo "  To activate the project profile:"
    echo "    dotclaude activate $DESIRED_PROFILE"
    echo ""
fi

exit 0
