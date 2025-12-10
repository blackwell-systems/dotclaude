package cli

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func newCheckBranchesCmd() *cobra.Command {
	var defaultBranch string

	cmd := &cobra.Command{
		Use:   "check-branches",
		Short: "Check which branches are behind main",
		Long:  "Quick check to see which feature branches are ahead of or behind the default branch.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if we're in a git repository
			if err := exec.Command("git", "rev-parse", "--git-dir").Run(); err != nil {
				return fmt.Errorf("not in a git repository")
			}

			fmt.Printf("Checking branches against %s...\n", defaultBranch)
			fmt.Println()

			// Fetch from origin
			if Verbose {
				fmt.Println("Fetching from origin...")
			}
			fetchCmd := exec.Command("git", "fetch", "origin", "--quiet")
			if err := fetchCmd.Run(); err != nil {
				// Non-fatal - continue even if fetch fails
				if Verbose {
					fmt.Printf("Warning: git fetch failed: %v\n", err)
				}
			}

			// Get all local branches
			branchCmd := exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/heads/")
			branchOutput, err := branchCmd.Output()
			if err != nil {
				return fmt.Errorf("failed to list branches: %w", err)
			}

			branches := strings.Split(strings.TrimSpace(string(branchOutput)), "\n")
			hasDivergent := false

			for _, branch := range branches {
				branch = strings.TrimSpace(branch)
				if branch == "" || branch == "main" || branch == "master" {
					continue
				}

				// Check how many commits ahead/behind
				behind, err := getCommitCount(branch, defaultBranch)
				if err != nil {
					continue // Skip branches with errors
				}

				ahead, err := getCommitCount(defaultBranch, branch)
				if err != nil {
					continue // Skip branches with errors
				}

				// Only show branches that have diverged
				if ahead > 0 || behind > 0 {
					fmt.Printf("  %-30s %d ahead, %d behind\n", branch, ahead, behind)
					hasDivergent = true
				}
			}

			if !hasDivergent {
				fmt.Println("  All branches are up to date")
			}

			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVarP(&defaultBranch, "base", "b", "main", "base branch to compare against")

	return cmd
}

// getCommitCount returns the number of commits between two branches.
// Usage: getCommitCount("branch", "main") returns commits in main not in branch
func getCommitCount(from, to string) (int, error) {
	countCmd := exec.Command("git", "rev-list", "--count", fmt.Sprintf("%s..%s", from, to))
	output, err := countCmd.Output()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}

	return count, nil
}
