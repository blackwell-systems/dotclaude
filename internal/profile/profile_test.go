package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateProfileName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		// Valid names
		{"simple lowercase", "myprofile", false},
		{"simple uppercase", "MyProfile", false},
		{"with numbers", "profile123", false},
		{"with hyphens", "my-profile", false},
		{"with underscores", "my_profile", false},
		{"mixed", "My-Profile_123", false},
		{"single char", "a", false},

		// Invalid names
		{"empty", "", true},
		{"with spaces", "my profile", true},
		{"with dots", "my.profile", true},
		{"with slash", "my/profile", true},
		{"with backslash", "my\\profile", true},
		{"with special chars", "my@profile!", true},
		{"with unicode", "プロファイル", true},
		{"starts with space", " profile", true},
		{"ends with space", "profile ", true},
		{"path traversal", "../etc", true},
		{"too long", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", true}, // 65 chars
		{"max length", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false},  // 64 chars
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProfileName(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateProfileName(%q) error = %v, wantError %v", tt.input, err, tt.wantError)
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	repoDir := "/test/repo"
	claudeDir := "/test/claude"

	mgr := NewManager(repoDir, claudeDir)

	if mgr.RepoDir != repoDir {
		t.Errorf("RepoDir = %q, want %q", mgr.RepoDir, repoDir)
	}

	if mgr.ClaudeDir != claudeDir {
		t.Errorf("ClaudeDir = %q, want %q", mgr.ClaudeDir, claudeDir)
	}

	expectedProfilesDir := filepath.Join(repoDir, "profiles")
	if mgr.ProfilesDir != expectedProfilesDir {
		t.Errorf("ProfilesDir = %q, want %q", mgr.ProfilesDir, expectedProfilesDir)
	}

	expectedStateFile := filepath.Join(claudeDir, ".current-profile")
	if mgr.StateFile != expectedStateFile {
		t.Errorf("StateFile = %q, want %q", mgr.StateFile, expectedStateFile)
	}
}

func TestGetActiveProfileName(t *testing.T) {
	// Create temp directory
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

	t.Run("no state file", func(t *testing.T) {
		result := mgr.GetActiveProfileName()
		if result != "" {
			t.Errorf("GetActiveProfileName() = %q, want empty string", result)
		}
	})

	t.Run("with state file", func(t *testing.T) {
		stateFile := filepath.Join(claudeDir, ".current-profile")
		if err := os.WriteFile(stateFile, []byte("test-profile\n"), 0644); err != nil {
			t.Fatal(err)
		}

		result := mgr.GetActiveProfileName()
		if result != "test-profile" {
			t.Errorf("GetActiveProfileName() = %q, want %q", result, "test-profile")
		}
	})

	t.Run("with whitespace", func(t *testing.T) {
		stateFile := filepath.Join(claudeDir, ".current-profile")
		if err := os.WriteFile(stateFile, []byte("  profile-with-spaces  \n"), 0644); err != nil {
			t.Fatal(err)
		}

		result := mgr.GetActiveProfileName()
		if result != "profile-with-spaces" {
			t.Errorf("GetActiveProfileName() = %q, want %q", result, "profile-with-spaces")
		}
	})
}

func TestProfileExists(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	profilesDir := filepath.Join(tmpDir, "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a profile directory
	existingProfile := filepath.Join(profilesDir, "existing-profile")
	if err := os.MkdirAll(existingProfile, 0755); err != nil {
		t.Fatal(err)
	}

	mgr := NewManager(tmpDir, filepath.Join(tmpDir, ".claude"))

	t.Run("existing profile", func(t *testing.T) {
		if !mgr.ProfileExists("existing-profile") {
			t.Error("ProfileExists() = false, want true for existing profile")
		}
	})

	t.Run("non-existing profile", func(t *testing.T) {
		if mgr.ProfileExists("non-existing") {
			t.Error("ProfileExists() = true, want false for non-existing profile")
		}
	})
}

func TestListProfiles(t *testing.T) {
	// Create temp directory
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

	t.Run("empty profiles directory", func(t *testing.T) {
		profiles, err := mgr.ListProfiles()
		if err != nil {
			t.Fatalf("ListProfiles() error = %v", err)
		}
		if len(profiles) != 0 {
			t.Errorf("ListProfiles() returned %d profiles, want 0", len(profiles))
		}
	})

	t.Run("with profiles", func(t *testing.T) {
		// Create some profiles
		profilesDir := filepath.Join(tmpDir, "profiles")
		for _, name := range []string{"alpha", "beta", "gamma"} {
			profileDir := filepath.Join(profilesDir, name)
			if err := os.MkdirAll(profileDir, 0755); err != nil {
				t.Fatal(err)
			}
		}

		profiles, err := mgr.ListProfiles()
		if err != nil {
			t.Fatalf("ListProfiles() error = %v", err)
		}

		if len(profiles) != 3 {
			t.Errorf("ListProfiles() returned %d profiles, want 3", len(profiles))
		}

		// Should be sorted alphabetically
		expectedOrder := []string{"alpha", "beta", "gamma"}
		for i, p := range profiles {
			if p.Name != expectedOrder[i] {
				t.Errorf("profiles[%d].Name = %q, want %q", i, p.Name, expectedOrder[i])
			}
		}
	})

	t.Run("with active profile", func(t *testing.T) {
		// Set active profile
		stateFile := filepath.Join(claudeDir, ".current-profile")
		if err := os.WriteFile(stateFile, []byte("beta"), 0644); err != nil {
			t.Fatal(err)
		}

		profiles, err := mgr.ListProfiles()
		if err != nil {
			t.Fatalf("ListProfiles() error = %v", err)
		}

		for _, p := range profiles {
			expectedActive := p.Name == "beta"
			if p.IsActive != expectedActive {
				t.Errorf("profile %q IsActive = %v, want %v", p.Name, p.IsActive, expectedActive)
			}
		}
	})

	t.Run("filters out files", func(t *testing.T) {
		// Create a file in profiles directory (should be ignored)
		profilesDir := filepath.Join(tmpDir, "profiles")
		filePath := filepath.Join(profilesDir, "not-a-profile.txt")
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		profiles, err := mgr.ListProfiles()
		if err != nil {
			t.Fatalf("ListProfiles() error = %v", err)
		}

		// Should still have only 3 profiles
		if len(profiles) != 3 {
			t.Errorf("ListProfiles() returned %d profiles, want 3", len(profiles))
		}
	})
}

func TestGetActiveProfile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	claudeDir := filepath.Join(tmpDir, ".claude")
	profilesDir := filepath.Join(tmpDir, "profiles")

	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		t.Fatal(err)
	}

	mgr := NewManager(tmpDir, claudeDir)

	t.Run("no active profile", func(t *testing.T) {
		profile, err := mgr.GetActiveProfile()
		if err != nil {
			t.Fatalf("GetActiveProfile() error = %v", err)
		}
		if profile != nil {
			t.Errorf("GetActiveProfile() = %v, want nil", profile)
		}
	})

	t.Run("active profile exists", func(t *testing.T) {
		// Create profile directory
		profileDir := filepath.Join(profilesDir, "my-profile")
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Set as active
		stateFile := filepath.Join(claudeDir, ".current-profile")
		if err := os.WriteFile(stateFile, []byte("my-profile"), 0644); err != nil {
			t.Fatal(err)
		}

		profile, err := mgr.GetActiveProfile()
		if err != nil {
			t.Fatalf("GetActiveProfile() error = %v", err)
		}

		if profile == nil {
			t.Fatal("GetActiveProfile() = nil, want profile")
		}

		if profile.Name != "my-profile" {
			t.Errorf("profile.Name = %q, want %q", profile.Name, "my-profile")
		}

		if !profile.IsActive {
			t.Error("profile.IsActive = false, want true")
		}
	})

	t.Run("active profile does not exist", func(t *testing.T) {
		// Set non-existent profile as active
		stateFile := filepath.Join(claudeDir, ".current-profile")
		if err := os.WriteFile(stateFile, []byte("deleted-profile"), 0644); err != nil {
			t.Fatal(err)
		}

		profile, err := mgr.GetActiveProfile()
		if err != nil {
			t.Fatalf("GetActiveProfile() error = %v", err)
		}
		if profile != nil {
			t.Errorf("GetActiveProfile() = %v, want nil", profile)
		}
	})
}

func TestListProfiles_EdgeCases(t *testing.T) {
	t.Run("handles profiles with Info errors gracefully", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		profilesDir := filepath.Join(tmpDir, "profiles")
		if err := os.MkdirAll(profilesDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create a valid profile
		goodProfile := filepath.Join(profilesDir, "good-profile")
		if err := os.MkdirAll(goodProfile, 0755); err != nil {
			t.Fatal(err)
		}

		mgr := NewManager(tmpDir, filepath.Join(tmpDir, ".claude"))
		profiles, err := mgr.ListProfiles()

		// Should succeed
		if err != nil {
			t.Fatalf("ListProfiles() error = %v", err)
		}

		// Should have the good profile
		if len(profiles) != 1 {
			t.Errorf("ListProfiles() returned %d profiles, want 1", len(profiles))
		}
	})

	t.Run("creates profiles directory if missing", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		mgr := NewManager(tmpDir, filepath.Join(tmpDir, ".claude"))

		// profiles directory doesn't exist yet
		profilesDir := filepath.Join(tmpDir, "profiles")
		if _, err := os.Stat(profilesDir); !os.IsNotExist(err) {
			t.Fatal("profiles directory should not exist before ListProfiles")
		}

		profiles, err := mgr.ListProfiles()
		if err != nil {
			t.Fatalf("ListProfiles() error = %v", err)
		}

		if len(profiles) != 0 {
			t.Errorf("ListProfiles() should return empty slice for new profiles dir")
		}

		// Verify directory was created
		if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
			t.Error("ListProfiles() should create profiles directory")
		}
	})
}
