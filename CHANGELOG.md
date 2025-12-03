# Changelog

All notable changes to dotclaude will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
- `sync_profiles_json()` - auto-generates `~/.claude/profiles.json` for dotfiles integration
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
