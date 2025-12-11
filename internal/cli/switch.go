package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

func newSwitchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "switch",
		Aliases: []string{"select"},
		Short:   "Interactively switch between profiles",
		Long: `Display a numbered list of profiles and switch to the selected one.

This is an interactive command that prompts you to choose a profile.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := profile.NewManager(RepoDir, ClaudeDir)

			// Get all profiles
			profiles, err := mgr.ListProfiles()
			if err != nil {
				return fmt.Errorf("failed to list profiles: %w", err)
			}

			if len(profiles) == 0 {
				fmt.Println("No profiles found.")
				fmt.Println()
				fmt.Println("Create your first profile:")
				fmt.Println("  dotclaude create my-project")
				return nil
			}

			// Print header
			fmt.Println()
			fmt.Println("╭─────────────────────────────────────────────────────────────╮")
			fmt.Println("│  Select a Profile                                           │")
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()

			// Display profiles with numbers
			for i, p := range profiles {
				if p.IsActive {
					fmt.Printf("  [%d] \033[1;32m%s\033[0m (active)\n", i+1, p.Name)
				} else {
					fmt.Printf("  [%d] %s\n", i+1, p.Name)
				}
			}

			fmt.Println()

			// Prompt for selection
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter profile number (or 'q' to quit): ")
			choice, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			choice = strings.TrimSpace(choice)

			// Handle quit
			if choice == "q" || choice == "Q" {
				fmt.Println("Cancelled.")
				return nil
			}

			// Parse selection
			selection, err := strconv.Atoi(choice)
			if err != nil || selection < 1 || selection > len(profiles) {
				return fmt.Errorf("invalid selection: %s", choice)
			}

			selectedProfile := profiles[selection-1]

			// Check if already active
			if selectedProfile.IsActive {
				fmt.Printf("\nProfile '%s' is already active.\n", selectedProfile.Name)
				return nil
			}

			// Activate the selected profile
			fmt.Println()
			fmt.Printf("Activating profile: %s\n", selectedProfile.Name)
			fmt.Println()

			if err := mgr.Activate(selectedProfile.Name); err != nil {
				return fmt.Errorf("failed to activate profile: %w", err)
			}

			// Success message
			fmt.Println("╭─────────────────────────────────────────────────────────────╮")
			fmt.Printf("│  ✓ Switched to: %-43s│\n", selectedProfile.Name)
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()

			return nil
		},
	}

	return cmd
}
