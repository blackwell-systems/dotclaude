package hooks

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// builtInSessionInfo displays session start information
func builtInSessionInfo(r *Runner) error {
	fmt.Println("=== Claude Code Session Started ===")
	fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))

	// Working directory
	cwd, err := os.Getwd()
	if err == nil {
		fmt.Printf("Working directory: %s\n", cwd)
	}

	// Git branch info
	branch, err := getGitBranch()
	if err == nil && branch != "" {
		fmt.Printf("Git branch: %s\n", branch)

		// Check if branch is behind main
		if branch != "main" && branch != "master" && branch != "unknown" {
			behind := getCommitsBehind(branch)
			if behind > 0 {
				fmt.Printf("Warning: Branch is %d commits behind main - consider running: dotclaude sync\n", behind)
			}
		}
	}

	return nil
}

// builtInCheckDotclaude checks for .dotclaude file and profile mismatch
func builtInCheckDotclaude(r *Runner) error {
	// Check for .dotclaude file in current directory
	dotclaudePath := ".dotclaude"
	if _, err := os.Stat(dotclaudePath); os.IsNotExist(err) {
		return nil // No .dotclaude file, nothing to do
	}

	// Read desired profile from .dotclaude
	desiredProfile, err := readDotclaudeProfile(dotclaudePath)
	if err != nil {
		return nil // Can't read, skip silently
	}

	if desiredProfile == "" {
		return nil // No profile specified
	}

	// Validate profile name (security: prevent path traversal)
	if !isValidProfileName(desiredProfile) {
		fmt.Printf("\nWarning: Invalid profile name in .dotclaude: %s\n", desiredProfile)
		fmt.Println("   Profile names must contain only letters, numbers, hyphens, and underscores")
		return nil
	}

	// Check if profile exists
	profileDir := filepath.Join(r.RepoDir, "profiles", desiredProfile)
	if _, err := os.Stat(profileDir); os.IsNotExist(err) {
		fmt.Printf("\nWarning: Profile '%s' specified in .dotclaude not found\n", desiredProfile)
		fmt.Println("   Available profiles: dotclaude list")
		return nil
	}

	// Get current profile
	currentProfile := ""
	currentProfilePath := filepath.Join(r.ClaudeDir, ".current-profile")
	if data, err := os.ReadFile(currentProfilePath); err == nil {
		currentProfile = strings.TrimSpace(string(data))
	}

	// Compare profiles
	if desiredProfile != currentProfile {
		fmt.Println("")
		fmt.Println("+-------------------------------------------------------------+")
		fmt.Println("|  Profile Mismatch Detected                                  |")
		fmt.Println("+-------------------------------------------------------------+")
		fmt.Println("")
		fmt.Printf("  This project uses:    %s\n", desiredProfile)
		if currentProfile == "" {
			fmt.Println("  Currently active:     none")
		} else {
			fmt.Printf("  Currently active:     %s\n", currentProfile)
		}
		fmt.Println("")
		fmt.Printf("  To activate the project profile:\n")
		fmt.Printf("    dotclaude activate %s\n", desiredProfile)
		fmt.Println("")
	}

	return nil
}

// builtInGitTips provides helpful git workflow tips after git operations
func builtInGitTips(r *Runner) error {
	// This hook is triggered after Bash tool use
	// Check TOOL_USE_ARGS environment variable for git commands
	toolArgs := os.Getenv("TOOL_USE_ARGS")
	if toolArgs == "" {
		return nil
	}

	// Check if this was a git checkout or pull to main
	if strings.Contains(toolArgs, "git checkout main") ||
		strings.Contains(toolArgs, "git checkout master") ||
		strings.Contains(toolArgs, "git pull") {

		// Check if we're in a git repo
		if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err == nil {
			fmt.Println("Tip: Feature branches may be behind main. Run 'dotclaude sync' to check.")
		}
	}

	return nil
}

// readDotclaudeProfile reads the profile name from a .dotclaude file
func readDotclaudeProfile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// YAML format: profile: my-profile
		if strings.HasPrefix(line, "profile:") {
			value := strings.TrimPrefix(line, "profile:")
			value = strings.TrimSpace(value)
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			return value, nil
		}

		// Shell format: profile=my-profile
		if strings.HasPrefix(line, "profile=") {
			value := strings.TrimPrefix(line, "profile=")
			value = strings.TrimSpace(value)
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			return value, nil
		}
	}

	return "", nil
}

// isValidProfileName validates a profile name for security
func isValidProfileName(name string) bool {
	if name == "" {
		return false
	}
	// Only allow alphanumeric, hyphens, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched
}

// getGitBranch returns the current git branch name
func getGitBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getCommitsBehind returns how many commits the current branch is behind main
func getCommitsBehind(branch string) int {
	// Try main first, then master
	for _, baseBranch := range []string{"main", "master"} {
		cmd := exec.Command("git", "rev-list", "--count", fmt.Sprintf("%s..%s", branch, baseBranch))
		output, err := cmd.Output()
		if err == nil {
			var count int
			fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &count)
			return count
		}
	}
	return 0
}
