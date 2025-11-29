#!/bin/bash

# Sync Feature Branch with Main
# Keeps feature branches up-to-date after PRs are merged

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Get current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
DEFAULT_BRANCH=${1:-main}

echo -e "${GREEN}=== Feature Branch Sync Tool ===${NC}"

# Check if we're in a git repo
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# If on main/master, list feature branches that need syncing
if [[ "$CURRENT_BRANCH" == "main" || "$CURRENT_BRANCH" == "master" ]]; then
    echo -e "${YELLOW}Currently on $CURRENT_BRANCH${NC}"
    echo ""
    echo "Feature branches that are behind:"

    # Fetch latest
    git fetch origin --quiet

    # Find branches behind main
    BEHIND_BRANCHES=$(git for-each-ref --format='%(refname:short)' refs/heads/ | \
        while read branch; do
            if [[ "$branch" != "main" && "$branch" != "master" ]]; then
                BEHIND=$(git rev-list --count $branch..$DEFAULT_BRANCH)
                if [ "$BEHIND" -gt 0 ]; then
                    echo "  - $branch (behind by $BEHIND commits)"
                fi
            fi
        done)

    if [ -z "$BEHIND_BRANCHES" ]; then
        echo -e "${GREEN}All feature branches are up-to-date!${NC}"
        exit 0
    else
        echo "$BEHIND_BRANCHES"
        echo ""
        echo "To sync a branch, run:"
        echo "  git checkout <branch-name>"
        echo "  sync-feature-branch"
    fi
    exit 0
fi

# We're on a feature branch
echo -e "Current branch: ${YELLOW}$CURRENT_BRANCH${NC}"

# Check if branch is behind main
BEHIND=$(git rev-list --count $CURRENT_BRANCH..$DEFAULT_BRANCH 2>/dev/null || echo "0")
AHEAD=$(git rev-list --count $DEFAULT_BRANCH..$CURRENT_BRANCH 2>/dev/null || echo "0")

echo "Status: $AHEAD commits ahead, $BEHIND commits behind $DEFAULT_BRANCH"

if [ "$BEHIND" -eq 0 ]; then
    echo -e "${GREEN}✓ Branch is already up-to-date with $DEFAULT_BRANCH${NC}"
    exit 0
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: You have uncommitted changes. Commit or stash them first.${NC}"
    git status --short
    exit 1
fi

echo ""
echo -e "${YELLOW}Branch is $BEHIND commits behind $DEFAULT_BRANCH${NC}"
echo ""
echo "Choose sync method:"
echo "  1) Rebase (cleaner history, requires force push)"
echo "  2) Merge (preserves history, no force push needed)"
echo "  3) Cancel"
echo ""
read -p "Selection (1/2/3): " -n 1 -r
echo ""

case $REPLY in
    1)
        echo -e "${GREEN}Rebasing $CURRENT_BRANCH onto $DEFAULT_BRANCH...${NC}"

        # Fetch latest
        git fetch origin

        # Rebase
        if git rebase origin/$DEFAULT_BRANCH; then
            echo -e "${GREEN}✓ Rebase successful${NC}"
            echo ""
            echo "To push changes:"
            echo -e "${YELLOW}  git push --force-with-lease${NC}"
            echo ""
            read -p "Push now? (y/N): " -n 1 -r
            echo ""
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                git push --force-with-lease
                echo -e "${GREEN}✓ Branch synced and pushed${NC}"
            fi
        else
            echo -e "${RED}Rebase failed. Resolve conflicts and run:${NC}"
            echo "  git rebase --continue"
            echo "  git push --force-with-lease"
            exit 1
        fi
        ;;
    2)
        echo -e "${GREEN}Merging $DEFAULT_BRANCH into $CURRENT_BRANCH...${NC}"

        # Fetch latest
        git fetch origin

        # Merge
        if git merge origin/$DEFAULT_BRANCH -m "Merge $DEFAULT_BRANCH into $CURRENT_BRANCH"; then
            echo -e "${GREEN}✓ Merge successful${NC}"
            echo ""
            echo "To push changes:"
            echo -e "${YELLOW}  git push${NC}"
            echo ""
            read -p "Push now? (y/N): " -n 1 -r
            echo ""
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                git push
                echo -e "${GREEN}✓ Branch synced and pushed${NC}"
            fi
        else
            echo -e "${RED}Merge failed. Resolve conflicts and run:${NC}"
            echo "  git merge --continue"
            echo "  git push"
            exit 1
        fi
        ;;
    3)
        echo "Cancelled"
        exit 0
        ;;
    *)
        echo "Invalid selection"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}=== Sync Complete ===${NC}"
