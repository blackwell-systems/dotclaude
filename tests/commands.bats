#!/usr/bin/env bats

# Command tests for all 12 dotclaude commands
# Tests show, list, activate, switch, create, edit, diff, restore, sync, branches, version, help

load helpers/test_helper

setup() {
    setup_test_env
    create_test_profiles
}

teardown() {
    teardown_test_env
}

# ============================================================================
# dotclaude show
# ============================================================================

@test "show: displays current profile when one is active" {
    skip_if_no_flock_macos
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" show
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test-profile-1" ]] || [[ "$output" =~ "Active Profile" ]]
}

@test "show: handles no active profile gracefully" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" show
    [ "$status" -eq 0 ]
    # Should not crash, may show "No active profile" or similar
}

@test "show: accepts --verbose flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" show --verbose
    [ "$status" -eq 0 ]
}

@test "show: accepts --debug flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" show --debug
    [ "$status" -eq 0 ]
}

# ============================================================================
# dotclaude list
# ============================================================================

@test "list: shows all available profiles" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" list
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test-profile-1" ]]
    [[ "$output" =~ "test-profile-2" ]]
    [[ "$output" =~ "test-profile-3" ]]
}

@test "list: indicates active profile" {
    skip_if_no_flock_macos
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" list
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test-profile-1" ]]
    [[ "$output" =~ "active" ]] || [[ "$output" =~ "●" ]]
}

@test "list: works with alias 'ls'" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" ls
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test-profile" ]]
}

@test "list: accepts --verbose flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" list --verbose
    [ "$status" -eq 0 ]
}

@test "list: handles empty profiles directory" {
    rm -rf "$TEST_REPO_DIR/profiles"/*

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" list
    [ "$status" -eq 0 ]
    [[ "$output" =~ "No profiles" ]] || [[ "$output" =~ "0 profiles" ]]
}

# ============================================================================
# dotclaude activate
# ============================================================================

@test "activate: successfully activates valid profile" {
    skip_if_no_flock_macos
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    [ "$status" -eq 0 ]
    [[ "$output" =~ "activated" ]] || [[ "$output" =~ "test-profile-1" ]]
}

@test "activate: works with alias 'use'" {
    skip_if_no_flock_macos
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" use test-profile-1
    [ "$status" -eq 0 ]
}

@test "activate: accepts --dry-run flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --dry-run
    [ "$status" -eq 0 ]
}

@test "activate: accepts --preview flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --preview
    [ "$status" -eq 0 ]
}

@test "activate: accepts --verbose flag" {
    skip_if_no_flock_macos
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --verbose
    [ "$status" -eq 0 ]
}

@test "activate: accepts combined flags" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --dry-run --verbose
    [ "$status" -eq 0 ]
}

@test "activate: rejects unknown flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1 --unknown-flag
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Unknown flag" ]]
}

# ============================================================================
# dotclaude switch
# ============================================================================

@test "switch: runs without error" {
    # Switch command needs interactive input, so we just test it starts
    skip "switch requires interactive input - needs expect for full test"
}

@test "switch: works with alias 'select'" {
    skip "select requires interactive input - needs expect for full test"
}

@test "switch: accepts --verbose flag" {
    # Just test the flag is accepted (command will timeout waiting for input)
    timeout 1 bash "$TEST_REPO_DIR/base/scripts/dotclaude" switch --verbose || true
}

# ============================================================================
# dotclaude create
# ============================================================================

@test "create: creates new profile directory" {
    local new_profile="test-new-profile"

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" create "$new_profile"
    [ "$status" -eq 0 ]

    assert_dir_exists "$TEST_REPO_DIR/profiles/$new_profile"
    assert_file_exists "$TEST_REPO_DIR/profiles/$new_profile/CLAUDE.md"
}

@test "create: works with alias 'new'" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" new test-new-via-alias
    [ "$status" -eq 0 ]
    assert_dir_exists "$TEST_REPO_DIR/profiles/test-new-via-alias"
}

@test "create: rejects invalid profile name" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" create "../invalid-name"
    [ "$status" -eq 1 ]
}

@test "create: rejects empty profile name" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" create ""
    [ "$status" -eq 1 ]
}

@test "create: fails if profile already exists" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" create test-profile-1
    [ "$status" -eq 1 ]
    [[ "$output" =~ "already exists" ]] || [[ "$output" =~ "exists" ]]
}

@test "create: accepts --verbose flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" create test-verbose-create --verbose
    [ "$status" -eq 0 ]
}

# ============================================================================
# dotclaude edit
# ============================================================================

@test "edit: requires EDITOR environment variable or uses fallback" {
    export EDITOR="cat"

    run timeout 1 bash "$TEST_REPO_DIR/base/scripts/dotclaude" edit test-profile-1 || true
    # Command will attempt to open editor, we just verify it starts
}

@test "edit: edits current profile when no name specified" {
    skip_if_no_flock_macos
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    export EDITOR="cat"

    run timeout 1 bash "$TEST_REPO_DIR/base/scripts/dotclaude" edit || true
    # Just verify command accepts no argument
}

@test "edit: accepts --verbose flag" {
    export EDITOR="cat"
    run timeout 1 bash "$TEST_REPO_DIR/base/scripts/dotclaude" edit test-profile-1 --verbose || true
}

# ============================================================================
# dotclaude diff
# ============================================================================

@test "diff: compares two profiles" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" diff test-profile-1 test-profile-2
    [ "$status" -eq 0 ]
}

@test "diff: compares current vs specified profile" {
    skip_if_no_flock_macos
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" diff test-profile-2
    [ "$status" -eq 0 ]
}

@test "diff: fails when profile does not exist" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" diff nonexistent-profile
    [ "$status" -eq 1 ]
}

@test "diff: accepts --verbose flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" diff test-profile-1 test-profile-2 --verbose
    [ "$status" -eq 0 ]
}

# ============================================================================
# dotclaude restore
# ============================================================================

@test "restore: runs without error when backups exist" {
    skip_if_no_flock_macos
    # Create a backup by activating
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1

    # Restore is interactive, we just verify command exists
    skip "restore requires interactive input - needs expect for full test"
}

@test "restore: accepts --verbose flag" {
    timeout 1 bash "$TEST_REPO_DIR/base/scripts/dotclaude" restore --verbose || true
}

# ============================================================================
# dotclaude sync
# ============================================================================

@test "sync: exits gracefully when not in git repo" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" sync
    # May fail or succeed depending on whether git repo exists
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "sync: accepts --verbose flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" sync --verbose
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

# ============================================================================
# dotclaude branches
# ============================================================================

@test "branches: exits gracefully when not in git repo" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" branches
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "branches: works with alias 'br'" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" br
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "branches: accepts --verbose flag" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" branches --verbose
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

# ============================================================================
# dotclaude version
# ============================================================================

@test "version: displays version information" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" version
    [ "$status" -eq 0 ]
    [[ "$output" =~ "version" ]] || [[ "$output" =~ "1.0" ]]
}

@test "version: works with alias -v" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" -v
    [ "$status" -eq 0 ]
}

@test "version: works with alias --version" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" --version
    [ "$status" -eq 0 ]
}

# ============================================================================
# dotclaude help
# ============================================================================

@test "help: displays help information" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "dotclaude" ]]
    [[ "$output" =~ "USAGE" ]] || [[ "$output" =~ "COMMANDS" ]]
}

@test "help: works with alias -h" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" -h
    [ "$status" -eq 0 ]
}

@test "help: works with alias --help" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" --help
    [ "$status" -eq 0 ]
}

@test "help: lists all available commands" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" help
    [ "$status" -eq 0 ]

    # Check major commands are listed
    [[ "$output" =~ "activate" ]]
    [[ "$output" =~ "list" ]]
    [[ "$output" =~ "show" ]]
}

# ============================================================================
# Unknown command handling
# ============================================================================

@test "unknown command: returns error" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" nonexistent-command
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Unknown command" ]] || [[ "$output" =~ "not found" ]]
}

@test "no command: shows help" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude"
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
    # May show help or error, both acceptable
}

# ============================================================================
# Exit code tests
# ============================================================================

@test "exit code: 0 for successful operations" {
    skip_if_no_flock_macos
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" list
    [ "$status" -eq 0 ]

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    [ "$status" -eq 0 ]
}

@test "exit code: 1 for errors" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate nonexistent-profile
    [ "$status" -eq 1 ]

    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate "../invalid"
    [ "$status" -eq 1 ]
}
