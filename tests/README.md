# dotclaude Test Suite

Comprehensive test suite for dotclaude.

## Overview

### Go Tests (Primary)

The Go implementation includes unit tests in the `internal/` directories. These are the **primary tests** run in CI.

```bash
# Run all Go tests
go test ./...

# Run with coverage
go test ./... -cover

# Run with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/profile/...
go test ./internal/cli/...
go test ./internal/hooks/...
```

#### What Go Tests Cover

- **Profile validation**: Name validation, path traversal prevention, injection attacks
- **Profile management**: Activation, creation, deletion, listing
- **Backup system**: Automatic backup, rotation, restore
- **Settings handling**: JSON merging, configuration
- **Hook execution**: Pre/post tool hooks, session hooks
- **CLI commands**: All command handlers, flags, output formatting

### BATS Tests (Legacy - Not Run in CI)

⚠️ **Note**: The BATS tests in this directory are legacy tests for the archived shell implementation. They are **not executed** in CI since the shell scripts they test have been archived to `archive/`.

The shell functions like `validate_profile_name`, `validate_directory`, etc. no longer exist. The equivalent functionality is now implemented in Go and tested by Go unit tests.

The BATS test files remain in the repository for historical reference and as documentation of the expected behaviors, but they will fail if run directly.

## Running Tests

### Go Tests (Recommended)

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test ./... -cover

# Run with race detector
go test -race ./...

# Run specific package
go test -v ./internal/profile/...
```

### Prerequisites for BATS (Legacy)

If you need to run the legacy BATS tests for reference:

```bash
# Using npm
npm install -g bats

# Or using Homebrew (macOS)
brew install bats-core

# Or using apt (Ubuntu/Debian)
sudo apt-get install bats
```

⚠️ **Note**: BATS tests will fail because they test archived shell functions.

## Test Structure

```
# Go Tests (Primary - Run in CI)
internal/
├── profile/
│   └── *_test.go          # Profile management tests
├── cli/
│   └── *_test.go          # CLI command tests
└── hooks/
    └── *_test.go          # Hook system tests

# BATS Tests (Legacy - Not Run in CI)
tests/
├── security.bats           # Legacy security validation tests
├── integration.bats        # Legacy workflow tests
├── commands.bats           # Legacy command tests
├── helpers/
│   └── test_helper.bash   # Setup/teardown utilities
└── fixtures/              # Test data (empty - generated dynamically)
```

## Test Categories

### Security Tests (`security.bats`)

Tests all validation functions (originally in shell, now implemented in Go):

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

**Example (legacy BATS test):**
```bash
@test "activate: merges base + profile into ~/.claude/CLAUDE.md" {
    run dotclaude activate test-profile-1
    [ "$status" -eq 0 ]

    assert_file_exists "$TEST_CLAUDE_DIR/CLAUDE.md"
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Base Configuration"
    assert_file_contains "$TEST_CLAUDE_DIR/CLAUDE.md" "Test Profile 1"
}
```

**Go Test Example:**
```go
func TestActivateProfile(t *testing.T) {
    m := profile.NewManager(testRepoDir, testClaudeDir)
    err := m.Activate("test-profile-1", false)
    assert.NoError(t, err)

    content, _ := os.ReadFile(filepath.Join(testClaudeDir, "CLAUDE.md"))
    assert.Contains(t, string(content), "Base Configuration")
    assert.Contains(t, string(content), "Test Profile 1")
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

**Example (legacy BATS test):**
```bash
@test "list: shows all available profiles" {
    run dotclaude list
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

1. **Go Tests** (Ubuntu + macOS)
   - Unit tests for all packages
   - Race condition detection
   - Coverage reporting

2. **Go & Shell Linting**
   - `go vet` for Go code analysis
   - `gofmt` for Go code formatting
   - `shellcheck` for install script and test helpers

3. **Installation Test**
   - Install script execution
   - Profile creation
   - Profile activation

4. **Coverage Report**
   - Go test coverage percentage
   - Test package summary

**Note**: BATS tests are not run in CI as they test the archived shell implementation.

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

### Add Go Test (Recommended)

```go
// In internal/profile/manager_test.go

func TestNewFeature(t *testing.T) {
    m := profile.NewManager(testRepoDir, testClaudeDir)

    // Arrange
    err := m.Activate("test-profile-1", false)
    require.NoError(t, err)

    // Act
    result, err := m.SomeNewMethod()

    // Assert
    require.NoError(t, err)
    assert.Contains(t, result, "expected-content")
}
```

### Add Legacy BATS Test

```bash
# In tests/integration.bats

setup() {
    setup_test_env
    create_test_profiles
}

@test "new workflow: description" {
    # Arrange
    dotclaude activate test-profile-1

    # Act
    run dotclaude some-command

    # Assert
    [ "$status" -eq 0 ]
    assert_file_contains "$TEST_CLAUDE_DIR/some-file" "expected-content"
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
# For Go binary
chmod +x bin/dotclaude

# For legacy shell tests (if running)
chmod +x archive/dotclaude-shell
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
