#!/bin/bash

# Claude Code Profile Management Functions
# Add these to your ~/.bashrc or ~/.zshrc

# Activate a profile
activate-profile() {
    if [ -f "$HOME/.claude/scripts/activate-profile.sh" ]; then
        bash "$HOME/.claude/scripts/activate-profile.sh" "$@"
    else
        echo "Error: activate-profile.sh not found"
        echo "Run the install script from your dotclaude repo"
        return 1
    fi
}

# Show current active profile
show-profile() {
    echo "╭─────────────────────────────────────────────────────────────╮"
    echo "│  🌲 dotclaude                                               │"
    echo "╰─────────────────────────────────────────────────────────────╯"
    echo ""

    if [ -f "$HOME/.claude/.current-profile" ]; then
        local PROFILE=$(cat "$HOME/.claude/.current-profile")
        echo "  Active profile: $PROFILE"
        echo ""

        if [ -f "$HOME/.claude/CLAUDE.md" ]; then
            local LINES=$(wc -l < "$HOME/.claude/CLAUDE.md")
            echo "  Configuration:"
            echo "    • CLAUDE.md: $LINES lines"
        fi

        if [ -f "$HOME/.claude/settings.json" ]; then
            echo "    • settings.json: configured"
        fi
    else
        echo "  No profile currently active"
        echo ""
        echo "  Run: activate-profile <profile-name>"
    fi

    echo ""
}

# List available profiles
list-profiles() {
    local REPO_DIR="${DOTCLAUDE_REPO_DIR:-$HOME/code/dotclaude}"

    if [ ! -d "$REPO_DIR/profiles" ]; then
        echo "Error: Profiles directory not found at $REPO_DIR/profiles"
        echo "Set DOTCLAUDE_REPO_DIR environment variable if repo is in a different location"
        return 1
    fi

    echo "╭─────────────────────────────────────────────────────────────╮"
    echo "│  🌲 dotclaude                                               │"
    echo "╰─────────────────────────────────────────────────────────────╯"
    echo ""

    local CURRENT=""
    if [ -f "$HOME/.claude/.current-profile" ]; then
        CURRENT=$(cat "$HOME/.claude/.current-profile")
    fi

    echo "  Profiles available:"
    for profile_dir in "$REPO_DIR/profiles"/*; do
        if [ -d "$profile_dir" ]; then
            local profile=$(basename "$profile_dir")
            if [ "$profile" = "$CURRENT" ]; then
                echo "    • $profile (active)"
            else
                echo "    • $profile"
            fi
        fi
    done

    echo ""
    echo "╭─────────────────────────────────────────────────────────────╮"
    echo "│  🍃 Tip: Use 'activate-profile <name>' to switch contexts  │"
    echo "╰─────────────────────────────────────────────────────────────╯"
}

# Export functions
export -f activate-profile 2>/dev/null || true
export -f show-profile 2>/dev/null || true
export -f list-profiles 2>/dev/null || true

echo "Claude Code profile management loaded:"
echo "  activate-profile <name>  - Activate a profile"
echo "  show-profile             - Show current active profile"
echo "  list-profiles            - List all available profiles"
