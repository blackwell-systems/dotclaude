package cli

import (
	"fmt"
	"os"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

func newActivateCmd() *cobra.Command {
	var dryRun bool
	var verbose bool

	cmd := &cobra.Command{
		Use:     "activate <profile-name>",
		Short:   "Activate a profile",
		Long:    "Activate a dotclaude profile by merging base + profile configuration.",
		Aliases: []string{"use"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("profile name is required")
			}
			if len(args) > 1 {
				return fmt.Errorf("only one profile name allowed, got %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			// Get flag values (handles aliases)
			preview, _ := cmd.Flags().GetBool("preview")
			debug, _ := cmd.Flags().GetBool("debug")
			if preview {
				dryRun = true
			}
			if debug {
				verbose = true
			}

			mgr := profile.NewManager(RepoDir, ClaudeDir)

			// Validate profile name
			if err := profile.ValidateProfileName(profileName); err != nil {
				return fmt.Errorf("Invalid profile name: %w", err)
			}

			// Check if profile exists
			if !mgr.ProfileExists(profileName) {
				return fmt.Errorf("profile '%s' not found", profileName)
			}

			// Get current active profile
			currentProfile := mgr.GetActiveProfileName()

			// Handle dry-run mode
			if dryRun {
				return showPreview(mgr, profileName, currentProfile, verbose)
			}

			// Handle verbose mode
			if verbose {
				fmt.Printf("[DEBUG] RepoDir: %s\n", RepoDir)
				fmt.Printf("[DEBUG] ClaudeDir: %s\n", ClaudeDir)
				fmt.Printf("[DEBUG] ProfilesDir: %s\n", mgr.ProfilesDir)
				fmt.Printf("[DEBUG] Current profile: %s\n", currentProfile)
				fmt.Printf("[DEBUG] Target profile: %s\n", profileName)
				fmt.Println()
			}

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

	// Add flags
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show preview without making changes")
	cmd.Flags().Bool("preview", false, "Alias for --dry-run")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show debug output")
	cmd.Flags().Bool("debug", false, "Alias for --verbose")

	return cmd
}

// showPreview displays what would happen without making changes
func showPreview(mgr *profile.Manager, profileName, currentProfile string, verbose bool) error {
	fmt.Println()
	fmt.Println("╭─────────────────────────────────────────────────────────────╮")
	fmt.Printf("│  DRY RUN - Preview Mode                                     │\n")
	fmt.Println("╰─────────────────────────────────────────────────────────────╯")
	fmt.Println()

	fmt.Printf("Would activate profile: %s\n", profileName)
	fmt.Println()

	// Show current state
	if currentProfile != "" {
		fmt.Printf("Current profile: %s\n", currentProfile)
		if currentProfile != profileName {
			fmt.Println("Action: Switch profiles (backup will be created)")
		} else {
			fmt.Println("Action: Update current profile in place")
		}
	} else {
		fmt.Println("Current profile: None")
		fmt.Println("Action: First activation")
	}
	fmt.Println()

	// Show what would be merged
	fmt.Println("Files that would be merged:")
	fmt.Printf("  • base/CLAUDE.md\n")
	fmt.Printf("  • profiles/%s/CLAUDE.md\n", profileName)
	fmt.Println()

	// Show settings
	fmt.Println("Settings:")
	profileSettingsPath := fmt.Sprintf("%s/profiles/%s/settings.json", mgr.RepoDir, profileName)
	baseSettingsPath := fmt.Sprintf("%s/base/settings.json", mgr.RepoDir)

	if fileExists(profileSettingsPath) {
		fmt.Printf("  • Would use profile settings: profiles/%s/settings.json\n", profileName)
	} else if fileExists(baseSettingsPath) {
		fmt.Println("  • Would use base settings: base/settings.json")
	} else {
		fmt.Println("  • No settings found")
	}
	fmt.Println()

	// Show verbose details if requested
	if verbose {
		fmt.Println("[DEBUG] Preview Details:")
		fmt.Printf("[DEBUG]   Profile path: %s/profiles/%s\n", mgr.RepoDir, profileName)
		fmt.Printf("[DEBUG]   Claude dir: %s\n", mgr.ClaudeDir)
		fmt.Printf("[DEBUG]   Would create backup: %v\n", currentProfile != "" && currentProfile != profileName)
		fmt.Println()
	}

	fmt.Println("╭─────────────────────────────────────────────────────────────╮")
	fmt.Println("│  No changes made (dry-run mode)                            │")
	fmt.Println("╰─────────────────────────────────────────────────────────────╯")
	fmt.Println()
	fmt.Println("To apply these changes, run without --dry-run:")
	fmt.Printf("  dotclaude activate %s\n", profileName)
	fmt.Println()

	return nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
