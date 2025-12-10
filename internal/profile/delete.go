package profile

import (
	"fmt"
	"os"
)

// Delete removes a profile.
func (m *Manager) Delete(name string) error {
	// Validate profile name
	if err := ValidateProfileName(name); err != nil {
		return err
	}

	// Check if profile exists
	if !m.ProfileExists(name) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	// Check if profile is currently active
	activeProfile := m.GetActiveProfileName()
	if activeProfile == name {
		return fmt.Errorf("cannot delete active profile '%s' (deactivate it first)", name)
	}

	// Delete profile directory
	profilePath := m.ProfilesDir + "/" + name
	if err := os.RemoveAll(profilePath); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}
