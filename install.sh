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

# Inline validation (Go binary handles all profile operations)
validate_directory() {
    if [ -L "$1" ]; then
        echo -e "${RED}Error: Path is a symlink: $1${NC}" >&2
        return 1
    fi
    return 0
}

echo -e "${BLUE}=== dotclaude installer ===${NC}"
echo "Repo: $REPO_DIR"
echo "Target: $CLAUDE_DIR"
echo ""

# Create directories
mkdir -p "$CLAUDE_DIR/agents"
mkdir -p "$HOME/.local/bin"

# Build and install dotclaude binary
echo "[1/2] Building and installing dotclaude..."

if command -v go >/dev/null 2>&1; then
    # Build the binary
    if [ ! -f "$REPO_DIR/bin/dotclaude" ] || [ "$FORCE_INSTALL" = "true" ]; then
        echo "  Building dotclaude binary..."
        (cd "$REPO_DIR" && make build 2>/dev/null) || {
            echo -e "  ${YELLOW}⚠${NC}  Make not available, using go build directly"
            mkdir -p "$REPO_DIR/bin"
            go build -o "$REPO_DIR/bin/dotclaude" "$REPO_DIR/cmd/dotclaude"
        }
        echo -e "  ${GREEN}✓${NC} Built: $REPO_DIR/bin/dotclaude"
    fi

    # Install to ~/.local/bin
    if [ -f "$REPO_DIR/bin/dotclaude" ]; then
        if [ -f "$HOME/.local/bin/dotclaude" ] && [ "$FORCE_INSTALL" = "false" ]; then
            if cmp -s "$REPO_DIR/bin/dotclaude" "$HOME/.local/bin/dotclaude"; then
                echo -e "  ${GREEN}✓${NC} dotclaude already up-to-date in ~/.local/bin"
            else
                cp "$REPO_DIR/bin/dotclaude" "$HOME/.local/bin/dotclaude"
                chmod +x "$HOME/.local/bin/dotclaude"
                echo -e "  ${GREEN}✓${NC} Updated ~/.local/bin/dotclaude"
            fi
        else
            cp "$REPO_DIR/bin/dotclaude" "$HOME/.local/bin/dotclaude"
            chmod +x "$HOME/.local/bin/dotclaude"
            echo -e "  ${GREEN}✓${NC} Installed to ~/.local/bin/dotclaude"
        fi
    fi
else
    echo -e "  ${YELLOW}⚠${NC}  Go not installed"
    echo "     Option 1: Install Go from https://go.dev/dl/"
    echo "     Option 2: Download pre-built binary from:"
    echo "               https://github.com/blackwell-systems/dotclaude/releases"
fi

# Check PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo -e "  ${YELLOW}⚠${NC}  ~/.local/bin is not in your PATH"
    echo "     Add this to your ~/.bashrc or ~/.zshrc:"
    echo "     export PATH=\"\$HOME/.local/bin:\$PATH\""
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

# Create first profile (mandatory)
echo -e "${BLUE}=== Create Your First Profile ===${NC}"
echo ""
if [ "$NON_INTERACTIVE" = "false" ]; then
    echo "Let's create your first profile. This will include:"
    echo "  • Tech stack preferences (backend, frontend, testing)"
    echo "  • Coding standards (TypeScript, API design, error handling)"
    echo "  • Project workflows and best practices"
    echo ""

    while true; do
        read -p "Profile name (e.g., my-project, work, personal): " PROFILE_NAME

        if [ -z "$PROFILE_NAME" ]; then
            echo -e "${YELLOW}Profile name cannot be empty${NC}"
            continue
        fi

        # Validate profile name
        if [[ ! "$PROFILE_NAME" =~ ^[a-zA-Z0-9_-]+$ ]]; then
            echo -e "${YELLOW}Profile name must contain only letters, numbers, hyphens, and underscores${NC}"
            continue
        fi

        break
    done

    echo ""
    echo "Creating profile: $PROFILE_NAME"
    "$HOME/.local/bin/dotclaude" create "$PROFILE_NAME"

    echo ""
    echo "Activating profile: $PROFILE_NAME"
    "$HOME/.local/bin/dotclaude" activate "$PROFILE_NAME"
else
    echo ""
    echo "Non-interactive mode: Create your first profile with:"
    echo "  dotclaude create <profile-name>"
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
        {
            echo ""
            echo "# dotclaude repository location"
            echo "export DOTCLAUDE_REPO_DIR=\"$REPO_DIR\""
        } >> "$SHELL_RC"
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
echo -e "${BLUE}=== Validating Installation ===${NC}"
echo ""

# Validation checks
VALIDATION_PASS=true

# Check 1: dotclaude CLI in PATH
if command -v dotclaude >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} dotclaude CLI installed and in PATH"
else
    echo -e "${RED}✗${NC} dotclaude CLI not found in PATH"
    echo "  Add ~/.local/bin to your PATH"
    VALIDATION_PASS=false
fi

# Check 2: Binary exists in repo
if [ -f "$REPO_DIR/bin/dotclaude" ]; then
    echo -e "${GREEN}✓${NC} Binary built at $REPO_DIR/bin/dotclaude"
else
    echo -e "${YELLOW}⚠${NC}  Binary not in repo (OK if installed from release)"
fi

# Check 3: Repository accessible
if [ -d "$REPO_DIR/base" ] && [ -f "$REPO_DIR/base/CLAUDE.md" ]; then
    echo -e "${GREEN}✓${NC} Repository accessible at: $REPO_DIR"
else
    echo -e "${YELLOW}⚠${NC}  Repository not found at: $REPO_DIR"
    echo "  Set DOTCLAUDE_REPO_DIR to correct location"
fi

# Check 4: Examples available
if [ -d "$REPO_DIR/examples/sample-profile" ]; then
    echo -e "${GREEN}✓${NC} Sample profile available for templates"
else
    echo -e "${YELLOW}⚠${NC}  Sample profile not found (optional)"
fi

echo ""

if [ "$VALIDATION_PASS" = "true" ]; then
    echo -e "${GREEN}╭─────────────────────────────────────────────────────────────╮${NC}"
    echo -e "${GREEN}│  ✓ Installation Complete - All Checks Passed               │${NC}"
    echo -e "${GREEN}╰─────────────────────────────────────────────────────────────╯${NC}"
else
    echo -e "${YELLOW}╭─────────────────────────────────────────────────────────────╮${NC}"
    echo -e "${YELLOW}│  ⚠  Installation Complete - Some Checks Failed             │${NC}"
    echo -e "${YELLOW}╰─────────────────────────────────────────────────────────────╯${NC}"
fi

echo ""
echo -e "${BLUE}=== Next Steps ===${NC}"
echo ""
if [ "$NON_INTERACTIVE" = "false" ] && [ -n "$PROFILE_NAME" ]; then
    echo "  Your profile '$PROFILE_NAME' is ready! Here's what to do next:"
    echo ""
    echo "  1. Customize your profile for your project:"
    echo -e "     ${GREEN}dotclaude edit $PROFILE_NAME${NC}"
    echo ""
    echo "  2. Verify it's active:"
    echo -e "     ${GREEN}dotclaude show${NC}"
    echo ""
    echo "  3. Create more profiles as needed:"
    echo -e "     ${GREEN}dotclaude create work-project${NC}"
    echo ""
    echo "Other useful commands:"
    echo -e "  ${GREEN}dotclaude list${NC}        List all profiles"
    echo -e "  ${GREEN}dotclaude switch${NC}      Interactive switcher"
    echo -e "  ${GREEN}dotclaude help${NC}        Show all commands"
else
    echo "  1. Create your first profile:"
    echo -e "     ${GREEN}dotclaude create my-project${NC}"
    echo ""
    echo "  2. Activate it:"
    echo -e "     ${GREEN}dotclaude activate my-project${NC}"
    echo ""
    echo "  3. Customize it for your project:"
    echo -e "     ${GREEN}dotclaude edit my-project${NC}"
    echo ""
    echo "Other useful commands:"
    echo -e "  ${GREEN}dotclaude list${NC}        List all profiles"
    echo -e "  ${GREEN}dotclaude switch${NC}      Interactive switcher"
    echo -e "  ${GREEN}dotclaude help${NC}        Show all commands"
fi
echo ""
echo "╭─────────────────────────────────────────────────────────────╮"
echo "│  📖 Documentation: https://blackwell-systems.github.io/dotclaude  │"
echo "╰─────────────────────────────────────────────────────────────╯"
echo ""
