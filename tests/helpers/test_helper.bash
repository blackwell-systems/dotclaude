#!/bin/bash

# Test helper functions for dotclaude test suite
# Provides setup, teardown, and utility functions for bats tests

# Test directories
export TEST_TEMP_DIR="${BATS_TEST_TMPDIR:-/tmp}/dotclaude-test-$$"
export TEST_REPO_DIR="$TEST_TEMP_DIR/dotclaude"
export TEST_CLAUDE_DIR="$TEST_TEMP_DIR/.claude"
export TEST_HOME="$TEST_TEMP_DIR/home"

# Setup test environment before each test
setup_test_env() {
    # Create test directories
    mkdir -p "$TEST_REPO_DIR"/{base,profiles,examples}
    mkdir -p "$TEST_CLAUDE_DIR"
    mkdir -p "$TEST_HOME"

    # Create minimal base/CLAUDE.md
    cat > "$TEST_REPO_DIR/base/CLAUDE.md" <<'EOF'
# Base Configuration

This is the base configuration that applies to all profiles.

## Universal Standards
- Git workflow
- Security practices
- Tool usage
EOF

    # Create minimal base/settings.json
    cat > "$TEST_REPO_DIR/base/settings.json" <<'EOF'
{
  "version": "1.0.0"
}
EOF

    # Copy actual scripts to test environment
    if [ -d "$BATS_TEST_DIRNAME/../base/scripts" ]; then
        cp -r "$BATS_TEST_DIRNAME/../base/scripts" "$TEST_REPO_DIR/base/"
    fi

    # Set environment variables for tests
    export DOTCLAUDE_REPO_DIR="$TEST_REPO_DIR"
    export HOME="$TEST_HOME"
    export CLAUDE_DIR="$TEST_CLAUDE_DIR"

    # Make dotclaude command available
    export PATH="$TEST_REPO_DIR/base/scripts:$PATH"
}

# Cleanup test environment after each test
teardown_test_env() {
    if [ -d "$TEST_TEMP_DIR" ]; then
        rm -rf "$TEST_TEMP_DIR"
    fi
}

# Create a test profile
create_test_profile() {
    local profile_name="$1"
    local content="${2:-# Profile: $profile_name}"

    mkdir -p "$TEST_REPO_DIR/profiles/$profile_name"
    cat > "$TEST_REPO_DIR/profiles/$profile_name/CLAUDE.md" <<EOF
$content
EOF
}

# Create multiple test profiles
create_test_profiles() {
    create_test_profile "test-profile-1" "# Test Profile 1\n\nNode.js project"
    create_test_profile "test-profile-2" "# Test Profile 2\n\nPython project"
    create_test_profile "test-profile-3" "# Test Profile 3\n\nRust project"
}

# Assert file exists
assert_file_exists() {
    local file="$1"
    if [ ! -f "$file" ]; then
        echo "Expected file to exist: $file" >&2
        return 1
    fi
}

# Assert file does not exist
assert_file_not_exists() {
    local file="$1"
    if [ -f "$file" ]; then
        echo "Expected file to not exist: $file" >&2
        return 1
    fi
}

# Assert file contains string
assert_file_contains() {
    local file="$1"
    local pattern="$2"

    if [ ! -f "$file" ]; then
        echo "File does not exist: $file" >&2
        return 1
    fi

    if ! grep -q "$pattern" "$file"; then
        echo "File does not contain '$pattern': $file" >&2
        return 1
    fi
}

# Assert file does not contain string
assert_file_not_contains() {
    local file="$1"
    local pattern="$2"

    if [ ! -f "$file" ]; then
        echo "File does not exist: $file" >&2
        return 1
    fi

    if grep -q "$pattern" "$file"; then
        echo "File should not contain '$pattern': $file" >&2
        return 1
    fi
}

# Assert directory exists
assert_dir_exists() {
    local dir="$1"
    if [ ! -d "$dir" ]; then
        echo "Expected directory to exist: $dir" >&2
        return 1
    fi
}

# Assert output contains string
assert_output_contains() {
    local pattern="$1"
    if [[ ! "$output" =~ $pattern ]]; then
        echo "Expected output to contain '$pattern'" >&2
        echo "Actual output: $output" >&2
        return 1
    fi
}

# Assert output matches regex
assert_output_matches() {
    local pattern="$1"
    if [[ ! "$output" =~ $pattern ]]; then
        echo "Expected output to match pattern: $pattern" >&2
        echo "Actual output: $output" >&2
        return 1
    fi
}

# Skip test if command not available
skip_if_no_command() {
    local cmd="$1"
    if ! command -v "$cmd" &> /dev/null; then
        skip "Command not available: $cmd"
    fi
}

# Skip test on macOS if flock not available (required for activation)
skip_if_no_flock_macos() {
    if [[ "$OSTYPE" == "darwin"* ]] && ! command -v flock &> /dev/null; then
        skip "flock not available on macOS (activation requires file locking)"
    fi
}

# Mock git for tests that need it
setup_git_mock() {
    mkdir -p "$TEST_REPO_DIR/.git"
    git config --global user.name "Test User" 2>/dev/null || true
    git config --global user.email "test@example.com" 2>/dev/null || true
}

# Load validation library for testing
load_validation_lib() {
    if [ -f "$TEST_REPO_DIR/base/scripts/lib/validation.sh" ]; then
        source "$TEST_REPO_DIR/base/scripts/lib/validation.sh"
    elif [ -f "$BATS_TEST_DIRNAME/../base/scripts/lib/validation.sh" ]; then
        source "$BATS_TEST_DIRNAME/../base/scripts/lib/validation.sh"
    fi
}

# Create symlink for testing symlink attacks
create_test_symlink() {
    local target="$1"
    local link="$2"
    ln -s "$target" "$link"
}

# Generate random profile name
random_profile_name() {
    echo "test-profile-$(date +%s)-$$"
}

# Export functions for use in tests
export -f setup_test_env
export -f teardown_test_env
export -f create_test_profile
export -f create_test_profiles
export -f assert_file_exists
export -f assert_file_not_exists
export -f assert_file_contains
export -f assert_file_not_contains
export -f assert_dir_exists
export -f assert_output_contains
export -f assert_output_matches
export -f skip_if_no_command
export -f setup_git_mock
export -f load_validation_lib
export -f create_test_symlink
export -f random_profile_name
