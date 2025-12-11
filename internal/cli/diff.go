package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/blackwell-systems/dotclaude/internal/profile"
	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff [profile1] [profile2]",
		Short: "Compare two profiles",
		Long: `Compare the CLAUDE.md files of two profiles.

If only one profile is provided, compares it to the currently active profile.
If no profiles are provided, shows an error.`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := profile.NewManager(RepoDir, ClaudeDir)

			var profile1Path, profile2Path string
			var profile1Name, profile2Name string

			if len(args) == 1 {
				// Compare provided profile with current active profile
				activeProfile, err := mgr.GetActiveProfile()
				if err != nil {
					return fmt.Errorf("failed to get active profile: %w", err)
				}
				if activeProfile == nil {
					return fmt.Errorf("no active profile to compare against")
				}

				profile1Name = activeProfile.Name
				profile1Path = fmt.Sprintf("%s/profiles/%s/CLAUDE.md", RepoDir, activeProfile.Name)
				profile2Name = args[0]
				profile2Path = fmt.Sprintf("%s/profiles/%s/CLAUDE.md", RepoDir, args[0])
			} else {
				// Compare two specified profiles
				profile1Name = args[0]
				profile1Path = fmt.Sprintf("%s/profiles/%s/CLAUDE.md", RepoDir, args[0])
				profile2Name = args[1]
				profile2Path = fmt.Sprintf("%s/profiles/%s/CLAUDE.md", RepoDir, args[1])
			}

			// Check if profiles exist
			if !mgr.ProfileExists(profile1Name) {
				return fmt.Errorf("profile '%s' does not exist", profile1Name)
			}
			if !mgr.ProfileExists(profile2Name) {
				return fmt.Errorf("profile '%s' does not exist", profile2Name)
			}

			// Check if diff is available
			diffCmd := exec.Command("diff", "-u", profile1Path, profile2Path)
			diffCmd.Stdout = os.Stdout
			diffCmd.Stderr = os.Stderr

			// Run diff (exit code 1 means differences found, which is expected)
			if err := diffCmd.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					// Exit code 1 means differences found (normal for diff)
					if exitErr.ExitCode() == 1 {
						return nil
					}
				}
				return fmt.Errorf("failed to run diff: %w", err)
			}

			return nil
		},
	}
}
