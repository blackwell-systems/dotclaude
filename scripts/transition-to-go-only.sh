#!/bin/bash
# scripts/transition-to-go-only.sh
# Phase 7: Transition to Go-only (Option 2 - Direct Binary)
#
# DO NOT RUN THIS UNTIL:
# - Phase 6 soft launch complete
# - 1-2 weeks of production use without issues
# - All validation tests passing
#
# This script removes the wrapper and makes Go the only implementation.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_DIR"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╭─────────────────────────────────────────────────────────────╮${NC}"
echo -e "${BLUE}│  Phase 7: Transition to Go-Only (Option 2)                 │${NC}"
echo -e "${BLUE}╰─────────────────────────────────────────────────────────────╯${NC}"
echo

echo -e "${YELLOW}⚠  WARNING: This will remove the shell fallback!${NC}"
echo
echo "Prerequisites:"
echo "  - Phase 6 soft launch complete"
echo "  - 1-2 weeks of production use without issues"
echo "  - All validation tests passing"
echo
echo "Changes that will be made:"
echo "  1. Archive shell scripts to archive/"
echo "  2. Remove wrapper script"
echo "  3. Copy Go binary to base/scripts/dotclaude"
echo "  4. Update version to 1.0.0"
echo "  5. Commit changes"
echo
read -p "Continue? (yes/NO): " -r
echo

if [ "$REPLY" != "yes" ]; then
    echo "Cancelled. Shell fallback preserved."
    exit 0
fi

echo -e "${GREEN}→${NC} Step 1/5: Archiving shell scripts..."
mkdir -p archive/
git mv base/scripts/dotclaude-shell archive/ 2>/dev/null || true
git mv base/scripts/shell-functions.sh archive/ 2>/dev/null || true
git mv base/scripts/sync-feature-branch.sh archive/ 2>/dev/null || true
echo "  ✓ Shell scripts archived"

echo -e "${GREEN}→${NC} Step 2/5: Building latest Go binary..."
make build
echo "  ✓ Go binary built"

echo -e "${GREEN}→${NC} Step 3/5: Replacing wrapper with Go binary..."
rm -f base/scripts/dotclaude
cp bin/dotclaude-go base/scripts/dotclaude
chmod +x base/scripts/dotclaude
echo "  ✓ Wrapper replaced with direct Go binary"

echo -e "${GREEN}→${NC} Step 4/5: Updating version to 1.0.0..."
sed -i 's/Version = "1.0.0-alpha.[0-9]"/Version = "1.0.0"/' internal/cli/root.go
echo "  ✓ Version updated"

echo -e "${GREEN}→${NC} Step 5/5: Rebuilding with new version..."
make build
echo "  ✓ Binary rebuilt"

echo
echo -e "${BLUE}╭─────────────────────────────────────────────────────────────╮${NC}"
echo -e "${BLUE}│  ✓ Transition Complete                                      │${NC}"
echo -e "${BLUE}╰─────────────────────────────────────────────────────────────╯${NC}"
echo
echo "Changes staged. Review with:"
echo "  git status"
echo "  git diff --cached"
echo
echo "If everything looks good:"
echo "  git commit -m 'feat: Transition to Go-only (Option 2)'"
echo "  git tag v1.0.0"
echo
echo "To test:"
echo "  ./base/scripts/dotclaude version"
echo "  ./base/scripts/dotclaude --help"
echo
echo -e "${YELLOW}Note: Shell scripts preserved in archive/ for reference${NC}"
