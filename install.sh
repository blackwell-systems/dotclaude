#!/bin/bash

# dotclaude installer
# Installs base configuration and scripts, then activates a profile

set -e

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLAUDE_DIR="$HOME/.claude"
BASE_DIR="$REPO_DIR/base"
PROFILES_DIR="$REPO_DIR/profiles"

# Parse flags
FORCE_INSTALL=false
NON_INTERACTIVE=false

for arg in "$@"; do
    case "$arg" in
        --force) FORCE_INSTALL=true ;;
        --non-interactive) NON_INTERACTIVE=true ;;
        --help)
            echo "Usage: ./install.sh [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --force            Force overwrite of existing files"
            echo "  --non-interactive  Run without prompting for input"
            echo "  --help             Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $arg"
            echo "Run './install.sh --help' for usage"
            exit 1
            ;;
    esac
done

# Check if running in non-interactive environment
if [ ! -t 0 ]; then
    NON_INTERACTIVE=true
fi

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== dotclaude installer ===${NC}"
echo "Repo: $REPO_DIR"
echo "Target: $CLAUDE_DIR"
echo ""

# Create directories
mkdir -p "$CLAUDE_DIR/agents"
mkdir -p "$CLAUDE_DIR/scripts"
mkdir -p "$HOME/.local/bin"

# Install dotclaude CLI
echo "[1/3] Installing dotclaude CLI..."
if [ -f "$BASE_DIR/scripts/dotclaude" ]; then
    # Check if CLI already exists
    if [ -f "$HOME/.local/bin/dotclaude" ] && [ "$FORCE_INSTALL" = "false" ]; then
        # Check if they're different
        if cmp -s "$BASE_DIR/scripts/dotclaude" "$HOME/.local/bin/dotclaude"; then
            echo -e "  ${GREEN}✓${NC} dotclaude CLI already up-to-date"
        else
            echo -e "  ${YELLOW}⚠${NC}  dotclaude CLI already exists (use --force to overwrite)"
            echo "     Existing: ~/.local/bin/dotclaude"
            echo "     New version: $BASE_DIR/scripts/dotclaude"
        fi
    else
        cp "$BASE_DIR/scripts/dotclaude" "$HOME/.local/bin/dotclaude"
        chmod +x "$HOME/.local/bin/dotclaude"
        echo -e "  ${GREEN}✓${NC} Installed to ~/.local/bin/dotclaude"
    fi

    # Check if ~/.local/bin is in PATH
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo -e "  ${YELLOW}⚠${NC}  ~/.local/bin is not in your PATH"
        echo "     Add this to your ~/.bashrc or ~/.zshrc:"
        echo "     export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
else
    echo "  ${YELLOW}⚠${NC}  dotclaude CLI not found in base/scripts"
fi

# Install scripts (needed for profile activation)
echo ""
echo "[2/3] Installing management scripts..."
if [ -d "$BASE_DIR/scripts" ]; then
    cp -r "$BASE_DIR/scripts/"* "$CLAUDE_DIR/scripts/"
    chmod +x "$CLAUDE_DIR/scripts/"*.sh
    echo -e "  ${GREEN}✓${NC} Installed scripts to ~/.claude/scripts/"
else
    echo "  No scripts found in base"
fi

# Install agents
echo ""
echo "[3/3] Installing global agents..."
if [ -d "$BASE_DIR/agents" ]; then
    for agent_dir in "$BASE_DIR/agents"/*; do
        if [ -d "$agent_dir" ]; then
            agent_name=$(basename "$agent_dir")
            target_dir="$CLAUDE_DIR/agents/$agent_name"

            if [ -d "$target_dir" ]; then
                if [ "$NON_INTERACTIVE" = "true" ] || [ "$FORCE_INSTALL" = "true" ]; then
                    # Auto-overwrite in non-interactive or force mode
                    rm -rf "$target_dir"
                    cp -r "$agent_dir" "$target_dir"
                    echo -e "  ${GREEN}✓${NC} Installed agent: $agent_name"
                else
                    echo -e "  ${YELLOW}Warning: Agent '$agent_name' already exists${NC}"
                    read -p "  Overwrite? (y/N): " -n 1 -r
                    echo
                    if [[ $REPLY =~ ^[Yy]$ ]]; then
                        rm -rf "$target_dir"
                        cp -r "$agent_dir" "$target_dir"
                        echo -e "  ${GREEN}✓${NC} Installed agent: $agent_name"
                    else
                        echo "  ✗ Skipped agent: $agent_name"
                    fi
                fi
            else
                cp -r "$agent_dir" "$target_dir"
                echo -e "  ${GREEN}✓${NC} Installed agent: $agent_name"
            fi
        fi
    done
else
    echo "  No agents found in base"
fi

echo ""
echo -e "${GREEN}=== Base Installation Complete ===${NC}"
echo ""

# Ask which profile to activate
echo -e "${BLUE}=== Profile Selection ===${NC}"
echo ""
echo "Available profiles:"
for profile_dir in "$PROFILES_DIR"/*; do
    if [ -d "$profile_dir" ]; then
        echo "  - $(basename "$profile_dir")"
    fi
done
if [ "$NON_INTERACTIVE" = "false" ]; then
    echo ""
    read -p "Which profile would you like to activate? (or 'skip' to skip): " PROFILE_NAME
    echo ""

    if [[ "$PROFILE_NAME" != "skip" && -n "$PROFILE_NAME" ]]; then
        if [ -d "$PROFILES_DIR/$PROFILE_NAME" ]; then
            echo "Activating profile: $PROFILE_NAME"
            bash "$CLAUDE_DIR/scripts/activate-profile.sh" "$PROFILE_NAME"
        else
            echo -e "${YELLOW}Profile '$PROFILE_NAME' not found. Skipping profile activation.${NC}"
            echo "You can activate a profile later with:"
            echo "  dotclaude activate <profile-name>"
        fi
    else
        echo "Skipped profile activation"
        echo "To activate a profile later:"
        echo "  dotclaude activate <profile-name>"
    fi
else
    echo ""
    echo "Skipping profile activation (non-interactive mode)"
    echo "To activate a profile later:"
    echo "  dotclaude activate <profile-name>"
fi

echo ""
echo -e "${BLUE}=== Setup Complete ===${NC}"
echo ""
echo "Set DOTCLAUDE_REPO_DIR in your shell (optional):"
echo ""
echo -e "${YELLOW}  # Add to ~/.bashrc or ~/.zshrc${NC}"
echo -e "${YELLOW}  export DOTCLAUDE_REPO_DIR=\"$REPO_DIR\"${NC}"
echo ""
echo -e "${BLUE}=== Getting Started ===${NC}"
echo ""
echo "Try these commands:"
echo ""
echo "  ${GREEN}dotclaude show${NC}              Show current profile"
echo "  ${GREEN}dotclaude list${NC}              List available profiles"
echo "  ${GREEN}dotclaude switch${NC}            Interactive profile switcher"
echo "  ${GREEN}dotclaude help${NC}              Show all commands"
echo ""
echo -e "${GREEN}╭─────────────────────────────────────────────────────────────╮${NC}"
echo -e "${GREEN}│  ✓ Installation Complete                                    │${NC}"
echo -e "${GREEN}╰─────────────────────────────────────────────────────────────╯${NC}"
echo ""
echo "╭─────────────────────────────────────────────────────────────╮"
echo "│  🍃 Tip: Run 'dotclaude help' to see all commands          │"
echo "╰─────────────────────────────────────────────────────────────╯"
