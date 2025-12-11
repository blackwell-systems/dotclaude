package cli

import (
	"fmt"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "create <profile-name>",
		Aliases: []string{"new"},
		Short:   "Create a new profile",
		Long:    "Create a new dotclaude profile from the template.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			mgr := profile.NewManager(RepoDir, ClaudeDir)

			// Create the profile
			if err := mgr.Create(profileName); err != nil {
				return err
			}

			// Success message
			fmt.Println()
			fmt.Println("╭─────────────────────────────────────────────────────────────╮")
			fmt.Printf("│  ✓ Profile Created: %-39s│\n", Green(profileName))
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()
			fmt.Printf("Profile created at: %s/profiles/%s\n", RepoDir, profileName)
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Printf("  1. Edit profile:    dotclaude edit %s\n", profileName)
			fmt.Printf("  2. Activate it:     dotclaude activate %s\n", profileName)
			fmt.Println()

			return nil
		},
	}
}
