package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "delete <profile-name>",
		Short:   "Delete a profile",
		Long:    "Delete a dotclaude profile permanently.",
		Aliases: []string{"rm", "remove"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			mgr := profile.NewManager(RepoDir, ClaudeDir)

			// Check if profile exists
			if !mgr.ProfileExists(profileName) {
				return fmt.Errorf("profile '%s' does not exist", profileName)
			}

			// Confirm deletion unless --force
			if !force {
				fmt.Printf("\nAre you sure you want to delete profile '%s'? [y/N]: ", profileName)
				reader := bufio.NewReader(os.Stdin)
				response, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				response = strings.TrimSpace(strings.ToLower(response))
				if response != "y" && response != "yes" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			// Delete the profile
			if err := mgr.Delete(profileName); err != nil {
				return err
			}

			// Success message
			fmt.Println()
			fmt.Println("╭─────────────────────────────────────────────────────────────╮")
			fmt.Printf("│  ✓ Profile Deleted: %-43s│\n", profileName)
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "skip confirmation prompt")

	return cmd
}
