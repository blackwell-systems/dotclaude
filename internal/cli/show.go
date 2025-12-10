package cli

import (
	"fmt"
	"os"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show active profile",
		Long:  "Display information about the currently active profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := profile.NewManager(RepoDir, ClaudeDir)

			activeProfile, err := mgr.GetActiveProfile()
			if err != nil {
				return fmt.Errorf("failed to get active profile: %w", err)
			}

			if activeProfile == nil {
				fmt.Println("\n╭─────────────────────────────────────────────────────────────╮")
				fmt.Println("│  No Active Profile                                          │")
				fmt.Println("╰─────────────────────────────────────────────────────────────╯")
				fmt.Println()
				fmt.Println("No profile is currently active.")
				fmt.Println()
				fmt.Println("Activate a profile:")
				fmt.Println("  dotclaude activate <profile-name>")
				fmt.Println()
				fmt.Println("Or create a new profile:")
				fmt.Println("  dotclaude create <profile-name>")
				fmt.Println()
				return nil
			}

			// Display active profile info
			fmt.Println("\n╭─────────────────────────────────────────────────────────────╮")
			fmt.Println("│  Active Profile                                             │")
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()
			fmt.Printf("  Profile:  \033[1;32m%s\033[0m\n", activeProfile.Name)
			fmt.Printf("  Location: %s\n", activeProfile.Path)
			fmt.Printf("  Modified: %s\n", activeProfile.LastModified.Format("2006-01-02 15:04:05"))
			fmt.Println()

			// Check if Claude directory exists
			if _, err := os.Stat(ClaudeDir); err == nil {
				fmt.Println("  Status:   ✓ Claude directory configured")
			} else {
				fmt.Println("  Status:   ⚠ Claude directory not found")
			}

			fmt.Println()

			return nil
		},
	}
}
