#!/bin/bash
# scripts/test-in-container.sh - Build and run Go test container

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_DIR"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╭─────────────────────────────────────────────────────────────╮${NC}"
echo -e "${BLUE}│  Building dotclaude Go Test Container                       │${NC}"
echo -e "${BLUE}╰─────────────────────────────────────────────────────────────╯${NC}"
echo

# Build the container
echo -e "${GREEN}→${NC} Building container image..."
docker build -f Dockerfile.go-test -t dotclaude-go-test .

echo
echo -e "${GREEN}✓${NC} Container built successfully!"
echo
echo -e "${BLUE}╭─────────────────────────────────────────────────────────────╮${NC}"
echo -e "${BLUE}│  Starting Container (--rm = auto-delete on exit)            │${NC}"
echo -e "${BLUE}╰─────────────────────────────────────────────────────────────╯${NC}"
echo

# Run the container
docker run -it --rm dotclaude-go-test
