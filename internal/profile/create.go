package profile

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Create creates a new profile from the template.
func (m *Manager) Create(name string) error {
	// Validate profile name
	if err := ValidateProfileName(name); err != nil {
		return err
	}

	// Check if profile already exists
	if m.ProfileExists(name) {
		return fmt.Errorf("profile '%s' already exists", name)
	}

	// Ensure profiles directory exists
	if err := os.MkdirAll(m.ProfilesDir, 0755); err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	// Template location
	templateDir := filepath.Join(m.RepoDir, "examples", "sample-profile")
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("template not found at %s (is dotclaude repository corrupt?)", templateDir)
	}

	// Destination
	profileDir := filepath.Join(m.ProfilesDir, name)

	// Copy template to profile directory
	if err := copyDir(templateDir, profileDir); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// Initialize git repository in profile
	if err := initGitRepo(profileDir, name); err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	return nil
}

// copyDir recursively copies a directory.
func copyDir(src, dst string) error {
	// Get properties of source
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursive copy for directories
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file.
func copyFile(src, dst string) error {
	// Open source
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Set permissions
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return err
	}

	return nil
}

// initGitRepo initializes a git repository in the profile directory.
func initGitRepo(profileDir, profileName string) error {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		// Git not available, skip initialization (not fatal)
		return nil
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = profileDir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Add all files
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = profileDir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Initial commit
	commitMsg := fmt.Sprintf("Initial commit for profile: %s", profileName)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = profileDir
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
