#!/bin/bash

# Claude Code Global Configuration Installer
# Syncs configurations from this repo to ~/.claude/

set -e

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLAUDE_DIR="$HOME/.claude"

echo "=== Claude Code Global Configuration Installer ==="
echo "Repo: $REPO_DIR"
echo "Target: $CLAUDE_DIR"
echo ""

# Create ~/.claude if it doesn't exist
mkdir -p "$CLAUDE_DIR/agents"
mkdir -p "$CLAUDE_DIR/scripts"

# Function to backup existing file
backup_file() {
    local file="$1"
    if [ -f "$file" ]; then
        local backup="${file}.backup.$(date +%Y%m%d-%H%M%S)"
        echo "  Backing up existing file to: $backup"
        cp "$file" "$backup"
    fi
}

# Install CLAUDE.md
echo "[1/4] Installing global CLAUDE.md..."
if [ -f "$CLAUDE_DIR/CLAUDE.md" ] && [ -s "$CLAUDE_DIR/CLAUDE.md" ]; then
    echo "  Warning: $CLAUDE_DIR/CLAUDE.md already exists and is not empty"
    read -p "  Overwrite? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        backup_file "$CLAUDE_DIR/CLAUDE.md"
        cp "$REPO_DIR/global/CLAUDE.md" "$CLAUDE_DIR/CLAUDE.md"
        echo "  ✓ Installed CLAUDE.md"
    else
        echo "  ✗ Skipped CLAUDE.md"
    fi
else
    cp "$REPO_DIR/global/CLAUDE.md" "$CLAUDE_DIR/CLAUDE.md"
    echo "  ✓ Installed CLAUDE.md"
fi

# Install settings.json (merge, don't overwrite)
echo ""
echo "[2/4] Installing global settings.json..."
if [ -f "$CLAUDE_DIR/settings.json" ]; then
    echo "  Warning: $CLAUDE_DIR/settings.json already exists"
    echo "  Current settings:"
    cat "$CLAUDE_DIR/settings.json"
    echo ""
    echo "  New settings from repo:"
    cat "$REPO_DIR/global/settings.json"
    echo ""
    read -p "  Overwrite? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        backup_file "$CLAUDE_DIR/settings.json"
        cp "$REPO_DIR/global/settings.json" "$CLAUDE_DIR/settings.json"
        echo "  ✓ Installed settings.json"
    else
        echo "  ✗ Skipped settings.json"
        echo "  Tip: Manually merge hooks if needed"
    fi
else
    cp "$REPO_DIR/global/settings.json" "$CLAUDE_DIR/settings.json"
    echo "  ✓ Installed settings.json"
fi

# Install scripts
echo ""
echo "[3/4] Installing helper scripts..."
if [ -d "$REPO_DIR/global/scripts" ]; then
    cp -r "$REPO_DIR/global/scripts/"* "$CLAUDE_DIR/scripts/"
    chmod +x "$CLAUDE_DIR/scripts/"*.sh
    echo "  ✓ Installed scripts to ~/.claude/scripts/"
else
    echo "  No scripts found in repo"
fi

# Install agents
echo ""
echo "[4/4] Installing global agents..."
if [ -d "$REPO_DIR/global/agents" ]; then
    for agent_dir in "$REPO_DIR/global/agents"/*; do
        if [ -d "$agent_dir" ]; then
            agent_name=$(basename "$agent_dir")
            target_dir="$CLAUDE_DIR/agents/$agent_name"

            if [ -d "$target_dir" ]; then
                echo "  Warning: Agent '$agent_name' already exists"
                read -p "  Overwrite? (y/N): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    rm -rf "$target_dir"
                    cp -r "$agent_dir" "$target_dir"
                    echo "  ✓ Installed agent: $agent_name"
                else
                    echo "  ✗ Skipped agent: $agent_name"
                fi
            else
                cp -r "$agent_dir" "$target_dir"
                echo "  ✓ Installed agent: $agent_name"
            fi
        fi
    done
else
    echo "  No agents found in repo"
fi

echo ""
echo "=== Installation Complete ==="
echo ""
echo "Installed to: $CLAUDE_DIR"
echo ""
echo "To verify installation:"
echo "  ls -la ~/.claude/"
echo "  cat ~/.claude/CLAUDE.md"
echo "  cat ~/.claude/settings.json"
echo "  ls ~/.claude/agents/"
echo "  ls ~/.claude/scripts/"
echo ""
echo "=== Shell Functions Setup ==="
echo ""
echo "To enable git workflow helper functions, add this to your ~/.bashrc or ~/.zshrc:"
echo ""
echo "  # Claude Code git workflow functions"
echo "  if [ -f \"\$HOME/.claude/scripts/shell-functions.sh\" ]; then"
echo "      source \"\$HOME/.claude/scripts/shell-functions.sh\""
echo "  fi"
echo ""
echo "Then restart your shell or run: source ~/.bashrc"
echo ""
echo "Available commands after setup:"
echo "  - sync-feature-branch      Sync current branch with main"
echo "  - check-branches           Check all branches status"
echo "  - pr-merged                Workflow after PR is merged"
echo "  - list-feature-branches    List all feature branches"
echo ""
echo "Backups are stored with .backup.TIMESTAMP extension"
