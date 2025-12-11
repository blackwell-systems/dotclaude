#!/bin/bash

# dotclaude validation library
# Defensive programming utilities for safe script execution

# Colors for error messages
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Validate profile name (prevent path traversal and injection)
validate_profile_name() {
    local name="$1"
    local context="${2:-Profile name}"

    if [ -z "$name" ]; then
        echo -e "${RED}Error: $context cannot be empty${NC}" >&2
        return 1
    fi

    # Only allow alphanumeric, hyphens, and underscores
    if [[ ! "$name" =~ ^[a-zA-Z0-9_-]+$ ]]; then
        echo -e "${RED}Error: Invalid $context: '$name'${NC}" >&2
        echo "  Profile names must contain only:" >&2
        echo "    • Letters (a-z, A-Z)" >&2
        echo "    • Numbers (0-9)" >&2
        echo "    • Hyphens (-)" >&2
        echo "    • Underscores (_)" >&2
        echo "" >&2
        echo "  Examples of valid names:" >&2
        echo "    • my-profile" >&2
        echo "    • dev_environment" >&2
        echo "    • profile2024" >&2
        return 1
    fi

    # Additional safety: check for suspicious patterns
    if [[ "$name" =~ \.\. ]] || [[ "$name" =~ / ]]; then
        echo -e "${RED}Error: $context contains path traversal characters${NC}" >&2
        return 1
    fi

    return 0
}

# Validate directory path (prevent symlink attacks)
validate_directory() {
    local dir="$1"
    local context="${2:-Directory}"

    if [ ! -e "$dir" ]; then
        # Path doesn't exist - safe to use
        return 0
    fi

    if [ ! -d "$dir" ]; then
        echo -e "${RED}Error: $context is not a directory: $dir${NC}" >&2
        return 1
    fi

    if [ -L "$dir" ]; then
        echo -e "${RED}Error: $context is a symlink (potential security risk): $dir${NC}" >&2
        echo "  Refusing to operate on symlinked directories" >&2
        return 1
    fi

    return 0
}

# Safe remove directory (checks for symlinks first)
safe_remove_directory() {
    local dir="$1"

    if [ ! -e "$dir" ]; then
        return 0  # Nothing to remove
    fi

    # Validate it's not a symlink
    if [ -L "$dir" ]; then
        echo -e "${RED}Error: Refusing to remove symlink: $dir${NC}" >&2
        return 1
    fi

    if [ ! -d "$dir" ]; then
        echo -e "${RED}Error: Not a directory: $dir${NC}" >&2
        return 1
    fi

    # Additional safety: ensure it's under expected parent
    local canonical_path
    canonical_path=$(cd "$dir" && pwd -P)

    # Must be under ~/.claude/agents/ or profiles/
    if [[ "$canonical_path" != "$HOME/.claude/agents/"* ]] && \
       [[ "$canonical_path" != *"/profiles/"* ]]; then
        echo -e "${RED}Error: Directory outside safe zones: $canonical_path${NC}" >&2
        return 1
    fi

    rm -rf "$dir"
}

# Check disk space before large operations
check_disk_space() {
    local target_dir="$1"
    local required_kb="${2:-1024}"  # Default 1MB

    local available
    available=$(df -k "$target_dir" | tail -1 | awk '{print $4}')

    if [ "$available" -lt "$required_kb" ]; then
        echo -e "${RED}Error: Insufficient disk space${NC}" >&2
        echo "  Available: ${available}KB" >&2
        echo "  Required: ${required_kb}KB" >&2
        return 1
    fi

    return 0
}

# Acquire exclusive lock (prevent concurrent execution)
acquire_lock() {
    local lockfile="$1"
    local timeout="${2:-5}"

    # Create parent directory if it doesn't exist
    local lockdir
    lockdir="$(dirname "$lockfile")"
    if [ ! -d "$lockdir" ]; then
        mkdir -p "$lockdir" || return 1
    fi

    exec 200>"$lockfile"

    if ! flock -w "$timeout" 200; then
        echo -e "${YELLOW}Warning: Another dotclaude operation in progress${NC}" >&2
        echo "  Waiting for lock: $lockfile" >&2
        echo "  Timeout: ${timeout}s" >&2
        return 1
    fi

    return 0
}

# Release lock
release_lock() {
    exec 200>&-  # Close file descriptor
}

# Validate file doesn't contain sensitive patterns
check_sensitive_data() {
    local file="$1"

    if [ ! -f "$file" ]; then
        return 0
    fi

    # Check for common sensitive patterns
    local patterns=(
        "api[_-]?key"
        "secret"
        "password"
        "token"
        "private[_-]?key"
        "access[_-]?key"
    )

    for pattern in "${patterns[@]}"; do
        if grep -qi "$pattern" "$file" 2>/dev/null; then
            echo -e "${YELLOW}Warning: File may contain sensitive data: $file${NC}" >&2
            echo "  Pattern matched: $pattern" >&2
            return 1
        fi
    done

    return 0
}

# Validate required commands exist
require_commands() {
    local missing=()

    for cmd in "$@"; do
        if ! command -v "$cmd" &> /dev/null; then
            missing+=("$cmd")
        fi
    done

    if [ ${#missing[@]} -gt 0 ]; then
        echo -e "${RED}Error: Required commands not found:${NC}" >&2
        for cmd in "${missing[@]}"; do
            echo "  • $cmd" >&2
        done
        return 1
    fi

    return 0
}

# Sanitize string for safe output
sanitize_output() {
    local input="$1"
    # Remove ANSI escape codes and control characters
    echo "$input" | sed 's/\x1b\[[0-9;]*m//g' | tr -cd '[:print:]\n'
}

# Validate REPO_DIR is correctly set
validate_repo_dir() {
    local repo_dir="$1"

    if [ -z "$repo_dir" ]; then
        echo -e "${RED}Error: Repository directory not set${NC}" >&2
        return 1
    fi

    if [ ! -d "$repo_dir" ]; then
        echo -e "${RED}Error: Repository directory does not exist: $repo_dir${NC}" >&2
        return 1
    fi

    if [ ! -d "$repo_dir/base" ] || [ ! -d "$repo_dir/profiles" ]; then
        echo -e "${RED}Error: Invalid repository structure at: $repo_dir${NC}" >&2
        echo "  Expected: base/ and profiles/ directories" >&2
        return 1
    fi

    return 0
}

# Export functions for use in other scripts
export -f validate_profile_name 2>/dev/null || true
export -f validate_directory 2>/dev/null || true
export -f safe_remove_directory 2>/dev/null || true
export -f check_disk_space 2>/dev/null || true
export -f acquire_lock 2>/dev/null || true
export -f release_lock 2>/dev/null || true
export -f check_sensitive_data 2>/dev/null || true
export -f require_commands 2>/dev/null || true
export -f sanitize_output 2>/dev/null || true
export -f validate_repo_dir 2>/dev/null || true
