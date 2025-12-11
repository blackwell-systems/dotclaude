// Package profile handles dotclaude profile management.
package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Profile represents a dotclaude profile.
type Profile struct {
	Name         string
	Path         string
	IsActive     bool
	LastModified time.Time
}

// Manager handles profile operations.
type Manager struct {
	RepoDir     string
	ProfilesDir string
	ClaudeDir   string
	StateFile   string
}

// NewManager creates a new profile manager.
func NewManager(repoDir, claudeDir string) *Manager {
	return &Manager{
		RepoDir:     repoDir,
		ProfilesDir: filepath.Join(repoDir, "profiles"),
		ClaudeDir:   claudeDir,
		StateFile:   filepath.Join(claudeDir, ".current-profile"),
	}
}

// ListProfiles returns all available profiles.
func (m *Manager) ListProfiles() ([]*Profile, error) {
	// Check if profiles directory exists
	if _, err := os.Stat(m.ProfilesDir); os.IsNotExist(err) {
		// Create it if it doesn't exist
		if err := os.MkdirAll(m.ProfilesDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create profiles directory: %w", err)
		}
		return []*Profile{}, nil
	}

	// Read profiles directory
	entries, err := os.ReadDir(m.ProfilesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read profiles directory: %w", err)
	}

	// Get active profile name
	activeProfile := m.GetActiveProfileName()

	profiles := make([]*Profile, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		profilePath := filepath.Join(m.ProfilesDir, entry.Name())

		// Get modification time
		info, err := entry.Info()
		if err != nil {
			continue
		}

		profiles = append(profiles, &Profile{
			Name:         entry.Name(),
			Path:         profilePath,
			IsActive:     entry.Name() == activeProfile,
			LastModified: info.ModTime(),
		})
	}

	// Sort by name
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	return profiles, nil
}

// GetActiveProfileName returns the name of the currently active profile.
func (m *Manager) GetActiveProfileName() string {
	data, err := os.ReadFile(m.StateFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// GetActiveProfile returns the currently active profile, or nil if none.
func (m *Manager) GetActiveProfile() (*Profile, error) {
	name := m.GetActiveProfileName()
	if name == "" {
		return nil, nil
	}

	profilePath := filepath.Join(m.ProfilesDir, name)
	info, err := os.Stat(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return &Profile{
		Name:         name,
		Path:         profilePath,
		IsActive:     true,
		LastModified: info.ModTime(),
	}, nil
}

// ProfileExists checks if a profile with the given name exists.
func (m *Manager) ProfileExists(name string) bool {
	profilePath := filepath.Join(m.ProfilesDir, name)
	_, err := os.Stat(profilePath)
	return err == nil
}

// MaxProfileNameLength is the maximum allowed length for a profile name.
const MaxProfileNameLength = 64

// ValidateProfileName checks if a profile name is valid.
func ValidateProfileName(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	// Check length limit (filesystem compatibility)
	if len(name) > MaxProfileNameLength {
		return fmt.Errorf("profile name too long: %d characters (maximum %d allowed)", len(name), MaxProfileNameLength)
	}

	// Check for invalid characters
	// Allow: letters, numbers, hyphens, underscores
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_') {
			return fmt.Errorf("invalid profile name: %s (only letters, numbers, hyphens, and underscores allowed)", name)
		}
	}

	return nil
}
