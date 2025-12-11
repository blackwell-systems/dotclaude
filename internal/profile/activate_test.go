package profile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestActivate(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	claudeDir := filepath.Join(tmpDir, ".claude")
	mgr := NewManager(tmpDir, claudeDir)

	// Create a test profile
	profilesDir := filepath.Join(tmpDir, "profiles")
	testProfile := filepath.Join(profilesDir, "test-profile")
	if err := os.MkdirAll(testProfile, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testProfile, "CLAUDE.md"), []byte("# Test Profile\n"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("activate profile", func(t *testing.T) {
		err := mgr.Activate("test-profile")
		if err != nil {
			t.Fatalf("Activate() error = %v", err)
		}

		// Verify state file was updated
		active := mgr.GetActiveProfileName()
		if active != "test-profile" {
			t.Errorf("Active profile = %q, want %q", active, "test-profile")
		}

		// Verify CLAUDE.md was merged
		claudeMd := filepath.Join(claudeDir, "CLAUDE.md")
		content, err := os.ReadFile(claudeMd)
		if err != nil {
			t.Fatalf("Failed to read CLAUDE.md: %v", err)
		}

		// Should contain base content
		if !strings.Contains(string(content), "# Base Config") {
			t.Error("Merged CLAUDE.md should contain base content")
		}

		// Should contain profile content
		if !strings.Contains(string(content), "# Test Profile") {
			t.Error("Merged CLAUDE.md should contain profile content")
		}

		// Should contain separator with profile name
		if !strings.Contains(string(content), "# Profile: test-profile") {
			t.Error("Merged CLAUDE.md should contain profile separator")
		}
	})

	t.Run("activate non-existent profile", func(t *testing.T) {
		err := mgr.Activate("non-existent")
		if err == nil {
			t.Error("Activate() should error for non-existent profile")
		}
	})

	t.Run("activate with invalid name", func(t *testing.T) {
		err := mgr.Activate("invalid/name")
		if err == nil {
			t.Error("Activate() should error for invalid name")
		}
	})

	t.Run("switch profiles creates backup", func(t *testing.T) {
		// Create another profile
		profile2 := filepath.Join(profilesDir, "profile-2")
		if err := os.MkdirAll(profile2, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(profile2, "CLAUDE.md"), []byte("# Profile 2\n"), 0644); err != nil {
			t.Fatal(err)
		}

		// Ensure first profile is active
		mgr.Activate("test-profile")

		// Wait a moment so backup timestamps differ
		time.Sleep(10 * time.Millisecond)

		// Switch to second profile
		err := mgr.Activate("profile-2")
		if err != nil {
			t.Fatalf("Activate() error = %v", err)
		}

		// Verify backup was created
		backups, _ := filepath.Glob(filepath.Join(claudeDir, "CLAUDE.md.backup.*"))
		if len(backups) == 0 {
			t.Error("Backup should have been created when switching profiles")
		}
	})
}

func TestBackupFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	mgr := NewManager(tmpDir, claudeDir)

	t.Run("backup non-existent file", func(t *testing.T) {
		// Should not error for non-existent file
		err := mgr.backupFile("non-existent.txt")
		if err != nil {
			t.Errorf("backupFile() should not error for non-existent file: %v", err)
		}
	})

	t.Run("backup existing file", func(t *testing.T) {
		// Create file
		filePath := filepath.Join(claudeDir, "CLAUDE.md")
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}

		err := mgr.backupFile("CLAUDE.md")
		if err != nil {
			t.Fatalf("backupFile() error = %v", err)
		}

		// Verify backup was created
		backups, err := filepath.Glob(filepath.Join(claudeDir, "CLAUDE.md.backup.*"))
		if err != nil {
			t.Fatal(err)
		}

		if len(backups) == 0 {
			t.Error("Backup file should have been created")
		}

		// Verify backup content
		backupContent, err := os.ReadFile(backups[0])
		if err != nil {
			t.Fatal(err)
		}

		if string(backupContent) != "test content" {
			t.Errorf("Backup content = %q, want %q", string(backupContent), "test content")
		}
	})
}

func TestCleanupBackups(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	mgr := NewManager(tmpDir, claudeDir)

	// Create multiple backup files
	for i := 0; i < 7; i++ {
		timestamp := time.Now().Add(time.Duration(i) * time.Second).Format("20060102-150405")
		backupPath := filepath.Join(claudeDir, "CLAUDE.md.backup."+timestamp)
		if err := os.WriteFile(backupPath, []byte("backup"), 0644); err != nil {
			t.Fatal(err)
		}
		time.Sleep(10 * time.Millisecond) // Ensure different mod times
	}

	// Verify we have 7 backups
	backups, _ := filepath.Glob(filepath.Join(claudeDir, "CLAUDE.md.backup.*"))
	if len(backups) != 7 {
		t.Fatalf("Should have 7 backups before cleanup, got %d", len(backups))
	}

	// Cleanup, keeping only 3
	err = mgr.cleanupBackups("CLAUDE.md", 3)
	if err != nil {
		t.Fatalf("cleanupBackups() error = %v", err)
	}

	// Verify only 3 remain
	backups, _ = filepath.Glob(filepath.Join(claudeDir, "CLAUDE.md.backup.*"))
	if len(backups) != 3 {
		t.Errorf("Should have 3 backups after cleanup, got %d", len(backups))
	}
}

func TestMergeCLAUDEmd(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	claudeDir := filepath.Join(tmpDir, ".claude")
	mgr := NewManager(tmpDir, claudeDir)

	// Create a profile
	profilesDir := filepath.Join(tmpDir, "profiles")
	testProfile := filepath.Join(profilesDir, "merge-test")
	if err := os.MkdirAll(testProfile, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testProfile, "CLAUDE.md"), []byte("# Profile Content\n"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("merge CLAUDE.md", func(t *testing.T) {
		err := mgr.mergeCLAUDEmd("merge-test")
		if err != nil {
			t.Fatalf("mergeCLAUDEmd() error = %v", err)
		}

		// Read merged file
		content, err := os.ReadFile(filepath.Join(claudeDir, "CLAUDE.md"))
		if err != nil {
			t.Fatalf("Failed to read merged file: %v", err)
		}

		// Verify structure
		contentStr := string(content)

		if !strings.HasPrefix(contentStr, "# Base Config") {
			t.Error("Merged file should start with base content")
		}

		if !strings.Contains(contentStr, "# Profile: merge-test") {
			t.Error("Merged file should contain profile separator")
		}

		if !strings.Contains(contentStr, "# Profile Content") {
			t.Error("Merged file should contain profile content")
		}
	})

	t.Run("merge non-existent profile", func(t *testing.T) {
		err := mgr.mergeCLAUDEmd("non-existent")
		if err == nil {
			t.Error("mergeCLAUDEmd() should error for non-existent profile")
		}
	})
}

func TestApplySettings(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	claudeDir := filepath.Join(tmpDir, ".claude")
	mgr := NewManager(tmpDir, claudeDir)

	// Create a profile
	profilesDir := filepath.Join(tmpDir, "profiles")
	testProfile := filepath.Join(profilesDir, "settings-test")
	if err := os.MkdirAll(testProfile, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testProfile, "CLAUDE.md"), []byte("# Test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("apply profile settings", func(t *testing.T) {
		// Create profile-specific settings
		profileSettings := `{"profile": "custom"}`
		if err := os.WriteFile(filepath.Join(testProfile, "settings.json"), []byte(profileSettings), 0644); err != nil {
			t.Fatal(err)
		}

		err := mgr.applySettings("settings-test")
		if err != nil {
			t.Fatalf("applySettings() error = %v", err)
		}

		// Verify settings were applied
		content, err := os.ReadFile(filepath.Join(claudeDir, "settings.json"))
		if err != nil {
			t.Fatalf("Failed to read settings: %v", err)
		}

		if string(content) != profileSettings {
			t.Errorf("Settings = %q, want %q", string(content), profileSettings)
		}
	})

	t.Run("fallback to base settings", func(t *testing.T) {
		// Create profile without settings.json
		noSettingsProfile := filepath.Join(profilesDir, "no-settings")
		if err := os.MkdirAll(noSettingsProfile, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(noSettingsProfile, "CLAUDE.md"), []byte("# Test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		err := mgr.applySettings("no-settings")
		if err != nil {
			t.Fatalf("applySettings() error = %v", err)
		}

		// Verify base settings were applied
		content, err := os.ReadFile(filepath.Join(claudeDir, "settings.json"))
		if err != nil {
			t.Fatalf("Failed to read settings: %v", err)
		}

		if !strings.Contains(string(content), `"key": "value"`) {
			t.Error("Settings should contain base settings content")
		}
	})
}

func TestActivate_ErrorPaths(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	claudeDir := filepath.Join(tmpDir, ".claude")
	mgr := NewManager(tmpDir, claudeDir)

	t.Run("activate with empty name", func(t *testing.T) {
		err := mgr.Activate("")
		if err == nil {
			t.Error("Activate() should error for empty name")
		}
	})

	t.Run("activate with invalid characters", func(t *testing.T) {
		err := mgr.Activate("profile with spaces")
		if err == nil {
			t.Error("Activate() should error for invalid name")
		}
	})

	t.Run("activate creates claude directory", func(t *testing.T) {
		// Create a new temp dir without .claude
		newTmp, err := os.MkdirTemp("", "dotclaude-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(newTmp)

		// Copy repo structure
		baseDir := filepath.Join(newTmp, "base")
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(baseDir, "CLAUDE.md"), []byte("# Base\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(baseDir, "settings.json"), []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}

		profilesDir := filepath.Join(newTmp, "profiles", "test-profile")
		if err := os.MkdirAll(profilesDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(profilesDir, "CLAUDE.md"), []byte("# Test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		newClaudeDir := filepath.Join(newTmp, ".claude-new")
		mgr2 := NewManager(newTmp, newClaudeDir)

		err = mgr2.Activate("test-profile")
		if err != nil {
			t.Fatalf("Activate() should create claude dir: %v", err)
		}

		if _, err := os.Stat(newClaudeDir); os.IsNotExist(err) {
			t.Error("Activate() should have created claude directory")
		}
	})

	t.Run("activate same profile twice", func(t *testing.T) {
		// Create a profile
		profileDir := filepath.Join(tmpDir, "profiles", "same-twice")
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(profileDir, "CLAUDE.md"), []byte("# Test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		// Activate first time
		err := mgr.Activate("same-twice")
		if err != nil {
			t.Fatalf("First Activate() error: %v", err)
		}

		// Activate second time (same profile)
		err = mgr.Activate("same-twice")
		if err != nil {
			t.Fatalf("Second Activate() should succeed: %v", err)
		}
	})
}

func TestCleanupBackups_EdgeCases(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	mgr := NewManager(tmpDir, claudeDir)

	t.Run("cleanup with no backups", func(t *testing.T) {
		err := mgr.cleanupBackups("CLAUDE.md", 5)
		if err != nil {
			t.Errorf("cleanupBackups() should not error for no backups: %v", err)
		}
	})

	t.Run("cleanup with fewer than limit", func(t *testing.T) {
		// Create 2 backups
		for i := 0; i < 2; i++ {
			timestamp := time.Now().Add(time.Duration(i) * time.Second).Format("20060102-150405")
			backupPath := filepath.Join(claudeDir, "test.backup."+timestamp)
			if err := os.WriteFile(backupPath, []byte("backup"), 0644); err != nil {
				t.Fatal(err)
			}
		}

		err := mgr.cleanupBackups("test", 5)
		if err != nil {
			t.Errorf("cleanupBackups() error = %v", err)
		}

		// Should still have 2 backups
		backups, _ := filepath.Glob(filepath.Join(claudeDir, "test.backup.*"))
		if len(backups) != 2 {
			t.Errorf("Should have 2 backups, got %d", len(backups))
		}
	})

	t.Run("cleanup with exactly limit", func(t *testing.T) {
		// Clear old backups
		oldBackups, _ := filepath.Glob(filepath.Join(claudeDir, "exact.backup.*"))
		for _, b := range oldBackups {
			os.Remove(b)
		}

		// Create exactly 3 backups
		for i := 0; i < 3; i++ {
			timestamp := time.Now().Add(time.Duration(i) * time.Second).Format("20060102-150405")
			backupPath := filepath.Join(claudeDir, "exact.backup."+timestamp)
			if err := os.WriteFile(backupPath, []byte("backup"), 0644); err != nil {
				t.Fatal(err)
			}
			time.Sleep(10 * time.Millisecond)
		}

		err := mgr.cleanupBackups("exact", 3)
		if err != nil {
			t.Errorf("cleanupBackups() error = %v", err)
		}

		// Should still have 3 backups
		backups, _ := filepath.Glob(filepath.Join(claudeDir, "exact.backup.*"))
		if len(backups) != 3 {
			t.Errorf("Should have 3 backups, got %d", len(backups))
		}
	})
}

func TestMergeCLAUDEmd_ErrorPaths(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	claudeDir := filepath.Join(tmpDir, ".claude")
	mgr := NewManager(tmpDir, claudeDir)

	t.Run("missing base CLAUDE.md", func(t *testing.T) {
		// Remove base CLAUDE.md
		os.Remove(filepath.Join(tmpDir, "base", "CLAUDE.md"))

		// Create a profile
		profileDir := filepath.Join(tmpDir, "profiles", "test-merge")
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(profileDir, "CLAUDE.md"), []byte("# Test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		err := mgr.mergeCLAUDEmd("test-merge")
		if err == nil {
			t.Error("mergeCLAUDEmd() should error when base CLAUDE.md is missing")
		}
	})
}

func TestApplySettings_ErrorPaths(t *testing.T) {
	t.Run("missing both profile and base settings", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create minimal structure without settings
		baseDir := filepath.Join(tmpDir, "base")
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			t.Fatal(err)
		}

		profileDir := filepath.Join(tmpDir, "profiles", "no-settings-anywhere")
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}

		claudeDir := filepath.Join(tmpDir, ".claude")
		if err := os.MkdirAll(claudeDir, 0755); err != nil {
			t.Fatal(err)
		}

		mgr := NewManager(tmpDir, claudeDir)
		err = mgr.applySettings("no-settings-anywhere")
		if err == nil {
			t.Error("applySettings() should error when no settings exist")
		}
	})
}
