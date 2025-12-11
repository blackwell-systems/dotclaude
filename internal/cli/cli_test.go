package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// setupTestEnv creates a temp directory structure for testing CLI commands
func setupTestEnv(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "dotclaude-cli-test-*")
	if err != nil {
		t.Fatal(err)
	}

	// Create base directory
	baseDir := filepath.Join(tmpDir, "base")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(baseDir, "CLAUDE.md"), []byte("# Base Config\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(baseDir, "settings.json"), []byte(`{"key": "value"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create sample-profile template
	templateDir := filepath.Join(tmpDir, "examples", "sample-profile")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "CLAUDE.md"), []byte("# Sample Profile\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .claude directory
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create profiles directory
	profilesDir := filepath.Join(tmpDir, "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Save original values
	origRepoDir := RepoDir
	origClaudeDir := ClaudeDir
	origProfilesDir := ProfilesDir

	// Set test values
	RepoDir = tmpDir
	ClaudeDir = claudeDir
	ProfilesDir = profilesDir

	cleanup := func() {
		os.RemoveAll(tmpDir)
		// Restore original values
		RepoDir = origRepoDir
		ClaudeDir = origClaudeDir
		ProfilesDir = origProfilesDir
	}

	return tmpDir, cleanup
}

// executeCommand executes a command and returns any error
// Note: Commands print to os.Stdout directly, so output capture is limited
func executeCommand(cmd *cobra.Command, args ...string) error {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestVersionCmd(t *testing.T) {
	cmd := newVersionCmd()
	err := executeCommand(cmd)

	if err != nil {
		t.Fatalf("version command error: %v", err)
	}
}

func TestListCmd(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("empty profiles", func(t *testing.T) {
		cmd := newListCmd()
		err := executeCommand(cmd)

		if err != nil {
			t.Fatalf("list command error: %v", err)
		}
	})

	t.Run("with profiles", func(t *testing.T) {
		// Create some profiles
		for _, name := range []string{"alpha", "beta"} {
			profileDir := filepath.Join(ProfilesDir, name)
			if err := os.MkdirAll(profileDir, 0755); err != nil {
				t.Fatal(err)
			}
		}

		cmd := newListCmd()
		err := executeCommand(cmd)

		if err != nil {
			t.Fatalf("list command error: %v", err)
		}
	})
}

func TestShowCmd(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("no active profile", func(t *testing.T) {
		cmd := newShowCmd()
		err := executeCommand(cmd)

		if err != nil {
			t.Fatalf("show command error: %v", err)
		}
	})

	t.Run("with active profile", func(t *testing.T) {
		// Create and activate a profile
		profileDir := filepath.Join(ProfilesDir, "test-profile")
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(profileDir, "CLAUDE.md"), []byte("# Test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		// Set as active
		stateFile := filepath.Join(ClaudeDir, ".current-profile")
		if err := os.WriteFile(stateFile, []byte("test-profile"), 0644); err != nil {
			t.Fatal(err)
		}

		cmd := newShowCmd()
		err := executeCommand(cmd)

		if err != nil {
			t.Fatalf("show command error: %v", err)
		}
	})
}

func TestCreateCmd(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("create valid profile", func(t *testing.T) {
		cmd := newCreateCmd()
		err := executeCommand(cmd, "new-profile")

		if err != nil {
			t.Fatalf("create command error: %v", err)
		}

		// Verify profile directory exists
		profileDir := filepath.Join(ProfilesDir, "new-profile")
		if _, err := os.Stat(profileDir); os.IsNotExist(err) {
			t.Error("profile directory should have been created")
		}
	})

	t.Run("create without name", func(t *testing.T) {
		cmd := newCreateCmd()
		err := executeCommand(cmd)

		if err == nil {
			t.Error("create without name should error")
		}
	})

	t.Run("create duplicate profile", func(t *testing.T) {
		// Create a profile first
		profileDir := filepath.Join(ProfilesDir, "existing")
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}

		cmd := newCreateCmd()
		err := executeCommand(cmd, "existing")

		if err == nil {
			t.Error("creating duplicate profile should error")
		}
	})

	t.Run("create with invalid name", func(t *testing.T) {
		cmd := newCreateCmd()
		err := executeCommand(cmd, "invalid/name")

		if err == nil {
			t.Error("creating profile with invalid name should error")
		}
	})
}

func TestDeleteCmd(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("delete non-existent profile", func(t *testing.T) {
		cmd := newDeleteCmd()
		err := executeCommand(cmd, "--force", "non-existent")

		if err == nil {
			t.Error("deleting non-existent profile should error")
		}
	})

	t.Run("delete existing profile", func(t *testing.T) {
		// Create a profile
		profileDir := filepath.Join(ProfilesDir, "to-delete")
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}

		cmd := newDeleteCmd()
		err := executeCommand(cmd, "--force", "to-delete")

		if err != nil {
			t.Fatalf("delete command error: %v", err)
		}

		// Verify profile was deleted
		if _, err := os.Stat(profileDir); !os.IsNotExist(err) {
			t.Error("profile directory should have been deleted")
		}
	})

	t.Run("delete active profile", func(t *testing.T) {
		// Create and activate a profile
		profileDir := filepath.Join(ProfilesDir, "active-to-delete")
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Set as active
		stateFile := filepath.Join(ClaudeDir, ".current-profile")
		if err := os.WriteFile(stateFile, []byte("active-to-delete"), 0644); err != nil {
			t.Fatal(err)
		}

		cmd := newDeleteCmd()
		err := executeCommand(cmd, "--force", "active-to-delete")

		if err == nil {
			t.Error("deleting active profile should error")
		}
	})
}

func TestActivateCmd(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("activate existing profile", func(t *testing.T) {
		// Create a profile
		profileDir := filepath.Join(ProfilesDir, "to-activate")
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(profileDir, "CLAUDE.md"), []byte("# Test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		cmd := newActivateCmd()
		err := executeCommand(cmd, "to-activate")

		if err != nil {
			t.Fatalf("activate command error: %v", err)
		}

		// Verify state file
		stateFile := filepath.Join(ClaudeDir, ".current-profile")
		content, err := os.ReadFile(stateFile)
		if err != nil {
			t.Fatal(err)
		}
		if strings.TrimSpace(string(content)) != "to-activate" {
			t.Errorf("state file = %q, want %q", string(content), "to-activate")
		}
	})

	t.Run("activate non-existent profile", func(t *testing.T) {
		cmd := newActivateCmd()
		err := executeCommand(cmd, "non-existent")

		if err == nil {
			t.Error("activating non-existent profile should error")
		}
	})

	t.Run("activate without name", func(t *testing.T) {
		cmd := newActivateCmd()
		err := executeCommand(cmd)

		if err == nil {
			t.Error("activate without name should error")
		}
	})
}

func TestRestoreCmd(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Note: restore is interactive, so we just verify the command is set up correctly
	cmd := newRestoreCmd()

	// Verify basic command structure
	if cmd.Use != "restore" {
		t.Errorf("restore Use = %q, want %q", cmd.Use, "restore")
	}

	// Verify it takes no arguments (interactive)
	if cmd.Args == nil {
		t.Error("restore should have Args validation")
	}
}

func TestDiffCmd(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("diff with no profiles", func(t *testing.T) {
		cmd := newDiffCmd()
		err := executeCommand(cmd, "non-existent")

		if err == nil {
			t.Error("diff with non-existent profile should error")
		}
	})

	t.Run("diff two non-existent profiles", func(t *testing.T) {
		cmd := newDiffCmd()
		err := executeCommand(cmd, "profile1", "profile2")

		if err == nil {
			t.Error("diff with non-existent profiles should error")
		}
	})
}

func TestCheckBranchesCmd(t *testing.T) {
	cmd := newCheckBranchesCmd()

	// Verify aliases are set
	if len(cmd.Aliases) == 0 {
		t.Error("check-branches should have aliases")
	}

	hasBranches := false
	hasBr := false
	for _, alias := range cmd.Aliases {
		if alias == "branches" {
			hasBranches = true
		}
		if alias == "br" {
			hasBr = true
		}
	}

	if !hasBranches {
		t.Error("check-branches should have 'branches' alias")
	}
	if !hasBr {
		t.Error("check-branches should have 'br' alias")
	}
}

func TestSyncCmd(t *testing.T) {
	cmd := newSyncCmd()

	// Verify the --base flag exists
	flag := cmd.Flags().Lookup("base")
	if flag == nil {
		t.Error("sync should have --base flag")
	}

	// Verify short form
	shortFlag := cmd.Flags().ShorthandLookup("b")
	if shortFlag == nil {
		t.Error("sync should have -b short flag")
	}

	// Verify default value
	if flag.DefValue != "main" {
		t.Errorf("base flag default = %q, want %q", flag.DefValue, "main")
	}
}

func TestSwitchCmd(t *testing.T) {
	cmd := newSwitchCmd()

	// Verify aliases
	hasSelect := false
	for _, alias := range cmd.Aliases {
		if alias == "select" {
			hasSelect = true
			break
		}
	}

	if !hasSelect {
		t.Error("switch should have 'select' alias")
	}

	// Verify command Use
	if cmd.Use != "switch" {
		t.Errorf("switch Use = %q, want %q", cmd.Use, "switch")
	}
}

func TestEditCmd(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("edit without argument and no active profile", func(t *testing.T) {
		cmd := newEditCmd()
		err := executeCommand(cmd)

		// Should error because no active profile
		if err == nil {
			t.Error("edit without argument and no active profile should error")
		}
	})

	t.Run("edit command structure", func(t *testing.T) {
		cmd := newEditCmd()

		// Verify command accepts optional argument
		if cmd.Use != "edit [profile-name]" {
			t.Errorf("edit Use = %q, want 'edit [profile-name]'", cmd.Use)
		}
	})

	t.Run("edit non-existent profile", func(t *testing.T) {
		cmd := newEditCmd()
		err := executeCommand(cmd, "non-existent")

		// Should error because profile doesn't exist
		if err == nil {
			t.Error("edit non-existent profile should error")
		}
	})
}

func TestRootCmdVersion(t *testing.T) {
	// Verify version constant
	if Version == "" {
		t.Error("Version constant should not be empty")
	}

	if !strings.HasPrefix(Version, "1.0.0") {
		t.Errorf("Version = %q, expected to start with '1.0.0'", Version)
	}
}

func TestCommandRegistration(t *testing.T) {
	// Verify all commands are registered on rootCmd
	expectedCommands := []string{
		"version",
		"list",
		"show",
		"create",
		"delete",
		"edit",
		"activate",
		"switch",
		"restore",
		"check-branches",
		"sync",
		"diff",
	}

	registeredCommands := make(map[string]bool)
	for _, cmd := range rootCmd.Commands() {
		registeredCommands[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !registeredCommands[expected] {
			t.Errorf("Command %q should be registered", expected)
		}
	}
}
