package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDelete(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	claudeDir := filepath.Join(tmpDir, ".claude")
	mgr := NewManager(tmpDir, claudeDir)

	// Create a profile first
	profilesDir := filepath.Join(tmpDir, "profiles")
	testProfile := filepath.Join(profilesDir, "test-profile")
	if err := os.MkdirAll(testProfile, 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("delete existing profile", func(t *testing.T) {
		// Create a profile to delete
		deleteMe := filepath.Join(profilesDir, "delete-me")
		if err := os.MkdirAll(deleteMe, 0755); err != nil {
			t.Fatal(err)
		}

		err := mgr.Delete("delete-me")
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// Verify profile was deleted
		if _, err := os.Stat(deleteMe); !os.IsNotExist(err) {
			t.Error("Profile directory should have been deleted")
		}
	})

	t.Run("delete non-existent profile", func(t *testing.T) {
		err := mgr.Delete("non-existent")
		if err == nil {
			t.Error("Delete() should error for non-existent profile")
		}
	})

	t.Run("delete active profile", func(t *testing.T) {
		// Create a profile
		activeProfile := filepath.Join(profilesDir, "active-profile")
		if err := os.MkdirAll(activeProfile, 0755); err != nil {
			t.Fatal(err)
		}

		// Set as active
		stateFile := filepath.Join(claudeDir, ".current-profile")
		if err := os.WriteFile(stateFile, []byte("active-profile"), 0644); err != nil {
			t.Fatal(err)
		}

		err := mgr.Delete("active-profile")
		if err == nil {
			t.Error("Delete() should error for active profile")
		}

		// Profile should still exist
		if _, err := os.Stat(activeProfile); os.IsNotExist(err) {
			t.Error("Active profile should not have been deleted")
		}
	})

	t.Run("delete with invalid name", func(t *testing.T) {
		err := mgr.Delete("invalid/name")
		if err == nil {
			t.Error("Delete() should error for invalid name")
		}
	})

	t.Run("delete with empty name", func(t *testing.T) {
		err := mgr.Delete("")
		if err == nil {
			t.Error("Delete() should error for empty name")
		}
	})

	t.Run("delete profile with files", func(t *testing.T) {
		// Create a profile with nested content
		profileWithFiles := filepath.Join(profilesDir, "profile-with-files")
		if err := os.MkdirAll(profileWithFiles, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(profileWithFiles, "CLAUDE.md"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		subDir := filepath.Join(profileWithFiles, "subdir")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0644); err != nil {
			t.Fatal(err)
		}

		err := mgr.Delete("profile-with-files")
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// Verify everything was deleted
		if _, err := os.Stat(profileWithFiles); !os.IsNotExist(err) {
			t.Error("Profile directory should have been deleted")
		}
	})
}
