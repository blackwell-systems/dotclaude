package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	var editSettings bool

	cmd := &cobra.Command{
		Use:   "edit [profile-name]",
		Short: "Edit a profile",
		Long: `Open a profile's CLAUDE.md file in your editor.

If no profile name is provided, edits the currently active profile.`,
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

			// Get editor
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim" // Default to vim
			}

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

			editorCmd := exec.Command(editor, filePath)
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
