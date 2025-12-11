# Changelog

All notable changes to dotclaude will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0-beta.3] - 2025-12-11

### Summary

Bug fix release addressing issues discovered during comprehensive Go migration audit, plus addition of comprehensive Go unit tests.

### Added

**Go Unit Tests:**
- Profile package tests (83.3% coverage)
  - ValidateProfileName - 17 test cases for name validation
  - Manager operations - NewManager, ListProfiles, GetActiveProfile
  - Create/Delete/Activate operations with error paths
  - Backup and restore functionality tests
  - Edge cases and error path coverage
- CLI package tests (30.1% coverage)
  - Command registration and structure tests
  - Error handling for invalid inputs
  - Flag verification for all commands
- Total: 80+ new Go test cases

**switch command:**
- New interactive profile selector command
- Displays numbered list of profiles
- `select` alias for `switch` command
- Shows current active profile with indicator

**New command aliases:**
- `branches` alias for check-branches command
- `br` alias for check-branches command

**edit command enhancement:**
- Now works without argument (uses active profile)
- If no active profile, shows helpful error message

### Fixed

**Bug Fixes:**
- `check-branches` command now has `branches` and `br` aliases for discoverability
- `edit` command now supports editing active profile without specifying name
- Path construction in delete.go now uses `filepath.Join()` for cross-platform safety
- File permissions standardized to 0644 for all backup and config files

### Changed

**Code Quality:**
- Removed unused `getCommitCountForSync()` function from sync.go
- Standardized file permissions: backups and restored files both use 0644

### Documentation

- Complete rewrite of ARCHITECTURE.md for Go implementation
- Added Go package structure documentation
- Updated security architecture diagrams
- Added backend selection documentation

## [1.0.0-beta.2] - 2025-12-10

### Summary

Test-driven feature completion release. Comprehensive test suite (122 tests) exposed missing features and infrastructure gaps during migration. All gaps now closed, achieving 100% feature parity with CI/CD validation.

### Added

**diff command:**
- Compare CLAUDE.md files between two profiles
- Single argument compares with active profile
- Uses system diff with unified format (-u)
- Full error handling for non-existent profiles

**Missing flags:**
- `--dry-run` / `--preview` on activate - Preview changes without applying
- `--verbose` / `--debug` on activate - Show debug output during activation
- `--debug` on show command - Display internal paths and settings
- `--version` on root command - Cobra built-in version support

**New aliases:**
- `new` alias for create command

### Fixed

**CI/CD Stabilization:**
- Added Go 1.23 setup to all test jobs (ubuntu, macos, test-install)
- Binary builds before tests run
- Git configured for profile creation tests
- DOTCLAUDE_REPO_DIR exported in all test-install steps

**Test Infrastructure:**
- Go binary now copied to test temp directories
- Git config runs AFTER HOME export (critical ordering fix)
- Git commit failure made non-fatal in profile creation
- Tests updated to match Cobra's error message format (lowercase)

**Profile Creation:**
- Git commit no longer blocks profile creation if user not configured
- Profile successfully created even if git commit fails
- Graceful degradation for missing git config
- Best-effort git initialization

### Technical Details

**Environment Variable Ordering:**
Git's `--global` flag writes to `$HOME/.gitconfig`. Tests were setting
git config before changing HOME, causing config to be written to the
wrong location. Fixed by setting HOME first, then configuring git.

**Error Message Format:**
Tests expected shell-specific error messages but Cobra uses different
formatting. Updated tests to accept Cobra's format:
- "unknown flag" (lowercase)
- "Usage" / "Commands" (mixed case)

**Test Coverage Impact:**
The 122-test suite acted as a specification, exposing 9 missing features:
1. --dry-run flag
2. --preview flag
3. --verbose flag on activate
4. --debug flag on activate + show
5. --version flag on root
6. diff command
7. 'new' alias
8. Binary not in test directories
9. Git config ordering issue

### Changed

**Wrapper Behavior:**
- Default backend changed from `auto` to `go` in v1.0.0-beta.1
- Now uses Go implementation by default
- Shell available via `DOTCLAUDE_BACKEND=shell`

**Test Strategy:**
- CI tests run against Go implementation
- Validates feature parity automatically
- Catches regressions immediately

### Migration Progress

**Phase 6: ✅ COMPLETE**
- 11/11 commands implemented
- All flags present
- All aliases working
- CI/CD stable
- 122/122 tests passing

**Status:** Ready for production use with validation period

### Commits

- 8c65667 fix: Run git config after HOME export in test helper
- 610e39d fix: Make git commit failure non-fatal in profile creation
- 2628538 fix: Configure git in test helper for profile creation
- 1d1ee5c fix: Update tests to match Cobra's error message format
- 87a5f75 feat: Add missing features revealed by testing
- 3b3ac5f fix: Copy Go binary to test directories and export DOTCLAUDE_REPO_DIR
- ba3e4d9 feat: Add missing flags for feature parity with shell version
- ad4e929 ci: Add Go build support to test workflows

## [1.0.0-beta.1] - 2025-12-10

### Summary

Initial beta release with Go migration complete (10 commands). Wrapper
defaults to Go backend, shell fallback available.

## [1.0.0-alpha.5] - 2025-12-10

### Added
- **sync command** - FINAL COMMAND! Migration complete! 🎉
  - Syncs feature branches with main using rebase or merge
  - When on main: Lists feature branches that are behind
  - When on feature branch: Offers rebase or merge options
  - Checks for uncommitted changes before syncing
  - Interactive confirmation for pushing changes
  - Configurable base branch with `--base` flag
  - Helpful guidance for conflict resolution

### Progress
- **10/10 commands complete (100%)**
- **Go migration COMPLETE!**
- All shell commands successfully ported to Go
- Full functional parity achieved

### Migration Summary
After auditing the actual shell implementation:
- Initial plan had 13 commands (included non-existent ones)
- Actual shell version has 10 commands
- All 10 commands now implemented in Go
- **backup**: Not a command (automatic via activate)
- **deactivate**: Never existed in shell
- **feature-branch**: Never existed in shell (only `branches`/check-branches)

**Final Stats:**
- 10 commands implemented
- ~1,400 lines of Go code
- 7 hours total migration time
- 100% parity with shell version
- Strangler fig pattern successfully applied

**Branch:** `go-migration` (ready for testing and merge to main)

## [1.0.0-alpha.4] - 2025-12-10

### Added
- **check-branches command** - Check which branches are behind main
  - Lists all feature branches (excludes main/master)
  - Shows commits ahead/behind for each branch
  - Fetches from origin automatically
  - Configurable base branch with `--base` flag
  - Clean output showing only divergent branches

### Progress
- **9/12 commands complete (75%)**
- **check-branches** now functional
- Only 3 commands remaining (deactivate, sync, feature-branch)

**Migration Timeline:** 7 hours completed, 6-9 hours estimated remaining (~1 weekend)

**Branch:** `go-migration`

## [1.0.0-alpha.3] - 2025-12-10

### Added
- **restore command** - Interactive backup restoration
  - Lists all available backups sorted by modification time
  - Groups backups by type (CLAUDE.md vs settings.json)
  - Interactive selection with cancel option ('q')
  - Confirms overwrite before restoring
  - Creates backup of current file before restoring
  - Updates .current-profile marker when restoring CLAUDE.md
  - Handles missing backups with helpful message

- **Backup types in profile package**
  - `Backup` struct with Path, Filename, Timestamp, Size, Type
  - `ListBackups()` - Find and sort all backup files
  - `parseBackup()` - Extract backup metadata
  - `Restore()` - Restore with current file backup
  - `updateProfileFromCLAUDE()` - Auto-detect profile from content

### Changed
- **Migration plan clarification**: "backup" command removed from plan
  - Backups are automatic (created by activate command)
  - Only restore command needed for user interaction
  - Updated from 13 → 12 total commands

### Progress
- **8/12 commands complete (67%)**
- **restore** now functional with interactive UI
- Only 4 commands remaining (deactivate, sync, check-branches, feature-branch)

**Migration Timeline:** 6 hours completed, 7-11 hours estimated remaining (~1 weekend)

**Branch:** `go-migration`

## [1.0.0-alpha.2] - 2025-12-10

### Added
- **activate command** - Most complex command complete
  - Merges base + profile CLAUDE.md with separator header
  - Applies settings.json (profile or base fallback)
  - Creates timestamped backups on profile switch
  - Detects re-activation and updates in place
  - Keeps only 5 most recent backups per file
  - Updates .current-profile state file
  - Creates Claude directory if missing

- **Container Testing Environment** (`Dockerfile.go-test`)
  - Go 1.23 Alpine with full toolchain
  - Automated test script at `/root/test-activate.sh`
  - Forced Go backend for testing
  - Helper script: `scripts/test-in-container.sh`

### Progress
- **7/13 commands complete (54%)**
- **activate** now functional with full backup/merge logic
- **Profile Management** expanded with activation functions:
  - `Activate()` - Main activation with backup management
  - `mergeCLAUDEmd()` - Merge base + profile with separator
  - `applySettings()` - Apply profile or base settings
  - `backupFile()` - Timestamped backups with cleanup
  - `cleanupBackups()` - Keep only N most recent backups

### Changed
- State file standardized to `.current-profile` (was `.dotclaude-active`)
- Go module requirement changed from `go 1.25.5` to `go 1.23` for compatibility

### Testing
- ✅ First activation creates backups
- ✅ Re-activation updates in place
- ✅ Delete command prevents removing active profile
- ✅ Backup cleanup keeps only 5 most recent
- ✅ Settings fallback to base if profile has none

**Migration Timeline:** 5 hours completed, 8-13 hours estimated remaining (1-2 weekends)

**Branch:** `go-migration`

## [1.0.0-alpha.1] - 2025-12-10

### Added
- **Go Implementation** - Started migration from shell to Go using strangler fig pattern
  - Cobra CLI framework with full command structure
  - Smart wrapper (`base/scripts/dotclaude`) routes commands to Go or shell implementation
  - Environment variable control: `DOTCLAUDE_BACKEND=go|shell|auto`
  - Shell version preserved as `base/scripts/dotclaude-shell` for reference
  - See [GO-MIGRATION.md](GO-MIGRATION.md) for full migration plan and progress

- **Go Commands Implemented (6/13 - 46% complete)**
  - `version` - Display version and build information
  - `list` - List all profiles with active indicator
  - `show` - Show active profile information
  - `create` - Create new profile from template with git initialization
  - `delete` - Delete profile with confirmation prompt and safety checks
  - `edit` - Open CLAUDE.md or settings.json in $EDITOR

- **Profile Management Foundation**
  - Profile struct with metadata (Name, Path, IsActive, LastModified)
  - Manager with RepoDir, ProfilesDir, ClaudeDir, StateFile
  - Core operations: ListProfiles(), GetActiveProfile(), ProfileExists(), ValidateProfileName()
  - File operations: Create(), Delete(), copyDir(), copyFile(), initGitRepo()

- **Build System**
  - Makefile with build, test, clean, install targets
  - Go module: `github.com/blackwell-systems/dotclaude`
  - Dependencies: `github.com/spf13/cobra v1.10.2`

### Changed
- Original shell implementation renamed to `dotclaude-shell`
- New wrapper script enables transparent routing between implementations
- Profile creation now initializes git repository (Go version only)

### In Progress
- **Remaining Commands (7/13 - 54%)**
  - activate (HIGH priority - most complex)
  - deactivate
  - backup
  - restore
  - sync
  - check-branches
  - feature-branch

**Migration Timeline:** 3 hours completed, 10-16 hours estimated remaining (2-3 weekends)

**Branch:** `go-migration`

**Commits:**
- 8db41e1 - Foundation + version command
- ba96f9c - list and show commands
- 2d32d5c - create command
- e17314c - delete and edit commands

## [0.5.1] - 2025-12-02

### Fixed
- ANSI color codes now render properly in help/version output
  - Added `-e` flag to all echo statements using color variables
  - Set `TERM=xterm-256color` in Dockerfile.lite
  - Commands like `dotclaude help` now display with proper colors and formatting

## [0.5.0] - 2025-12-02

### Changed
- **BREAKING (Install Only):** Profile creation now mandatory during installation
  - Interactive install prompts for profile name (validated, cannot be empty)
  - Non-interactive install shows instructions for manual profile creation
  - Ensures every user has a working profile after installation
  - **Note:** Existing installations and profiles are unaffected - only impacts new installs

- **Removed fallback template** from `dotclaude create` command
  - Now fails loudly with helpful error if `examples/sample-profile` is missing
  - Provides clear diagnostics and reinstall instructions
  - Prevents silent failures with incomplete templates
  - Ensures all profiles start with comprehensive 250+ line template

### Improved
- Installation "Next Steps" now context-aware
  - Shows different messages if profile was created vs skipped (non-interactive)
  - Guides users to customize (not create) if profile already exists
  - Clearer post-install experience

**Upgrade Note:** Existing users on v0.4.0 can upgrade seamlessly. These changes only affect new installations. Your existing profiles continue to work without modification.

## [0.4.0] - 2025-12-02

### Added
- **curl | bash Install Support** - One-line install now works seamlessly
  - Auto-detects if running from repo or via curl
  - Clones to ~/code/dotclaude if needed
  - Re-execs from cloned location automatically
  - Fully backward compatible with manual git clone
- **Automatic Shell Configuration** - DOTCLAUDE_REPO_DIR added to shell RC automatically
  - Detects bash or zsh and updates appropriate RC file
  - Idempotent (checks if already present)
  - No manual configuration required
- **Install Validation Checks** - Post-install health validation
  - Verifies dotclaude CLI in PATH
  - Confirms management scripts installed
  - Checks repository accessibility
  - Validates sample profile availability
- **Profile Template Scaffolding** - `dotclaude create` now uses comprehensive template
  - Copies from examples/sample-profile automatically
  - Provides 200+ lines of examples and best practices
  - Fallback to basic template if sample-profile missing
  - Customizes profile name in template
- **Dockerfile.lite** - Lightweight Alpine container for safe testing
  - Pre-installs dotclaude CLI
  - Shows welcome message with command suggestions
  - Auto-deletes on exit (--rm flag)
  - 30-second trust verification before real install
- **ONBOARDING-AUDIT.md** - Comprehensive new user experience analysis
  - Documents all onboarding issues found
  - Testing notes and recommendations
  - Tracks resolution status

### Changed
- **Improved Post-Install Messaging** - Clear numbered "Next Steps" guide
  - Shows `dotclaude create` workflow (not manual cp -r)
  - Links to documentation
  - Removed confusing "optional" DOTCLAUDE_REPO_DIR messaging
- **Unified Install Instructions** - Consistent across all documentation
  - README uses curl | bash as primary method
  - GETTING-STARTED updated to match
  - All examples show `dotclaude create` instead of `cp -r examples/`
- **Updated Docker Test Instructions** - Added Dockerfile.lite option to README

### Fixed
- Install script now works when piped from curl (was broken)
- Profile creation now provides comprehensive starting template

## [0.3.0] - 2025-12-01

### Added
- `active` command - machine-readable profile name for scripting/integration
- `sync_profiles_json()` - auto-generates `~/.claude/profiles.json` for blackdot integration
- Auto-sync of profiles.json on `activate` and `create` commands
- Dockerfile for isolated testing without installation

## [0.2.0] - 2025-11-30

### Added
- Blackwell Systems™ branding with trademark badge
- Centralized Blackwell dark theme from GitHub Pages
- Comprehensive dark mode styling across all documentation
- Interactive click-to-copy for install command on coverpage
- Styled hamburger menu with blue theme color
- Horizontal dividers between BRAND.md sections
- CHANGELOG.md for tracking project changes
- Changelog link in documentation footer

### Changed
- Switched from local CSS to centralized blackwell-docs-theme
- Optimized logo from 1020KB to 56KB (94.5% reduction)
- Streamlined sidebar navigation from 31 to 22 items
- Improved coverpage bulletpoints for clarity and impact
- Enhanced button styling with visual hierarchy and animations
- Reorganized badges for better visual flow
- Italicized "definitive" in tagline with subtle color
- Made arrow bullets bolder (font-weight: 900)
- Two-tone color scheme for bullet text (bold vs regular)
- Updated "Happy coding!" message to remove "Tip:" prefix

### Fixed
- Sidebar "dotclaude" link now goes to coverpage instead of 404
- Install command box padding and text overflow issues
- Button hover text contrast for better readability
- nameLink uses '#/' for proper coverpage loading

## [0.1.0] - 2025-11-30

### Added
- **Profile Management System**
  - One-command profile switching with `dotclaude switch`
  - Layered configuration (base + profile overlays)
  - Auto-detection via `.dotclaude` files in project directories
  - Preview mode with `--dry-run` flag
  - Profile creation with `dotclaude create`
  - Profile editing with `dotclaude edit`

- **Backup & Restore**
  - Automatic versioned backups before profile switches
  - Restore previous configurations with `dotclaude restore`
  - Backup pruning to manage disk space

- **Git Workflow Tools**
  - Branch management for long-lived feature branches
  - Sync command to keep profiles in sync
  - Branch listing with `dotclaude branches`

- **Cross-Platform Support**
  - Linux (Ubuntu, Debian, Arch, Fedora)
  - macOS (Intel and Apple Silicon)
  - WSL2 support

- **Commands**
  - `show` - Display current active profile
  - `list` - List all available profiles
  - `activate` - Activate a profile
  - `switch` - Switch to a different profile
  - `create` - Create a new profile
  - `edit` - Edit profile configuration
  - `diff` - Compare profiles
  - `restore` - Restore from backup
  - `sync` - Sync profiles across machines
  - `branches` - Manage git branches
  - `version` - Show version information
  - `help` - Display help information

- **Testing & Quality**
  - 122 automated tests (security, integration, commands)
  - CI/CD via GitHub Actions
  - Cross-platform testing (Ubuntu + macOS)
  - Installation testing
  - Shellcheck linting
  - Security validation (path traversal, input sanitization)

- **Documentation**
  - Comprehensive README with quick start guide
  - GETTING-STARTED.md for new users
  - USAGE.md with detailed usage patterns
  - COMMANDS.md reference for all commands
  - USER-GUIDE.md for workflows
  - FAQ.md for common questions
  - TROUBLESHOOTING.md for debugging
  - Docsify-powered documentation site

- **Multi-Provider Support**
  - AWS Bedrock integration
  - Claude Max support
  - Provider-specific configurations per profile

- **Brand Assets**
  - Blackwell Systems™ trademark and branding
  - BRAND.md with usage guidelines
  - Custom badge assets
  - Professional documentation theme

### Security
- File locking to prevent concurrent operations
- Input validation for profile names
- Path traversal prevention
- Secure backup handling
- Automatic configuration verification

[Unreleased]: https://github.com/blackwell-systems/dotclaude/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/blackwell-systems/dotclaude/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/blackwell-systems/dotclaude/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/blackwell-systems/dotclaude/releases/tag/v0.1.0
