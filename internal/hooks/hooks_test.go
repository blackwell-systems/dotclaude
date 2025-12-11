package hooks

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewRunner(t *testing.T) {
	runner := NewRunner("/tmp/claude", "/tmp/repo")

	if runner.ClaudeDir != "/tmp/claude" {
		t.Errorf("ClaudeDir = %q, want %q", runner.ClaudeDir, "/tmp/claude")
	}

	if runner.RepoDir != "/tmp/repo" {
		t.Errorf("RepoDir = %q, want %q", runner.RepoDir, "/tmp/repo")
	}

	if runner.HooksDir != "/tmp/claude/hooks" {
		t.Errorf("HooksDir = %q, want %q", runner.HooksDir, "/tmp/claude/hooks")
	}

	// Check built-in hooks are registered
	if len(runner.BuiltIns[HookSessionStart]) == 0 {
		t.Error("Expected session-start built-in hooks to be registered")
	}

	if len(runner.BuiltIns[HookPostToolBash]) == 0 {
		t.Error("Expected post-tool-bash built-in hooks to be registered")
	}
}

func TestGetHookTypes(t *testing.T) {
	types := GetHookTypes()

	if len(types) != 5 {
		t.Errorf("Expected 5 hook types, got %d", len(types))
	}

	expected := map[HookType]bool{
		HookSessionStart: true,
		HookPostToolBash: true,
		HookPostToolEdit: true,
		HookPreToolBash:  true,
		HookPreToolEdit:  true,
	}

	for _, ht := range types {
		if !expected[ht] {
			t.Errorf("Unexpected hook type: %s", ht)
		}
	}
}

func TestExtractPriority(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"two digit prefix", "00-first.sh", 0},
		{"ten prefix", "10-second.sh", 10},
		{"fifty prefix", "50-middle.sh", 50},
		{"ninety-nine prefix", "99-last.sh", 99},
		{"no prefix", "myhook.sh", 50},
		{"single digit", "5-hook.sh", 50}, // Not valid two-digit, defaults to 50
		{"letters", "ab-hook.sh", 50},
		{"empty", "", 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPriority(tt.input)
			if got != tt.expected {
				t.Errorf("extractPriority(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsExecutable(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "hooks-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a non-executable file
	nonExecFile := filepath.Join(tmpDir, "nonexec.txt")
	if err := os.WriteFile(nonExecFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create an executable file (Unix) or .exe (Windows)
	var execFile string
	if runtime.GOOS == "windows" {
		execFile = filepath.Join(tmpDir, "exec.exe")
		if err := os.WriteFile(execFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	} else {
		execFile = filepath.Join(tmpDir, "exec.sh")
		if err := os.WriteFile(execFile, []byte("#!/bin/bash\necho test"), 0755); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("non-executable file", func(t *testing.T) {
		if isExecutable(nonExecFile) {
			t.Error("Expected non-executable file to return false")
		}
	})

	t.Run("executable file", func(t *testing.T) {
		if !isExecutable(execFile) {
			t.Error("Expected executable file to return true")
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		if isExecutable("/nonexistent/file") {
			t.Error("Expected non-existent file to return false")
		}
	})
}

func TestEnsureHooksDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "hooks-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	claudeDir := filepath.Join(tmpDir, ".claude")
	runner := NewRunner(claudeDir, tmpDir)

	if err := runner.EnsureHooksDir(); err != nil {
		t.Fatalf("EnsureHooksDir() error = %v", err)
	}

	// Verify all hook type directories were created
	for _, ht := range GetHookTypes() {
		dir := filepath.Join(claudeDir, "hooks", string(ht))
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to exist", dir)
		}
	}
}

func TestList(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "hooks-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	claudeDir := filepath.Join(tmpDir, ".claude")
	runner := NewRunner(claudeDir, tmpDir)

	// Ensure hooks dir exists
	if err := runner.EnsureHooksDir(); err != nil {
		t.Fatal(err)
	}

	// Create a custom hook
	hookDir := filepath.Join(claudeDir, "hooks", "session-start")
	customHook := filepath.Join(hookDir, "20-custom.sh")
	if err := os.WriteFile(customHook, []byte("#!/bin/bash\necho custom"), 0755); err != nil {
		t.Fatal(err)
	}

	hooks := runner.List(HookSessionStart)

	// Should have built-in hooks + custom hook
	if len(hooks) < 2 {
		t.Errorf("Expected at least 2 hooks, got %d", len(hooks))
	}

	// Check that built-in hooks are listed
	hasBuiltIn := false
	hasCustom := false
	for _, h := range hooks {
		if h.Type == "built-in" && h.Name == "session-info" {
			hasBuiltIn = true
		}
		if h.Type == "external" && h.Name == "20-custom.sh" {
			hasCustom = true
		}
	}

	if !hasBuiltIn {
		t.Error("Expected built-in session-info hook to be listed")
	}

	if !hasCustom {
		t.Error("Expected custom hook to be listed")
	}
}

func TestRun(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "hooks-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	claudeDir := filepath.Join(tmpDir, ".claude")
	runner := NewRunner(claudeDir, tmpDir)

	// Run should not error even with no hooks dir
	if err := runner.Run(HookPreToolEdit); err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// Ensure hooks dir and run again
	if err := runner.EnsureHooksDir(); err != nil {
		t.Fatal(err)
	}

	if err := runner.Run(HookPreToolEdit); err != nil {
		t.Errorf("Run() with hooks dir error = %v", err)
	}
}

func TestReadDotclaudeProfile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "hooks-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{"yaml format", "profile: my-profile", "my-profile"},
		{"yaml with quotes", "profile: \"my-profile\"", "my-profile"},
		{"yaml with single quotes", "profile: 'my-profile'", "my-profile"},
		{"shell format", "profile=my-profile", "my-profile"},
		{"shell with quotes", "profile=\"my-profile\"", "my-profile"},
		{"with comments", "# comment\nprofile: my-profile", "my-profile"},
		{"empty file", "", ""},
		{"no profile", "other: value", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, ".dotclaude-"+tt.name)
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			got, err := readDotclaudeProfile(path)
			if err != nil {
				t.Fatalf("readDotclaudeProfile() error = %v", err)
			}

			if got != tt.expected {
				t.Errorf("readDotclaudeProfile() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestIsValidProfileName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple", "myprofile", true},
		{"with hyphen", "my-profile", true},
		{"with underscore", "my_profile", true},
		{"with numbers", "profile123", true},
		{"mixed", "My-Profile_123", true},
		{"empty", "", false},
		{"with slash", "my/profile", false},
		{"with dots", "my.profile", false},
		{"path traversal", "../etc", false},
		{"with spaces", "my profile", false},
		{"with special", "my@profile", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidProfileName(tt.input)
			if got != tt.expected {
				t.Errorf("isValidProfileName(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
