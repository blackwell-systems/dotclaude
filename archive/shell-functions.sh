#!/bin/bash

# Claude Code Global Configuration - Shell Functions
# Add these to your ~/.bashrc or ~/.zshrc

# Sync feature branch with main
sync-feature-branch() {
    if [ -f "$HOME/.claude/scripts/sync-feature-branch.sh" ]; then
        bash "$HOME/.claude/scripts/sync-feature-branch.sh" "$@"
    else
        echo "Error: sync-feature-branch.sh not found"
        echo "Run: cd ~/code/CLAUDE && ./install.sh"
        return 1
    fi
}

# Quick check if branches are behind main
check-branches() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo "Not in a git repository"
        return 1
    fi

    local DEFAULT_BRANCH=${1:-main}
    echo "Checking branches against $DEFAULT_BRANCH..."
    echo ""

    git fetch origin --quiet

    git for-each-ref --format='%(refname:short)' refs/heads/ | \
        while read branch; do
            if [[ "$branch" != "main" && "$branch" != "master" ]]; then
                BEHIND=$(git rev-list --count $branch..$DEFAULT_BRANCH 2>/dev/null || echo "0")
                AHEAD=$(git rev-list --count $DEFAULT_BRANCH..$branch 2>/dev/null || echo "0")

                if [ "$BEHIND" -gt 0 ] || [ "$AHEAD" -gt 0 ]; then
                    printf "  %-30s %s ahead, %s behind\n" "$branch" "$AHEAD" "$BEHIND"
                fi
            fi
        done
}

# After PR merge, update feature branch and continue working
pr-merged() {
    local CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    local DEFAULT_BRANCH=${1:-main}

    echo "PR merged workflow:"
    echo "  1. Switching to $DEFAULT_BRANCH and pulling latest"
    echo "  2. Switching back to $CURRENT_BRANCH and syncing"
    echo ""
    read -p "Continue? (y/N): " -n 1 -r
    echo ""

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Cancelled"
        return 1
    fi

    # Update main
    git checkout "$DEFAULT_BRANCH" && \
    git pull && \
    # Back to feature branch
    git checkout "$CURRENT_BRANCH" && \
    # Sync it
    sync-feature-branch

    if [ $? -eq 0 ]; then
        echo ""
        echo "✓ Feature branch synced and ready for continued development"
    fi
}

# List all feature branches with their status
list-feature-branches() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo "Not in a git repository"
        return 1
    fi

    local DEFAULT_BRANCH=${1:-main}

    git fetch origin --quiet

    echo "Feature branches:"
    echo ""
    printf "  %-30s %-15s %-15s %s\n" "BRANCH" "AHEAD" "BEHIND" "LAST COMMIT"
    printf "  %-30s %-15s %-15s %s\n" "------" "-----" "------" "-----------"

    git for-each-ref --format='%(refname:short)|%(committerdate:relative)' refs/heads/ | \
        while IFS='|' read branch date; do
            if [[ "$branch" != "main" && "$branch" != "master" ]]; then
                BEHIND=$(git rev-list --count $branch..$DEFAULT_BRANCH 2>/dev/null || echo "0")
                AHEAD=$(git rev-list --count $DEFAULT_BRANCH..$branch 2>/dev/null || echo "0")
                printf "  %-30s %-15s %-15s %s\n" "$branch" "$AHEAD" "$BEHIND" "$date"
            fi
        done
}

# Export functions
export -f sync-feature-branch 2>/dev/null || true
export -f check-branches 2>/dev/null || true
export -f pr-merged 2>/dev/null || true
export -f list-feature-branches 2>/dev/null || true

echo "Claude Code git workflow functions loaded:"
echo "  sync-feature-branch  - Sync current branch with main"
echo "  check-branches       - Check all branches status"
echo "  pr-merged            - Workflow after PR is merged"
echo "  list-feature-branches - List all feature branches with status"
