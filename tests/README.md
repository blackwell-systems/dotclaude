# dotclaude Test Suite

Comprehensive test suite for dotclaude using bats-core (Bash Automated Testing System).

## Overview

The test suite contains 145+ tests covering:
- **Security tests** (45 tests): Validation functions, path traversal prevention, injection attacks
- **Integration tests** (42 tests): Profile activation, merging, backup/restore, profile switching
- **Command tests** (58 tests): All 12 commands, flags, exit codes

## Running Tests

### Prerequisites

Install bats-core:

```bash
# Using npm (recommended)
npm install -g bats

# Or using Homebrew (macOS)
brew install bats-core

# Or using apt (Ubuntu/Debian)
sudo apt-get install bats
```

### Run All Tests

```bash
# From repository root
bats tests/*.bats

# Or run individually
bats tests/security.bats
bats tests/integration.bats
bats tests/commands.bats
```

### Run Specific Test

```bash
# Run single test by name
bats tests/security.bats --filter "validate_profile_name: rejects path traversal"

# Run tests matching pattern
bats tests/integration.bats --filter "activate"
```

### Verbose Output

```bash
# Show all test output (including passing tests)
bats tests/security.bats --tap

# Show timing information
bats tests/security.bats --timing
```

## Test Structure

```
tests/
├── security.bats           # Security validation tests (45 tests)
├── integration.bats        # Core workflow tests (42 tests)
├── commands.bats           # Command tests (58 tests)
├── helpers/
│   └── test_helper.bash   # Setup/teardown and utility functions
└── fixtures/              # Test data (empty - generated dynamically)
```

## Test Categories

### Security Tests (`security.bats`)

Tests all validation functions in `base/scripts/lib/validation.sh`:

- **Path traversal prevention**: Rejects `../../../etc/passwd`
- **Symlink attack prevention**: Rejects symlinked directories
- **Input validation**: Rejects special characters, spaces, slashes
- **Command injection**: Prevents shell metacharacters
- **Safe directory removal**: Ensures only safe paths removed
- **Disk space checks**: Validates available space
- **File locking**: Prevents concurrent execution
- **Sensitive data detection**: Warns on API keys, passwords, tokens

**Example:**
```bash
@test "validate_profile_name: rejects path traversal with ../" {
    run validate_profile_name "../../../etc/passwd"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "path traversal" ]]
}
```

### Integration Tests (`integration.bats`)

Tests core workflows end-to-end:

- **Profile activation**: Base + profile merge correctly
- **Backup creation**: Original files backed up before overwrite
- **Profile switching**: Previous profile content replaced
- **Dry-run mode**: No files changed in preview mode
- **Debug mode**: Verbose output when requested
- **Settings.json**: Profile settings applied correctly
- **Concurrent execution**: File locking prevents conflicts
- **Edge cases**: Empty profiles, long content, special characters

**Example:**
```bash
@test "activate: merges base + profile into ~/.claude/CLAUDE.md" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1
    [ "$status" -eq 0 ]

    assert_file_exists "$TEST_CLAUDE_DIR/CLAUDE.md"
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Base Configuration"
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Test Profile 1"
}
```

### Command Tests (`commands.bats`)

Tests all 12 commands:

| Command | Tests |
|---------|-------|
| show | Display current profile, handle no profile, flags |
| list | List profiles, show active, empty directory |
| activate | Activate profile, dry-run, flags, validation |
| switch | Interactive switcher, flags |
| create | Create profile, validation, duplicates |
| edit | Open editor, current profile, flags |
| diff | Compare profiles, validation, flags |
| restore | Interactive restore, flags |
| sync | Git sync, non-repo handling, flags |
| branches | Branch status, alias, flags |
| version | Show version, aliases (-v, --version) |
| help | Show help, aliases (-h, --help), command list |

**Example:**
```bash
@test "list: shows all available profiles" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" list
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test-profile-1" ]]
    [[ "$output" =~ "test-profile-2" ]]
}
```

## Test Helpers

The test helper (`helpers/test_helper.bash`) provides:

### Setup/Teardown
- `setup_test_env()`: Creates isolated test environment
- `teardown_test_env()`: Cleans up after tests
- `create_test_profile()`: Creates test profile
- `create_test_profiles()`: Creates multiple test profiles

### Assertions
- `assert_file_exists()`: Verify file exists
- `assert_file_contains()`: Verify file contains pattern
- `assert_dir_exists()`: Verify directory exists
- `assert_output_contains()`: Verify command output contains pattern

### Utilities
- `load_validation_lib()`: Load validation.sh for testing
- `setup_git_mock()`: Mock git for git-dependent tests
- `random_profile_name()`: Generate unique profile name

## CI/CD Integration

Tests run automatically on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Manual workflow dispatch

### GitHub Actions Workflow

`.github/workflows/test.yml` runs:

1. **Test Suite** (Ubuntu + macOS)
   - Security tests
   - Integration tests
   - Command tests

2. **Shell Linting** (shellcheck)
   - All scripts in `base/scripts/`
   - All test helpers

3. **Installation Test**
   - Install script
   - Profile creation
   - Profile activation

4. **Coverage Report**
   - Test count summary
   - Success/failure status

### View Results

- GitHub Actions: Repository → Actions tab
- PR checks: Automatic status on pull requests
- Badges: Add to README.md (optional)

## Writing New Tests

### Add Security Test

```bash
# In tests/security.bats

@test "validate_profile_name: rejects new attack vector" {
    run validate_profile_name "malicious-input"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "error message" ]]
}
```

### Add Integration Test

```bash
# In tests/integration.bats

setup() {
    setup_test_env
    create_test_profiles
}

@test "new workflow: description" {
    # Arrange
    bash "$TEST_REPO_DIR/base/scripts/dotclaude" activate test-profile-1

    # Act
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" some-command

    # Assert
    [ "$status" -eq 0 ]
    assert_file_contains "$TEST_CLAUDE_DIR/some-file" "expected-content"
}
```

### Add Command Test

```bash
# In tests/commands.bats

@test "new-command: does expected thing" {
    run bash "$TEST_REPO_DIR/base/scripts/dotclaude" new-command
    [ "$status" -eq 0 ]
    [[ "$output" =~ "success message" ]]
}
```

## Debugging Tests

### Run Single Test

```bash
bats tests/security.bats --filter "validate_profile_name: rejects path traversal"
```

### Show All Output

```bash
bats tests/security.bats --tap
```

### Add Debug Output

```bash
@test "debugging example" {
    echo "Debug: variable=$variable" >&3
    run some_command
    echo "Debug: status=$status" >&3
    echo "Debug: output=$output" >&3
    [ "$status" -eq 0 ]
}
```

### Inspect Test Environment

```bash
@test "inspect test env" {
    echo "TEST_REPO_DIR=$TEST_REPO_DIR" >&3
    echo "TEST_CLAUDE_DIR=$TEST_CLAUDE_DIR" >&3
    ls -la "$TEST_REPO_DIR" >&3
}
```

## Test Coverage Goals

- ✅ **Security**: 100% of validation functions tested
- ✅ **Core workflows**: Activate, merge, backup, restore
- ✅ **All commands**: 12/12 commands tested
- ✅ **All flags**: --verbose, --debug, --dry-run, --preview
- ✅ **Exit codes**: Success (0) and error (1) cases
- 🟡 **Interactive commands**: Needs expect for full coverage
- 🟡 **Git integration**: Needs real git repo for full coverage

## Performance

Test suite runs in ~5-10 seconds:
- Security tests: ~2 seconds
- Integration tests: ~3-5 seconds
- Command tests: ~2-3 seconds

## Troubleshooting

### bats: command not found

```bash
npm install -g bats
```

### Tests fail with "permission denied"

```bash
chmod +x base/scripts/dotclaude
chmod +x base/scripts/lib/*.sh
```

### Tests leave temp files

Tests clean up automatically in `teardown()`. If interrupted:

```bash
rm -rf /tmp/dotclaude-test-*
```

### CI tests pass but local tests fail

Ensure environment matches CI:
```bash
export DOTCLAUDE_REPO_DIR="$(pwd)"
bats tests/*.bats
```

## Contributing

When adding new features:
1. Write tests first (TDD)
2. Ensure all tests pass: `bats tests/*.bats`
3. Add documentation to this README
4. Submit PR with tests included

## Resources

- [bats-core documentation](https://bats-core.readthedocs.io/)
- [Bash testing best practices](https://github.com/bats-core/bats-core#writing-tests)
- [GitHub Actions workflow syntax](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions)
