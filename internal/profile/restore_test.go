package profile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestListBackups(t *testing.T) {
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

	t.Run("no backups", func(t *testing.T) {
		backups, err := mgr.ListBackups()
		if err != nil {
			t.Fatalf("ListBackups() error = %v", err)
		}
		if len(backups) != 0 {
			t.Errorf("ListBackups() returned %d backups, want 0", len(backups))
		}
	})

	t.Run("with CLAUDE.md backups", func(t *testing.T) {
		// Create some backup files
		timestamps := []string{"20251201-120000", "20251202-130000", "20251203-140000"}
		for _, ts := range timestamps {
			backupPath := filepath.Join(claudeDir, "CLAUDE.md.backup."+ts)
			if err := os.WriteFile(backupPath, []byte("backup content"), 0644); err != nil {
				t.Fatal(err)
			}
			// Sleep to ensure different mod times
			time.Sleep(10 * time.Millisecond)
		}

		backups, err := mgr.ListBackups()
		if err != nil {
			t.Fatalf("ListBackups() error = %v", err)
		}

		if len(backups) != 3 {
			t.Errorf("ListBackups() returned %d backups, want 3", len(backups))
		}

		// All should be CLAUDE.md type
		for _, b := range backups {
			if b.Type != "CLAUDE.md" {
				t.Errorf("Backup type = %q, want %q", b.Type, "CLAUDE.md")
			}
		}
	})

	t.Run("with settings.json backups", func(t *testing.T) {
		// Create settings backup
		backupPath := filepath.Join(claudeDir, "settings.json.backup.20251201-100000")
		if err := os.WriteFile(backupPath, []byte(`{}`), 0644); err != nil {
			t.Fatal(err)
		}

		backups, err := mgr.ListBackups()
		if err != nil {
			t.Fatalf("ListBackups() error = %v", err)
		}

		// Should now have 4 backups total
		if len(backups) != 4 {
			t.Errorf("ListBackups() returned %d backups, want 4", len(backups))
		}

		// Check we have both types
		hasSettings := false
		hasClaude := false
		for _, b := range backups {
			if b.Type == "settings.json" {
				hasSettings = true
			}
			if b.Type == "CLAUDE.md" {
				hasClaude = true
			}
		}

		if !hasSettings {
			t.Error("Should have settings.json backup")
		}
		if !hasClaude {
			t.Error("Should have CLAUDE.md backup")
		}
	})
}

func TestRestore(t *testing.T) {
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

	t.Run("restore CLAUDE.md backup", func(t *testing.T) {
		// Create current CLAUDE.md
		currentPath := filepath.Join(claudeDir, "CLAUDE.md")
		if err := os.WriteFile(currentPath, []byte("current content"), 0644); err != nil {
			t.Fatal(err)
		}

		// Create backup
		backupPath := filepath.Join(claudeDir, "CLAUDE.md.backup.20251201-120000")
		backupContent := "# Profile: old-profile\n\nbackup content"
		if err := os.WriteFile(backupPath, []byte(backupContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Wait a moment
		time.Sleep(10 * time.Millisecond)

		// Restore
		err := mgr.Restore(backupPath)
		if err != nil {
			t.Fatalf("Restore() error = %v", err)
		}

		// Verify restored content
		restored, err := os.ReadFile(currentPath)
		if err != nil {
			t.Fatalf("Failed to read restored file: %v", err)
		}

		if string(restored) != backupContent {
			t.Errorf("Restored content = %q, want %q", string(restored), backupContent)
		}

		// Verify a new backup was created of the "current content"
		backups, _ := filepath.Glob(filepath.Join(claudeDir, "CLAUDE.md.backup.*"))
		if len(backups) < 2 {
			t.Error("Should have created backup of current file before restoring")
		}
	})

	t.Run("restore settings.json backup", func(t *testing.T) {
		// Create current settings.json
		currentPath := filepath.Join(claudeDir, "settings.json")
		if err := os.WriteFile(currentPath, []byte(`{"current": true}`), 0644); err != nil {
			t.Fatal(err)
		}

		// Create backup
		backupPath := filepath.Join(claudeDir, "settings.json.backup.20251201-120000")
		backupContent := `{"old": true}`
		if err := os.WriteFile(backupPath, []byte(backupContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Restore
		err := mgr.Restore(backupPath)
		if err != nil {
			t.Fatalf("Restore() error = %v", err)
		}

		// Verify restored content
		restored, err := os.ReadFile(currentPath)
		if err != nil {
			t.Fatalf("Failed to read restored file: %v", err)
		}

		if string(restored) != backupContent {
			t.Errorf("Restored content = %q, want %q", string(restored), backupContent)
		}
	})

	t.Run("restore non-existent backup", func(t *testing.T) {
		err := mgr.Restore("/non/existent/backup")
		if err == nil {
			t.Error("Restore() should error for non-existent backup")
		}
	})

	t.Run("restore invalid backup filename", func(t *testing.T) {
		// Create file with invalid backup name
		invalidPath := filepath.Join(claudeDir, "invalid.backup.20251201")
		if err := os.WriteFile(invalidPath, []byte("data"), 0644); err != nil {
			t.Fatal(err)
		}

		err := mgr.Restore(invalidPath)
		if err == nil {
			t.Error("Restore() should error for invalid backup filename")
		}
	})
}

func TestUpdateProfileFromCLAUDE(t *testing.T) {
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

	t.Run("update profile from CLAUDE.md", func(t *testing.T) {
		// Create CLAUDE.md with profile header
		claudePath := filepath.Join(claudeDir, "CLAUDE.md")
		content := "# Base content\n\n# Profile: extracted-profile\n\nProfile content"
		if err := os.WriteFile(claudePath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		err := mgr.updateProfileFromCLAUDE(claudePath)
		if err != nil {
			t.Fatalf("updateProfileFromCLAUDE() error = %v", err)
		}

		// Verify state file was updated
		active := mgr.GetActiveProfileName()
		if active != "extracted-profile" {
			t.Errorf("Active profile = %q, want %q", active, "extracted-profile")
		}
	})

	t.Run("no profile header", func(t *testing.T) {
		// Clear state file
		os.Remove(filepath.Join(claudeDir, ".current-profile"))

		// Create CLAUDE.md without profile header
		claudePath := filepath.Join(claudeDir, "CLAUDE-no-profile.md")
		content := "# Base content only\n\nNo profile marker"
		if err := os.WriteFile(claudePath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		err := mgr.updateProfileFromCLAUDE(claudePath)
		if err != nil {
			t.Fatalf("updateProfileFromCLAUDE() should not error: %v", err)
		}

		// State file should not have been created
		active := mgr.GetActiveProfileName()
		if active != "" {
			t.Errorf("Active profile = %q, want empty", active)
		}
	})
}

func TestParseBackup(t *testing.T) {
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

	t.Run("parse CLAUDE.md backup", func(t *testing.T) {
		// Create backup file
		backupPath := filepath.Join(claudeDir, "CLAUDE.md.backup.20251201-153045")
		if err := os.WriteFile(backupPath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		backup, err := mgr.parseBackup(backupPath, "CLAUDE.md")
		if err != nil {
			t.Fatalf("parseBackup() error = %v", err)
		}

		if backup.Path != backupPath {
			t.Errorf("Path = %q, want %q", backup.Path, backupPath)
		}

		if backup.Type != "CLAUDE.md" {
			t.Errorf("Type = %q, want %q", backup.Type, "CLAUDE.md")
		}

		if backup.Timestamp != "20251201-153045" {
			t.Errorf("Timestamp = %q, want %q", backup.Timestamp, "20251201-153045")
		}

		if !strings.Contains(backup.Filename, "CLAUDE.md.backup.") {
			t.Errorf("Filename = %q, should contain backup pattern", backup.Filename)
		}
	})

	t.Run("parse settings.json backup", func(t *testing.T) {
		backupPath := filepath.Join(claudeDir, "settings.json.backup.20251202-100000")
		if err := os.WriteFile(backupPath, []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}

		backup, err := mgr.parseBackup(backupPath, "settings.json")
		if err != nil {
			t.Fatalf("parseBackup() error = %v", err)
		}

		if backup.Type != "settings.json" {
			t.Errorf("Type = %q, want %q", backup.Type, "settings.json")
		}

		if backup.Timestamp != "20251202-100000" {
			t.Errorf("Timestamp = %q, want %q", backup.Timestamp, "20251202-100000")
		}
	})

	t.Run("parse non-existent file", func(t *testing.T) {
		_, err := mgr.parseBackup("/non/existent", "CLAUDE.md")
		if err == nil {
			t.Error("parseBackup() should error for non-existent file")
		}
	})
}
