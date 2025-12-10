package cli

import (
	"fmt"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

func newActivateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "activate <profile-name>",
		Short:   "Activate a profile",
		Long:    "Activate a dotclaude profile by merging base + profile configuration.",
		Aliases: []string{"use"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			mgr := profile.NewManager(RepoDir, ClaudeDir)

			// Check if profile exists
			if !mgr.ProfileExists(profileName) {
				return fmt.Errorf("profile '%s' does not exist", profileName)
			}

			// Get current active profile
			currentProfile := mgr.GetActiveProfileName()

			// Show status message
			fmt.Println()
			fmt.Println("╭─────────────────────────────────────────────────────────────╮")
			if currentProfile == profileName {
				fmt.Printf("│  Updating Profile: %-43s│\n", profileName)
			} else {
				fmt.Printf("│  Activating Profile: %-41s│\n", profileName)
			}
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()

			// Show backup message if switching profiles
			if currentProfile != "" && currentProfile != profileName {
				fmt.Printf("  [1/3] Backing up existing configuration...\n")
			} else if currentProfile == profileName {
				fmt.Printf("  [1/3] Already on '%s', updating in place\n", profileName)
			} else {
				fmt.Printf("  [1/3] No existing configuration to backup\n")
			}

			// Activate the profile
			if err := mgr.Activate(profileName); err != nil {
				return err
			}

			fmt.Printf("  [2/3] Merged base + profile configuration\n")
			fmt.Printf("  [3/3] Applied profile settings\n")

			// Success message
			fmt.Println()
			fmt.Println("╭─────────────────────────────────────────────────────────────╮")
			fmt.Printf("│  ✓ Profile Activated: %-41s│\n", profileName)
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()
			fmt.Printf("  Configuration deployed to: %s\n", ClaudeDir)
			fmt.Println()
			fmt.Println("  Verify with:")
			fmt.Println("    • dotclaude show")
			fmt.Printf("    • cat %s/CLAUDE.md\n", ClaudeDir)
			fmt.Println()

			return nil
		},
	}

	return cmd
}
