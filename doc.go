// Package dotclaude provides profile management for Claude Code.
//
// dotclaude allows you to manage multiple Claude Code configurations (profiles)
// for different contexts like open-source work, employer projects, or client work.
// Each profile contains a CLAUDE.md file with context-specific instructions and
// optional settings overrides.
//
// # Features
//
//   - Profile Management: Create, activate, delete, and list profiles
//   - Configuration Merging: Combine base config with profile-specific settings
//   - Backup and Restore: Automatic backups when switching profiles
//   - Cross-Platform: Works on Linux, macOS, and Windows
//
// # Installation
//
// Download the latest release for your platform from:
// https://github.com/blackwell-systems/dotclaude/releases
//
// Or install with Go:
//
//	go install github.com/blackwell-systems/dotclaude/cmd/dotclaude@latest
//
// # Usage
//
// Create a new profile:
//
//	dotclaude create my-profile
//
// Activate a profile:
//
//	dotclaude activate my-profile
//
// List available profiles:
//
//	dotclaude list
//
// See the full documentation at https://github.com/blackwell-systems/dotclaude
package dotclaude
