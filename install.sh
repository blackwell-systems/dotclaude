#!/bin/bash

# Claude Code Multi-Profile Configuration Installer
# Installs base configuration and scripts, then activates a profile

set -e

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLAUDE_DIR="$HOME/.claude"
BASE_DIR="$REPO_DIR/global/base"
PROFILES_DIR="$REPO_DIR/global/profiles"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Claude Code Multi-Profile Configuration Installer ===${NC}"
echo "Repo: $REPO_DIR"
echo "Target: $CLAUDE_DIR"
echo ""

# Create ~/.claude if it doesn't exist
mkdir -p "$CLAUDE_DIR/agents"
mkdir -p "$CLAUDE_DIR/scripts"

# Install scripts (needed for profile activation)
echo "[1/2] Installing management scripts..."
if [ -d "$BASE_DIR/scripts" ]; then
    cp -r "$BASE_DIR/scripts/"* "$CLAUDE_DIR/scripts/"
    chmod +x "$CLAUDE_DIR/scripts/"*.sh
    echo -e "  ${GREEN}✓${NC} Installed scripts to ~/.claude/scripts/"
else
    echo "  No scripts found in base"
fi

# Install agents
echo ""
echo "[2/2] Installing global agents..."
if [ -d "$BASE_DIR/agents" ]; then
    for agent_dir in "$BASE_DIR/agents"/*; do
        if [ -d "$agent_dir" ]; then
            agent_name=$(basename "$agent_dir")
            target_dir="$CLAUDE_DIR/agents/$agent_name"

            if [ -d "$target_dir" ]; then
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
        echo "  activate-profile <profile-name>"
    fi
else
    echo "Skipped profile activation"
    echo "To activate a profile later:"
    echo "  activate-profile <profile-name>"
fi

echo ""
echo -e "${BLUE}=== Shell Functions Setup ===${NC}"
echo ""
echo "To enable all helper functions, add to your ~/.bashrc or ~/.zshrc:"
echo ""
echo -e "${YELLOW}  # Claude Code functions${NC}"
echo -e "${YELLOW}  export CLAUDE_REPO_DIR=\"$REPO_DIR\"${NC}"
echo -e "${YELLOW}  if [ -f \"\$HOME/.claude/scripts/shell-functions.sh\" ]; then${NC}"
echo -e "${YELLOW}      source \"\$HOME/.claude/scripts/shell-functions.sh\"${NC}"
echo -e "${YELLOW}  fi${NC}"
echo -e "${YELLOW}  if [ -f \"\$HOME/.claude/scripts/profile-management.sh\" ]; then${NC}"
echo -e "${YELLOW}      source \"\$HOME/.claude/scripts/profile-management.sh\"${NC}"
echo -e "${YELLOW}  fi${NC}"
echo ""
echo "Then restart your shell or run: source ~/.bashrc"
echo ""
echo -e "${BLUE}=== Available Commands ===${NC}"
echo ""
echo "Profile management:"
echo "  - activate-profile <name>  Activate a profile"
echo "  - show-profile             Show current profile"
echo "  - list-profiles            List available profiles"
echo "  - switch-to-oss            Quick switch to blackwell-systems-oss"
echo "  - switch-to-blackwell      Quick switch to blackwell-systems"
echo "  - switch-to-work           Quick switch to best-western"
echo ""
echo "Git workflow:"
echo "  - sync-feature-branch      Sync current branch with main"
echo "  - check-branches           Check all branches status"
echo "  - pr-merged                Workflow after PR is merged"
echo "  - list-feature-branches    List all feature branches"
echo ""
echo -e "${GREEN}Installation complete!${NC}"
