// Package hooks provides a cross-platform hook runner system for dotclaude.
//
// Hooks are executable scripts or built-in commands that run at specific points
// in the Claude Code lifecycle (e.g., session start, post-tool use).
//
// Hook Directory Structure:
//
//	~/.claude/hooks/
//	├── session-start/
//	│   ├── 00-session-info      # Runs first (lower number = higher priority)
//	│   ├── 10-check-dotclaude
//	│   └── 50-custom.sh         # User-added hooks
//	├── post-tool-bash/
//	│   └── 10-git-tips
//	└── pre-tool-edit/
//	    └── ...
//
// Hooks are executed in alphabetical order (use numeric prefixes for ordering).
// Any executable file in the hook directory will be run.
package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// HookType represents the type of hook event
type HookType string

const (
	HookSessionStart  HookType = "session-start"
	HookPostToolBash  HookType = "post-tool-bash"
	HookPostToolEdit  HookType = "post-tool-edit"
	HookPreToolBash   HookType = "pre-tool-bash"
	HookPreToolEdit   HookType = "pre-tool-edit"
)

// Runner executes hooks for a given hook type
type Runner struct {
	HooksDir   string            // Directory containing hook subdirectories
	ClaudeDir  string            // ~/.claude directory
	RepoDir    string            // DOTCLAUDE_REPO_DIR
	BuiltIns   map[HookType][]BuiltInHook // Built-in hooks by type
	Env        map[string]string // Additional environment variables
}

// BuiltInHook represents a hook implemented in Go
type BuiltInHook struct {
	Name     string
	Priority int // Lower runs first
	Run      func(r *Runner) error
}

// NewRunner creates a new hook runner
func NewRunner(claudeDir, repoDir string) *Runner {
	hooksDir := filepath.Join(claudeDir, "hooks")

	r := &Runner{
		HooksDir:  hooksDir,
		ClaudeDir: claudeDir,
		RepoDir:   repoDir,
		BuiltIns:  make(map[HookType][]BuiltInHook),
		Env:       make(map[string]string),
	}

	// Register built-in hooks
	r.registerBuiltIns()

	return r
}

// registerBuiltIns registers all built-in hooks
func (r *Runner) registerBuiltIns() {
	// Session start hooks
	r.BuiltIns[HookSessionStart] = []BuiltInHook{
		{Name: "session-info", Priority: 0, Run: builtInSessionInfo},
		{Name: "check-dotclaude", Priority: 10, Run: builtInCheckDotclaude},
	}

	// Post-tool bash hooks
	r.BuiltIns[HookPostToolBash] = []BuiltInHook{
		{Name: "git-tips", Priority: 10, Run: builtInGitTips},
	}
}

// Run executes all hooks for the given hook type
func (r *Runner) Run(hookType HookType) error {
	var allHooks []hookEntry

	// Collect built-in hooks
	for _, builtin := range r.BuiltIns[hookType] {
		allHooks = append(allHooks, hookEntry{
			name:     fmt.Sprintf("%02d-%s", builtin.Priority, builtin.Name),
			priority: builtin.Priority,
			builtin:  &builtin,
		})
	}

	// Collect external hooks from directory
	hookDir := filepath.Join(r.HooksDir, string(hookType))
	if entries, err := os.ReadDir(hookDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			path := filepath.Join(hookDir, name)

			// Check if executable
			if !isExecutable(path) {
				continue
			}

			// Extract priority from name (e.g., "10-myhook" -> priority 10)
			priority := extractPriority(name)

			allHooks = append(allHooks, hookEntry{
				name:     name,
				priority: priority,
				path:     path,
			})
		}
	}

	// Sort by priority (name includes priority prefix, so alphabetical works)
	sort.Slice(allHooks, func(i, j int) bool {
		return allHooks[i].name < allHooks[j].name
	})

	// Execute hooks in order
	for _, hook := range allHooks {
		if hook.builtin != nil {
			if err := hook.builtin.Run(r); err != nil {
				// Log but don't fail - hooks shouldn't break the session
				fmt.Fprintf(os.Stderr, "Hook %s warning: %v\n", hook.name, err)
			}
		} else {
			if err := r.runExternal(hook.path); err != nil {
				fmt.Fprintf(os.Stderr, "Hook %s warning: %v\n", hook.name, err)
			}
		}
	}

	return nil
}

// hookEntry represents either a built-in or external hook
type hookEntry struct {
	name     string
	priority int
	builtin  *BuiltInHook
	path     string
}

// runExternal executes an external hook script
func (r *Runner) runExternal(path string) error {
	var cmd *exec.Cmd

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".ps1":
		// PowerShell script - try pwsh first (cross-platform), fall back to powershell (Windows)
		if pwsh, err := exec.LookPath("pwsh"); err == nil {
			cmd = exec.Command(pwsh, "-ExecutionPolicy", "Bypass", "-File", path)
		} else if powershell, err := exec.LookPath("powershell"); err == nil {
			cmd = exec.Command(powershell, "-ExecutionPolicy", "Bypass", "-File", path)
		} else {
			return fmt.Errorf("PowerShell not found - install PowerShell Core (pwsh) to run .ps1 hooks")
		}
	case ".sh", ".bash":
		// Shell script - check if bash is available
		if bash, err := exec.LookPath("bash"); err == nil {
			cmd = exec.Command(bash, path)
		} else if runtime.GOOS == "windows" {
			// On Windows, suggest Git Bash or WSL
			return fmt.Errorf("bash not found - install Git for Windows or WSL to run .sh hooks")
		} else {
			// On Unix, try sh as fallback
			cmd = exec.Command("sh", path)
		}
	case ".cmd", ".bat":
		// Windows batch scripts
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", path)
		} else {
			return fmt.Errorf(".cmd/.bat hooks only work on Windows")
		}
	default:
		// Try to execute directly (works for shebang scripts on Unix, .exe on Windows)
		cmd = exec.Command(path)
	}

	// Set up environment
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOTCLAUDE_REPO_DIR=%s", r.RepoDir))
	cmd.Env = append(cmd.Env, fmt.Sprintf("CLAUDE_DIR=%s", r.ClaudeDir))
	for k, v := range r.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Connect to stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// List returns information about all hooks for a given type
func (r *Runner) List(hookType HookType) []HookInfo {
	var hooks []HookInfo

	// Built-in hooks
	for _, builtin := range r.BuiltIns[hookType] {
		hooks = append(hooks, HookInfo{
			Name:     builtin.Name,
			Priority: builtin.Priority,
			Type:     "built-in",
			Enabled:  true,
		})
	}

	// External hooks
	hookDir := filepath.Join(r.HooksDir, string(hookType))
	if entries, err := os.ReadDir(hookDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			path := filepath.Join(hookDir, name)

			hooks = append(hooks, HookInfo{
				Name:     name,
				Priority: extractPriority(name),
				Type:     "external",
				Path:     path,
				Enabled:  isExecutable(path),
			})
		}
	}

	// Sort by priority
	sort.Slice(hooks, func(i, j int) bool {
		if hooks[i].Priority != hooks[j].Priority {
			return hooks[i].Priority < hooks[j].Priority
		}
		return hooks[i].Name < hooks[j].Name
	})

	return hooks
}

// HookInfo contains information about a hook
type HookInfo struct {
	Name     string
	Priority int
	Type     string // "built-in" or "external"
	Path     string // Empty for built-in
	Enabled  bool
}

// GetHookTypes returns all supported hook types
func GetHookTypes() []HookType {
	return []HookType{
		HookSessionStart,
		HookPostToolBash,
		HookPostToolEdit,
		HookPreToolBash,
		HookPreToolEdit,
	}
}

// isExecutable checks if a file is executable
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		// On Windows, check for executable extensions
		ext := strings.ToLower(filepath.Ext(path))
		return ext == ".exe" || ext == ".cmd" || ext == ".bat" || ext == ".ps1" || ext == ".sh"
	}

	// On Unix, check execute permission
	return info.Mode()&0111 != 0
}

// extractPriority extracts numeric priority from hook name
// e.g., "10-myhook.sh" -> 10, "myhook.sh" -> 50 (default)
func extractPriority(name string) int {
	if len(name) >= 2 && name[0] >= '0' && name[0] <= '9' && name[1] >= '0' && name[1] <= '9' {
		priority := int(name[0]-'0')*10 + int(name[1]-'0')
		return priority
	}
	return 50 // Default priority
}

// EnsureHooksDir creates the hooks directory structure
func (r *Runner) EnsureHooksDir() error {
	for _, hookType := range GetHookTypes() {
		dir := filepath.Join(r.HooksDir, string(hookType))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create hooks directory %s: %w", dir, err)
		}
	}
	return nil
}
