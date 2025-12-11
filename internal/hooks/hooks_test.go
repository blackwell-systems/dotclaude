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

func TestRunExternalErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "hooks-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	runner := NewRunner(tmpDir, tmpDir)

	t.Run("cmd on non-windows", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping on Windows")
		}

		cmdFile := filepath.Join(tmpDir, "test.cmd")
		if err := os.WriteFile(cmdFile, []byte("echo test"), 0755); err != nil {
			t.Fatal(err)
		}

		err := runner.runExternal(cmdFile)
		if err == nil {
			t.Error("Expected error running .cmd on non-Windows")
		}
	})

	t.Run("bat on non-windows", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping on Windows")
		}

		batFile := filepath.Join(tmpDir, "test.bat")
		if err := os.WriteFile(batFile, []byte("echo test"), 0755); err != nil {
			t.Fatal(err)
		}

		err := runner.runExternal(batFile)
		if err == nil {
			t.Error("Expected error running .bat on non-Windows")
		}
	})

	t.Run("non-existent script", func(t *testing.T) {
		err := runner.runExternal("/nonexistent/script.sh")
		if err == nil {
			t.Error("Expected error running non-existent script")
		}
	})
}

func TestHookEntry(t *testing.T) {
	t.Run("hook entry with builtin", func(t *testing.T) {
		builtIn := BuiltInHook{
			Name:     "test-hook",
			Priority: 10,
			Run:      func(r *Runner) error { return nil },
		}

		entry := hookEntry{
			name:     "10-test-hook",
			priority: 10,
			builtin:  &builtIn,
		}

		if entry.name != "10-test-hook" {
			t.Errorf("Expected name '10-test-hook', got %q", entry.name)
		}
		if entry.priority != 10 {
			t.Errorf("Expected priority 10, got %d", entry.priority)
		}
		if entry.builtin == nil {
			t.Error("Expected builtin to be set")
		}
	})

	t.Run("hook entry with external", func(t *testing.T) {
		entry := hookEntry{
			name:     "20-custom.sh",
			priority: 20,
			path:     "/path/to/20-custom.sh",
		}

		if entry.name != "20-custom.sh" {
			t.Errorf("Expected name '20-custom.sh', got %q", entry.name)
		}
		if entry.path != "/path/to/20-custom.sh" {
			t.Errorf("Expected path '/path/to/20-custom.sh', got %q", entry.path)
		}
	})
}

func TestHookInfo(t *testing.T) {
	t.Run("built-in hook info", func(t *testing.T) {
		info := HookInfo{
			Name:     "session-info",
			Priority: 0,
			Type:     "built-in",
			Enabled:  true,
		}

		if info.Name != "session-info" {
			t.Errorf("Expected Name 'session-info', got %q", info.Name)
		}
		if info.Type != "built-in" {
			t.Errorf("Expected Type 'built-in', got %q", info.Type)
		}
		if !info.Enabled {
			t.Error("Expected Enabled to be true")
		}
	})

	t.Run("external hook info", func(t *testing.T) {
		info := HookInfo{
			Name:     "20-custom.sh",
			Priority: 20,
			Type:     "external",
			Path:     "/path/to/hook",
			Enabled:  true,
		}

		if info.Path != "/path/to/hook" {
			t.Errorf("Expected Path '/path/to/hook', got %q", info.Path)
		}
	})
}

func TestRunWithCustomHook(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping shell script test on Windows")
	}

	tmpDir, err := os.MkdirTemp("", "hooks-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	claudeDir := filepath.Join(tmpDir, ".claude")
	runner := NewRunner(claudeDir, tmpDir)

	// Ensure hooks dir
	if err := runner.EnsureHooksDir(); err != nil {
		t.Fatal(err)
	}

	// Create a simple hook that creates a marker file
	hookDir := filepath.Join(claudeDir, "hooks", "pre-tool-edit")
	markerFile := filepath.Join(tmpDir, "hook-ran")
	hookScript := filepath.Join(hookDir, "50-marker.sh")

	script := "#!/bin/bash\ntouch " + markerFile
	if err := os.WriteFile(hookScript, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	// Run the hooks
	if err := runner.Run(HookPreToolEdit); err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// Check if marker file was created
	if _, err := os.Stat(markerFile); os.IsNotExist(err) {
		t.Error("Expected hook to create marker file")
	}
}

func TestListEmptyHookType(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "hooks-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	claudeDir := filepath.Join(tmpDir, ".claude")
	runner := NewRunner(claudeDir, tmpDir)

	// List hooks for a type with no external hooks (but has built-ins)
	hooks := runner.List(HookPostToolEdit)

	// Should only have built-in hooks (if any)
	for _, h := range hooks {
		if h.Type != "built-in" {
			t.Errorf("Expected only built-in hooks, got type %q", h.Type)
		}
	}
}

func TestRunnerEnv(t *testing.T) {
	runner := NewRunner("/tmp/claude", "/tmp/repo")

	// Test setting environment variables
	runner.Env["TEST_VAR"] = "test_value"

	if runner.Env["TEST_VAR"] != "test_value" {
		t.Errorf("Expected TEST_VAR to be 'test_value', got %q", runner.Env["TEST_VAR"])
	}
}
