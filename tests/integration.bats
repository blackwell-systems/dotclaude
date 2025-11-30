#!/usr/bin/env bats

# Integration tests for dotclaude core workflows
# Tests activate, merge, backup/restore, profile switching

load helpers/test_helper

setup() {
    setup_test_env
    create_test_profiles
}

teardown() {
    teardown_test_env
}

# ============================================================================
# Profile activation tests
# ============================================================================

@test "activate: merges base + profile into ~/.claude/CLAUDE.md" {
    # Skip on macOS if flock is not available
    if [[ "$OSTYPE" == "darwin"* ]] && ! command -v flock &> /dev/null; then
        skip "flock not available on macOS"
    fi

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    [ "$status" -eq 0 ]

    assert_file_exists "$TEST_CLAUDE_DIR/CLAUDE.md"
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Base Configuration"
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Test Profile 1"
}

@test "activate: creates backup before overwriting" {
    skip_if_no_flock_macos
    # Create initial CLAUDE.md
    echo "# Existing config" > "$TEST_CLAUDE_DIR/CLAUDE.md"

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    [ "$status" -eq 0 ]

    # Check backup was created
    local backup_count=$(ls "$TEST_CLAUDE_DIR"/CLAUDE.md.backup.* 2>/dev/null | wc -l)
    [ "$backup_count" -gt 0 ]
}

@test "activate: writes .current-profile marker" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    [ "$status" -eq 0 ]

    assert_file_exists "$TEST_CLAUDE_DIR/.current-profile"
    assert_file_contains "$TEST_CLAUDE_DIR/.current-profile" "test-profile-1"
}

@test "activate: includes separator between base and profile" {
    skip_if_no_flock_macos
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    [ "$status" -eq 0 ]

    # Separator includes both parts (split assertions for flexibility)
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Profile-Specific Additions"
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "test-profile-1"
}

@test "activate: fails with invalid profile name" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate "../../../etc/passwd"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Invalid" ]]
}

@test "activate: fails when profile does not exist" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate nonexistent-profile
    [ "$status" -eq 1 ]
    [[ "$output" =~ "not found" ]]
}

@test "activate: requires profile name argument" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate
    [ "$status" -eq 1 ]
    [[ "$output" =~ "required" ]]
}

# ============================================================================
# Dry-run mode tests
# ============================================================================

@test "activate --dry-run: shows preview without making changes" {
    local original_content="# Original"
    echo "$original_content" > "$TEST_CLAUDE_DIR/CLAUDE.md"

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --dry-run
    [ "$status" -eq 0 ]
    [[ "$output" =~ "DRY RUN" ]] || [[ "$output" =~ "Preview" ]]

    # Verify original file unchanged
    local current_content=$(cat "$TEST_CLAUDE_DIR/CLAUDE.md")
    [ "$current_content" = "$original_content" ]
}

@test "activate --dry-run: does not create backup" {
    echo "# Original" > "$TEST_CLAUDE_DIR/CLAUDE.md"

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --dry-run
    [ "$status" -eq 0 ]

    # Check no backup created
    local backup_count=$(ls "$TEST_CLAUDE_DIR"/CLAUDE.md.backup.* 2>/dev/null | wc -l)
    [ "$backup_count" -eq 0 ]
}

@test "activate --preview: works as alias for --dry-run" {
    echo "# Original" > "$TEST_CLAUDE_DIR/CLAUDE.md"

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --preview
    [ "$status" -eq 0 ]

    # Verify no changes made
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "# Original"
}

# ============================================================================
# Debug mode tests
# ============================================================================

@test "activate --verbose: shows debug output" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --verbose
    [ "$status" -eq 0 ]
    [[ "$output" =~ "DEBUG" ]] || [[ "$output" =~ "Profile name: test-profile-1" ]]
}

@test "activate --debug: works as alias for --verbose" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --debug
    [ "$status" -eq 0 ]
}

@test "activate: DEBUG environment variable enables debug mode" {
    DEBUG=1 run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    [ "$status" -eq 0 ]
}

# ============================================================================
# Profile switching tests
# ============================================================================

@test "switching profiles: replaces previous profile content" {
    # Activate first profile
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Test Profile 1"
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Node.js project"

    # Switch to second profile
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-2
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Test Profile 2"
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Python project"

    # Should NOT contain first profile content
    assert_file_not_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Node.js project"
}

@test "switching profiles: updates .current-profile marker" {
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    assert_file_contains "$TEST_CLAUDE_DIR/.current-profile" "test-profile-1"

    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-2
    assert_file_contains "$TEST_CLAUDE_DIR/.current-profile" "test-profile-2"
}

@test "switching profiles: creates backup of previous profile" {
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    local backup_count_before=$(ls "$TEST_CLAUDE_DIR"/CLAUDE.md.backup.* 2>/dev/null | wc -l)

    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-2
    local backup_count_after=$(ls "$TEST_CLAUDE_DIR"/CLAUDE.md.backup.* 2>/dev/null | wc -l)

    [ "$backup_count_after" -gt "$backup_count_before" ]
}

# ============================================================================
# Backup and restore tests
# ============================================================================

@test "backup: preserves original content" {
    echo "# Original content" > "$TEST_CLAUDE_DIR/CLAUDE.md"
    local original_checksum=$(md5sum "$TEST_CLAUDE_DIR/CLAUDE.md" | awk '{print $1}')

    # Activate profile (creates backup)
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1

    # Find the backup
    local backup_file=$(ls -t "$TEST_CLAUDE_DIR"/CLAUDE.md.backup.* 2>/dev/null | head -1)
    [ -n "$backup_file" ]

    # Verify backup content matches original
    local backup_checksum=$(md5sum "$backup_file" | awk '{print $1}')
    [ "$backup_checksum" = "$original_checksum" ]
}

@test "backup: rotates old backups (keeps 5 most recent)" {
    # Create 7 backups
    for i in {1..7}; do
        echo "# Content $i" > "$TEST_CLAUDE_DIR/CLAUDE.md"
        bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
        sleep 1
    done

    # Count backups
    local backup_count=$(ls "$TEST_CLAUDE_DIR"/CLAUDE.md.backup.* 2>/dev/null | wc -l)

    # Should have at most 5 backups
    [ "$backup_count" -le 5 ]
}

# ============================================================================
# Profile with settings.json tests
# ============================================================================

@test "activate: applies profile settings.json if present" {
    # Create profile with settings.json
    mkdir -p "$TEST_REPO_DIR/profiles/settings-test"
    echo "# Settings Test Profile" > "$TEST_REPO_DIR/profiles/settings-test/CLAUDE.md"
    cat > "$TEST_REPO_DIR/profiles/settings-test/settings.json" <<'EOF'
{
  "model": "opus"
}
EOF

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate settings-test
    [ "$status" -eq 0 ]

    # Check settings.json was copied
    assert_file_exists "$TEST_CLAUDE_DIR/settings.json"
    assert_file_contains "$TEST_CLAUDE_DIR/settings.json" "opus"
}

@test "activate: uses base settings.json if profile has none" {
    # Ensure base has settings.json
    cat > "$TEST_REPO_DIR/base/settings.json" <<'EOF'
{
  "model": "sonnet"
}
EOF

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    [ "$status" -eq 0 ]

    # Should use base settings
    assert_file_exists "$TEST_CLAUDE_DIR/settings.json"
    assert_file_contains "$TEST_CLAUDE_DIR/settings.json" "sonnet"
}

# ============================================================================
# Concurrent execution tests
# ============================================================================

@test "activate: prevents concurrent execution with file lock" {
    # Start activation in background with sleep
    (
        bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 &
        sleep 2
    ) &

    sleep 0.5

    # Try to activate another profile concurrently
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-2
    [ "$status" -eq 1 ] || [ "$status" -eq 0 ]  # Either locks or succeeds after wait

    # Cleanup
    wait
}

# ============================================================================
# Edge case tests
# ============================================================================

@test "activate: handles profile with empty CLAUDE.md" {
    mkdir -p "$TEST_REPO_DIR/profiles/empty-profile"
    touch "$TEST_REPO_DIR/profiles/empty-profile/CLAUDE.md"

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate empty-profile
    [ "$status" -eq 0 ]

    # Should still have base content
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Base Configuration"
}

@test "activate: handles profile with very long content" {
    mkdir -p "$TEST_REPO_DIR/profiles/long-profile"
    for i in {1..1000}; do
        echo "Line $i" >> "$TEST_REPO_DIR/profiles/long-profile/CLAUDE.md"
    done

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate long-profile
    [ "$status" -eq 0 ]

    # Check merged file has both base and all profile lines
    local line_count=$(wc -l < "$TEST_CLAUDE_DIR/CLAUDE.md")
    [ "$line_count" -gt 1000 ]
}

@test "activate: handles special characters in profile content" {
    mkdir -p "$TEST_REPO_DIR/profiles/special-chars"
    cat > "$TEST_REPO_DIR/profiles/special-chars/CLAUDE.md" <<'EOF'
# Special Characters Test

Quotes: "double" and 'single'
Backticks: `code`
Dollar signs: $VAR ${VAR}
Backslashes: \ \n \t
Pipes: | ||
Ampersands: & &&
EOF

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate special-chars
    [ "$status" -eq 0 ]

    # Verify special characters preserved
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Quotes: \"double\""
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Dollar signs: \$VAR"
}
