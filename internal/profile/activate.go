package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Activate activates a profile by merging base + profile configuration.
func (m *Manager) Activate(name string) error {
	// Validate profile name
	if err := ValidateProfileName(name); err != nil {
		return err
	}

	// Check if profile exists
	if !m.ProfileExists(name) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	// Get current active profile
	currentProfile := m.GetActiveProfileName()

	// Ensure Claude directory exists
	if err := os.MkdirAll(m.ClaudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create Claude directory: %w", err)
	}

	// Backup existing files if switching profiles
	if currentProfile != name {
		if err := m.backupFile("CLAUDE.md"); err != nil {
			return fmt.Errorf("failed to backup CLAUDE.md: %w", err)
		}
		if err := m.backupFile("settings.json"); err != nil {
			return fmt.Errorf("failed to backup settings.json: %w", err)
		}
	}

	// Merge base + profile CLAUDE.md
	if err := m.mergeCLAUDEmd(name); err != nil {
		return fmt.Errorf("failed to merge CLAUDE.md: %w", err)
	}

	// Apply settings.json
	if err := m.applySettings(name); err != nil {
		return fmt.Errorf("failed to apply settings: %w", err)
	}

	// Mark as active
	stateFile := filepath.Join(m.ClaudeDir, ".current-profile")
	if err := os.WriteFile(stateFile, []byte(name), 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// backupFile creates a timestamped backup of a file in the Claude directory.
// Keeps only the 5 most recent backups.
func (m *Manager) backupFile(filename string) error {
	sourcePath := filepath.Join(m.ClaudeDir, filename)

	// Skip if file doesn't exist
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return nil
	}

	// Create backup with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(m.ClaudeDir, fmt.Sprintf("%s.backup.%s", filename, timestamp))

	// Copy file
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	// Cleanup old backups (keep only 5 most recent)
	if err := m.cleanupBackups(filename, 5); err != nil {
		// Log but don't fail on cleanup errors
		fmt.Fprintf(os.Stderr, "warning: failed to cleanup old backups: %v\n", err)
	}

	return nil
}

// cleanupBackups removes old backup files, keeping only the N most recent.
func (m *Manager) cleanupBackups(filename string, keep int) error {
	pattern := filepath.Join(m.ClaudeDir, fmt.Sprintf("%s.backup.*", filename))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	// Sort by modification time (newest first)
	sort.Slice(matches, func(i, j int) bool {
		infoI, errI := os.Stat(matches[i])
		infoJ, errJ := os.Stat(matches[j])
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Remove old backups beyond the keep limit
	for i := keep; i < len(matches); i++ {
		if err := os.Remove(matches[i]); err != nil {
			return err
		}
	}

	return nil
}

// mergeCLAUDEmd merges base/CLAUDE.md + profile/CLAUDE.md into Claude directory.
func (m *Manager) mergeCLAUDEmd(profileName string) error {
	basePath := filepath.Join(m.RepoDir, "base", "CLAUDE.md")
	profilePath := filepath.Join(m.ProfilesDir, profileName, "CLAUDE.md")
	outputPath := filepath.Join(m.ClaudeDir, "CLAUDE.md")

	// Read base CLAUDE.md
	baseContent, err := os.ReadFile(basePath)
	if err != nil {
		return fmt.Errorf("failed to read base CLAUDE.md: %w", err)
	}

	// Read profile CLAUDE.md
	profileContent, err := os.ReadFile(profilePath)
	if err != nil {
		return fmt.Errorf("failed to read profile CLAUDE.md: %w", err)
	}

	// Merge with separator
	separator := fmt.Sprintf("\n\n# =========================================\n# Profile: %s\n# =========================================\n\n", profileName)
	merged := string(baseContent) + separator + string(profileContent)

	// Write merged content
	if err := os.WriteFile(outputPath, []byte(merged), 0644); err != nil {
		return fmt.Errorf("failed to write merged CLAUDE.md: %w", err)
	}

	return nil
}

// applySettings copies settings.json from profile (or base if profile doesn't have one).
func (m *Manager) applySettings(profileName string) error {
	profileSettingsPath := filepath.Join(m.ProfilesDir, profileName, "settings.json")
	baseSettingsPath := filepath.Join(m.RepoDir, "base", "settings.json")
	outputPath := filepath.Join(m.ClaudeDir, "settings.json")

	// Try profile settings first
	sourcePath := profileSettingsPath
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		// Fall back to base settings
		sourcePath = baseSettingsPath
	}

	// Read settings
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	// Write settings
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}
