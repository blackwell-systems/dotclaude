package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Backup represents a backup file.
type Backup struct {
	Path      string
	Filename  string
	Timestamp string
	Size      int64
	Type      string // "CLAUDE.md" or "settings.json"
}

// ListBackups returns all available backup files sorted by modification time (newest first).
func (m *Manager) ListBackups() ([]*Backup, error) {
	// Find CLAUDE.md backups
	claudePattern := filepath.Join(m.ClaudeDir, "CLAUDE.md.backup.*")
	claudeMatches, err := filepath.Glob(claudePattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find CLAUDE.md backups: %w", err)
	}

	// Find settings.json backups
	settingsPattern := filepath.Join(m.ClaudeDir, "settings.json.backup.*")
	settingsMatches, err := filepath.Glob(settingsPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find settings.json backups: %w", err)
	}

	var backups []*Backup

	// Process CLAUDE.md backups
	for _, path := range claudeMatches {
		backup, err := m.parseBackup(path, "CLAUDE.md")
		if err != nil {
			continue // Skip invalid backups
		}
		backups = append(backups, backup)
	}

	// Process settings.json backups
	for _, path := range settingsMatches {
		backup, err := m.parseBackup(path, "settings.json")
		if err != nil {
			continue // Skip invalid backups
		}
		backups = append(backups, backup)
	}

	// Sort by modification time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		infoI, errI := os.Stat(backups[i].Path)
		infoJ, errJ := os.Stat(backups[j].Path)
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().After(infoJ.ModTime())
	})

	return backups, nil
}

// parseBackup extracts backup information from a backup file path.
func (m *Manager) parseBackup(path, backupType string) (*Backup, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	filename := filepath.Base(path)

	// Extract timestamp from filename
	// Format: CLAUDE.md.backup.20251210-155544
	var timestamp string
	if backupType == "CLAUDE.md" {
		timestamp = strings.TrimPrefix(filename, "CLAUDE.md.backup.")
	} else {
		timestamp = strings.TrimPrefix(filename, "settings.json.backup.")
	}

	return &Backup{
		Path:      path,
		Filename:  filename,
		Timestamp: timestamp,
		Size:      info.Size(),
		Type:      backupType,
	}, nil
}

// Restore restores a backup file to its original location.
// It creates a backup of the current file before restoring.
func (m *Manager) Restore(backupPath string) error {
	// Validate backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupPath)
	}

	// Determine target file based on backup filename
	filename := filepath.Base(backupPath)
	var targetPath string

	if strings.HasPrefix(filename, "CLAUDE.md.backup.") {
		targetPath = filepath.Join(m.ClaudeDir, "CLAUDE.md")
	} else if strings.HasPrefix(filename, "settings.json.backup.") {
		targetPath = filepath.Join(m.ClaudeDir, "settings.json")
	} else {
		return fmt.Errorf("invalid backup filename: %s", filename)
	}

	// Create backup of current file before restoring
	if _, err := os.Stat(targetPath); err == nil {
		timestamp := time.Now().Format("20060102-150405")
		currentBackup := targetPath + ".backup." + timestamp

		data, err := os.ReadFile(targetPath)
		if err != nil {
			return fmt.Errorf("failed to read current file: %w", err)
		}

		if err := os.WriteFile(currentBackup, data, 0644); err != nil {
			return fmt.Errorf("failed to backup current file: %w", err)
		}
	}

	// Restore the backup
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// If restoring CLAUDE.md, try to update .current-profile marker
	if strings.HasPrefix(filename, "CLAUDE.md.backup.") {
		if err := m.updateProfileFromCLAUDE(targetPath); err != nil {
			// Log but don't fail on this error
			fmt.Fprintf(os.Stderr, "warning: could not update active profile: %v\n", err)
		}
	}

	return nil
}

// updateProfileFromCLAUDE attempts to extract the profile name from CLAUDE.md content
// and update the .current-profile marker.
func (m *Manager) updateProfileFromCLAUDE(claudePath string) error {
	data, err := os.ReadFile(claudePath)
	if err != nil {
		return err
	}

	// Look for "# Profile: <name>" pattern
	re := regexp.MustCompile(`(?m)^# Profile: ([^\s]+)`)
	matches := re.FindStringSubmatch(string(data))

	if len(matches) > 1 {
		profileName := matches[1]
		stateFile := filepath.Join(m.ClaudeDir, ".current-profile")
		return os.WriteFile(stateFile, []byte(profileName), 0644)
	}

	return nil
}
