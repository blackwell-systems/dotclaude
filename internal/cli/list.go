package cli

import (
	"fmt"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List all profiles",
		Long:    "Display all available dotclaude profiles with their status.",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := profile.NewManager(RepoDir, ClaudeDir)

			profiles, err := mgr.ListProfiles()
			if err != nil {
				return fmt.Errorf("failed to list profiles: %w", err)
			}

			if len(profiles) == 0 {
				fmt.Println("No profiles found.")
				fmt.Printf("\nCreate your first profile:\n")
				fmt.Printf("  dotclaude create my-project\n")
				return nil
			}

			// Print header
			fmt.Println("\n╭─────────────────────────────────────────────────────────────╮")
			fmt.Println("│  Available Profiles                                         │")
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()

			// Print profiles
			for _, p := range profiles {
				if p.IsActive {
					fmt.Printf("  ▶ \033[1;32m%s\033[0m (active)\n", p.Name)
				} else {
					fmt.Printf("    %s\n", p.Name)
				}
			}

			fmt.Println()
			fmt.Printf("Total: %d profile(s)\n", len(profiles))
			fmt.Println()

			return nil
		},
	}
}

func init() {
	// This will be called when root.go initializes
}
