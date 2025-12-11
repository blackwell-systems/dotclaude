package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

// getEditor returns the appropriate editor for the current platform
func getEditor() string {
	// Check environment variables first (cross-platform)
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}

	// Platform-specific defaults
	if runtime.GOOS == "windows" {
		// Windows: try VS Code, then notepad
		if _, err := exec.LookPath("code"); err == nil {
			return "code --wait"
		}
		if _, err := exec.LookPath("notepad++"); err == nil {
			return "notepad++"
		}
		return "notepad"
	}

	// Unix: try common editors
	editors := []string{"code", "vim", "nano", "vi"}
	for _, editor := range editors {
		if _, err := exec.LookPath(editor); err == nil {
			if editor == "code" {
				return "code --wait"
			}
			return editor
		}
	}

	return "vi" // Ultimate fallback
}

func newEditCmd() *cobra.Command {
	var editSettings bool

	cmd := &cobra.Command{
		Use:   "edit [profile-name]",
		Short: "Edit a profile",
		Long: `Open a profile's CLAUDE.md file in your editor.

If no profile name is provided, edits the currently active profile.

Editor selection (in order of priority):
  1. $EDITOR environment variable
  2. $VISUAL environment variable
  3. Platform default (Windows: notepad, Unix: vim/nano)`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := profile.NewManager(RepoDir, ClaudeDir)

			var profileName string
			if len(args) == 0 {
				// Use active profile
				profileName = mgr.GetActiveProfileName()
				if profileName == "" {
					return fmt.Errorf("no active profile. Specify a profile name or activate one first")
				}
			} else {
				profileName = args[0]
			}

			// Check if profile exists
			if !mgr.ProfileExists(profileName) {
				return fmt.Errorf("profile '%s' does not exist", profileName)
			}

			// Get editor (cross-platform)
			editor := getEditor()

			// Determine file to edit
			var filePath string
			if editSettings {
				filePath = filepath.Join(ProfilesDir, profileName, "settings.json")
			} else {
				filePath = filepath.Join(ProfilesDir, profileName, "CLAUDE.md")
			}

			// Check if file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", filePath)
			}

			// Open editor
			fmt.Printf("Opening %s in %s...\n", filepath.Base(filePath), editor)

			// Parse editor command (may have arguments like "code --wait")
			parts := strings.Fields(editor)
			editorBin := parts[0]
			editorArgs := append(parts[1:], filePath)

			editorCmd := exec.Command(editorBin, editorArgs...)
			editorCmd.Stdin = os.Stdin
			editorCmd.Stdout = os.Stdout
			editorCmd.Stderr = os.Stderr

			if err := editorCmd.Run(); err != nil {
				return fmt.Errorf("editor exited with error: %w", err)
			}

			fmt.Println()
			fmt.Println("╭─────────────────────────────────────────────────────────────╮")
			fmt.Println("│  ✓ Edit Complete                                            │")
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()

			return nil
		},
	}

	cmd.Flags().BoolVarP(&editSettings, "settings", "s", false, "edit settings.json instead of CLAUDE.md")

	return cmd
}
