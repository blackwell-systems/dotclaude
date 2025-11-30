#!/usr/bin/env bats

# Security tests for dotclaude validation functions
# Tests path traversal prevention, symlink attacks, input validation

load helpers/test_helper

setup() {
    setup_test_env
    load_validation_lib
}

teardown() {
    teardown_test_env
}

# ============================================================================
# validate_profile_name tests
# ============================================================================

@test "validate_profile_name: accepts valid alphanumeric name" {
    run validate_profile_name "myproject"
    [ "$status" -eq 0 ]
}

@test "validate_profile_name: accepts name with hyphens" {
    run validate_profile_name "my-project"
    [ "$status" -eq 0 ]
}

@test "validate_profile_name: accepts name with underscores" {
    run validate_profile_name "my_project"
    [ "$status" -eq 0 ]
}

@test "validate_profile_name: accepts mixed alphanumeric with separators" {
    run validate_profile_name "my-project_2024"
    [ "$status" -eq 0 ]
}

@test "validate_profile_name: rejects path traversal with ../" {
    run validate_profile_name "../../../etc/passwd"
    [ "$status" -eq 1 ]
    # Accept either "path traversal" or "Invalid" error message
    [[ "$output" =~ "path traversal" ]] || [[ "$output" =~ "Invalid" ]]
}

@test "validate_profile_name: rejects path traversal with .." {
    run validate_profile_name ".."
    [ "$status" -eq 1 ]
}

@test "validate_profile_name: rejects absolute path" {
    run validate_profile_name "/etc/passwd"
    [ "$status" -eq 1 ]
}

@test "validate_profile_name: rejects name with slashes" {
    run validate_profile_name "my/project"
    [ "$status" -eq 1 ]
}

@test "validate_profile_name: rejects name with spaces" {
    run validate_profile_name "my project"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Invalid" ]]
}

@test "validate_profile_name: rejects empty name" {
    run validate_profile_name ""
    [ "$status" -eq 1 ]
    [[ "$output" =~ "cannot be empty" ]]
}

@test "validate_profile_name: rejects name with special characters" {
    run validate_profile_name "my@project"
    [ "$status" -eq 1 ]
}

@test "validate_profile_name: rejects shell injection attempt" {
    run validate_profile_name "test; rm -rf /"
    [ "$status" -eq 1 ]
}

@test "validate_profile_name: rejects command substitution" {
    run validate_profile_name "test\$(rm -rf /)"
    [ "$status" -eq 1 ]
}

# ============================================================================
# validate_directory tests
# ============================================================================

@test "validate_directory: accepts valid directory" {
    mkdir -p "$TEST_TEMP_DIR/valid-dir"
    run validate_directory "$TEST_TEMP_DIR/valid-dir"
    [ "$status" -eq 0 ]
}

@test "validate_directory: accepts non-existent path" {
    run validate_directory "$TEST_TEMP_DIR/does-not-exist"
    [ "$status" -eq 0 ]
}

@test "validate_directory: rejects file as directory" {
    touch "$TEST_TEMP_DIR/file"
    run validate_directory "$TEST_TEMP_DIR/file"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "not a directory" ]]
}

@test "validate_directory: rejects symlinked directory" {
    mkdir -p "$TEST_TEMP_DIR/real-dir"
    ln -s "$TEST_TEMP_DIR/real-dir" "$TEST_TEMP_DIR/symlink-dir"

    run validate_directory "$TEST_TEMP_DIR/symlink-dir"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "symlink" ]]
}

# ============================================================================
# safe_remove_directory tests
# ============================================================================

@test "safe_remove_directory: removes valid directory" {
    mkdir -p "$TEST_REPO_DIR/profiles/test-profile"
    [ -d "$TEST_REPO_DIR/profiles/test-profile" ]

    run safe_remove_directory "$TEST_REPO_DIR/profiles/test-profile"
    [ "$status" -eq 0 ]
    [ ! -d "$TEST_REPO_DIR/profiles/test-profile" ]
}

@test "safe_remove_directory: succeeds on non-existent directory" {
    run safe_remove_directory "$TEST_TEMP_DIR/does-not-exist"
    [ "$status" -eq 0 ]
}

@test "safe_remove_directory: refuses to remove symlink" {
    mkdir -p "$TEST_TEMP_DIR/real-dir"
    ln -s "$TEST_TEMP_DIR/real-dir" "$TEST_TEMP_DIR/symlink"

    run safe_remove_directory "$TEST_TEMP_DIR/symlink"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Refusing to remove symlink" ]]
}

@test "safe_remove_directory: refuses to remove file" {
    touch "$TEST_TEMP_DIR/file"

    run safe_remove_directory "$TEST_TEMP_DIR/file"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Not a directory" ]]
}

@test "safe_remove_directory: refuses unsafe path outside profiles/" {
    mkdir -p "$TEST_TEMP_DIR/unsafe-dir"

    run safe_remove_directory "$TEST_TEMP_DIR/unsafe-dir"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "outside safe zones" ]]
}

# ============================================================================
# check_disk_space tests
# ============================================================================

@test "check_disk_space: succeeds when sufficient space available" {
    run check_disk_space "$TEST_TEMP_DIR" 1
    [ "$status" -eq 0 ]
}

@test "check_disk_space: fails when insufficient space (unrealistic requirement)" {
    # Request 1TB - should fail on most systems
    run check_disk_space "$TEST_TEMP_DIR" 1073741824
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Insufficient disk space" ]]
}

# ============================================================================
# acquire_lock and release_lock tests
# ============================================================================

@test "acquire_lock: successfully acquires lock" {
    # Skip on macOS if flock is not available
    if [[ "$OSTYPE" == "darwin"* ]] && ! command -v flock &> /dev/null; then
        skip "flock not available on macOS"
    fi

    local lockfile="$TEST_TEMP_DIR/test.lock"

    run acquire_lock "$lockfile" 1
    [ "$status" -eq 0 ]
}

@test "acquire_lock: fails when lock already held" {
    local lockfile="$TEST_TEMP_DIR/test.lock"

    # Acquire lock in subshell
    (
        exec 200>"$lockfile"
        flock -n 200
        sleep 2
    ) &

    sleep 0.5

    # Try to acquire same lock with short timeout
    run acquire_lock "$lockfile" 1
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Another dotclaude operation" ]]
}

# ============================================================================
# check_sensitive_data tests
# ============================================================================

@test "check_sensitive_data: passes clean file" {
    echo "# Configuration file" > "$TEST_TEMP_DIR/config.md"

    run check_sensitive_data "$TEST_TEMP_DIR/config.md"
    [ "$status" -eq 0 ]
}

@test "check_sensitive_data: warns on api_key pattern" {
    echo "api_key=secret123" > "$TEST_TEMP_DIR/config.md"

    run check_sensitive_data "$TEST_TEMP_DIR/config.md"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "may contain sensitive data" ]]
}

@test "check_sensitive_data: warns on password pattern" {
    echo "password: mysecret" > "$TEST_TEMP_DIR/config.md"

    run check_sensitive_data "$TEST_TEMP_DIR/config.md"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "may contain sensitive data" ]]
}

@test "check_sensitive_data: warns on token pattern" {
    echo "AUTH_TOKEN=xyz123" > "$TEST_TEMP_DIR/config.md"

    run check_sensitive_data "$TEST_TEMP_DIR/config.md"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "may contain sensitive data" ]]
}

@test "check_sensitive_data: warns on secret pattern" {
    echo "client_secret: abc" > "$TEST_TEMP_DIR/config.md"

    run check_sensitive_data "$TEST_TEMP_DIR/config.md"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "may contain sensitive data" ]]
}

# ============================================================================
# require_commands tests
# ============================================================================

@test "require_commands: succeeds when commands exist" {
    run require_commands "bash" "cat" "echo"
    [ "$status" -eq 0 ]
}

@test "require_commands: fails when command missing" {
    run require_commands "bash" "nonexistent-command-xyz"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "nonexistent-command-xyz" ]]
}

# ============================================================================
# validate_repo_dir tests
# ============================================================================

@test "validate_repo_dir: accepts valid repo structure" {
    run validate_repo_dir "$TEST_REPO_DIR"
    [ "$status" -eq 0 ]
}

@test "validate_repo_dir: rejects empty path" {
    run validate_repo_dir ""
    [ "$status" -eq 1 ]
    [[ "$output" =~ "not set" ]]
}

@test "validate_repo_dir: rejects non-existent directory" {
    run validate_repo_dir "/nonexistent/path"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "does not exist" ]]
}

@test "validate_repo_dir: rejects directory without base/" {
    local invalid_repo="$TEST_TEMP_DIR/invalid-repo"
    mkdir -p "$invalid_repo/profiles"

    run validate_repo_dir "$invalid_repo"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Invalid repository structure" ]]
}

@test "validate_repo_dir: rejects directory without profiles/" {
    local invalid_repo="$TEST_TEMP_DIR/invalid-repo"
    mkdir -p "$invalid_repo/base"

    run validate_repo_dir "$invalid_repo"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Invalid repository structure" ]]
}

# ============================================================================
# sanitize_output tests
# ============================================================================

@test "sanitize_output: removes ANSI escape codes" {
    local input=$'\033[0;31mRed text\033[0m'

    run sanitize_output "$input"
    [ "$status" -eq 0 ]
    [[ "$output" == "Red text" ]]
}

@test "sanitize_output: removes control characters" {
    local input=$'Text with\x00null\x01chars'

    run sanitize_output "$input"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Text with" ]]
    # Check that output doesn't contain null bytes (platform-agnostic check)
    [ "${#output}" -lt "${#input}" ] || [ "$output" = "Text with" ]
}
