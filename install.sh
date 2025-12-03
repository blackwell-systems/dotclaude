#!/bin/bash

# dotclaude installer
# Installs base configuration and scripts, then activates a profile
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.sh | bash
#
# Or clone first:
#   git clone https://github.com/blackwell-systems/dotclaude.git ~/code/dotclaude
#   cd ~/code/dotclaude
#   ./install.sh

set -e

# Colors (defined early for clone message)
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Check if we're running from a git clone or need to clone
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" 2>/dev/null && pwd || echo "")"

if [ -z "$SCRIPT_DIR" ] || [ ! -f "$SCRIPT_DIR/base/CLAUDE.md" ]; then
    # Running via curl | bash, need to clone repo first
    CLONE_DIR="$HOME/code/dotclaude"

    echo -e "${BLUE}=== dotclaude installer ===${NC}"
    echo ""
    echo "Installing dotclaude to: $CLONE_DIR"
    echo ""

    # Check if directory already exists
    if [ -d "$CLONE_DIR" ]; then
        echo -e "${YELLOW}Directory already exists: $CLONE_DIR${NC}"
        read -p "Remove and re-clone? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$CLONE_DIR"
        else
            echo "Installation cancelled."
            echo "To install manually:"
            echo "  cd $CLONE_DIR"
            echo "  ./install.sh --force"
            exit 1
        fi
    fi

    # Clone the repository
    echo "Cloning repository..."
    if ! git clone https://github.com/blackwell-systems/dotclaude.git "$CLONE_DIR"; then
        echo -e "${RED}Failed to clone repository${NC}"
        echo "Please ensure git is installed and you have internet access"
        exit 1
    fi

    echo -e "${GREEN}✓${NC} Repository cloned to $CLONE_DIR"
    echo ""

    # Re-exec the script from the cloned location
    echo "Running installer from cloned repository..."
    echo ""
    exec bash "$CLONE_DIR/install.sh" "$@"
fi

# If we get here, we're running from a cloned repo
REPO_DIR="${DOTCLAUDE_REPO_DIR:-$SCRIPT_DIR}"
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

# Load validation library if available
if [ -f "$BASE_DIR/scripts/lib/validation.sh" ]; then
    source "$BASE_DIR/scripts/lib/validation.sh"
else
    # Fallback inline validation
    validate_directory() {
        if [ -L "$1" ]; then
            echo -e "${RED}Error: Path is a symlink: $1${NC}" >&2
            return 1
        fi
        return 0
    }
fi

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
                # Validate not a symlink before removing (symlink attack prevention)
                if [ -L "$target_dir" ]; then
                    echo -e "  ${RED}✗${NC} Agent '$agent_name' is a symlink (refusing to remove)"
                    continue
                fi

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

# Add DOTCLAUDE_REPO_DIR to shell RC if not already present
SHELL_RC=""
if [ -n "$BASH_VERSION" ]; then
    SHELL_RC="$HOME/.bashrc"
elif [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="$HOME/.zshrc"
else
    # Default to bashrc
    SHELL_RC="$HOME/.bashrc"
fi

if [ -f "$SHELL_RC" ]; then
    if ! grep -q "DOTCLAUDE_REPO_DIR" "$SHELL_RC"; then
        echo "" >> "$SHELL_RC"
        echo "# dotclaude repository location" >> "$SHELL_RC"
        echo "export DOTCLAUDE_REPO_DIR=\"$REPO_DIR\"" >> "$SHELL_RC"
        echo -e "${GREEN}✓${NC} Added DOTCLAUDE_REPO_DIR to $SHELL_RC"
        echo ""
        echo -e "${YELLOW}Run 'source $SHELL_RC' or restart your shell${NC}"
    else
        echo -e "${GREEN}✓${NC} DOTCLAUDE_REPO_DIR already configured"
    fi
else
    echo -e "${YELLOW}Note:${NC} Add to your shell configuration:"
    echo "  export DOTCLAUDE_REPO_DIR=\"$REPO_DIR\""
fi

echo ""
echo -e "${GREEN}╭─────────────────────────────────────────────────────────────╮${NC}"
echo -e "${GREEN}│  ✓ Installation Complete                                    │${NC}"
echo -e "${GREEN}╰─────────────────────────────────────────────────────────────╯${NC}"
echo ""
echo -e "${BLUE}=== Next Steps ===${NC}"
echo ""
echo "  1. Create your first profile:"
echo -e "     ${GREEN}dotclaude create my-project${NC}"
echo ""
echo "  2. Edit it to add your project context:"
echo -e "     ${GREEN}dotclaude edit my-project${NC}"
echo ""
echo "  3. Activate it:"
echo -e "     ${GREEN}dotclaude activate my-project${NC}"
echo ""
echo "  4. Verify it's active:"
echo -e "     ${GREEN}dotclaude show${NC}"
echo ""
echo "Other useful commands:"
echo -e "  ${GREEN}dotclaude list${NC}        List all profiles"
echo -e "  ${GREEN}dotclaude switch${NC}      Interactive switcher"
echo -e "  ${GREEN}dotclaude help${NC}        Show all commands"
echo ""
echo "╭─────────────────────────────────────────────────────────────╮"
echo "│  📖 Documentation: https://blackwell-systems.github.io/dotclaude  │"
echo "╰─────────────────────────────────────────────────────────────╯"
echo ""
