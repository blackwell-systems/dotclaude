package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	var defaultBranch string

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync feature branch with main",
		Long:  "Keeps feature branches up-to-date by rebasing or merging with the default branch.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if we're in a git repository
			if err := exec.Command("git", "rev-parse", "--git-dir").Run(); err != nil {
				return fmt.Errorf("not in a git repository")
			}

			// Get current branch
			currentCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
			currentOutput, err := currentCmd.Output()
			if err != nil {
				return fmt.Errorf("failed to get current branch: %w", err)
			}
			currentBranch := strings.TrimSpace(string(currentOutput))

			fmt.Println("=== Feature Branch Sync Tool ===")
			fmt.Println()

			// If on main/master, list feature branches that need syncing
			if currentBranch == "main" || currentBranch == "master" {
				fmt.Printf("Currently on %s\n", currentBranch)
				fmt.Println()
				fmt.Println("Feature branches that are behind:")
				fmt.Println()

				// Fetch latest
				fetchCmd := exec.Command("git", "fetch", "origin", "--quiet")
				if err := fetchCmd.Run(); err != nil {
					if Verbose {
						fmt.Printf("Warning: git fetch failed: %v\n", err)
					}
				}

				// Find branches behind main
				branchCmd := exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/heads/")
				branchOutput, err := branchCmd.Output()
				if err != nil {
					return fmt.Errorf("failed to list branches: %w", err)
				}

				branches := strings.Split(strings.TrimSpace(string(branchOutput)), "\n")
				foundBehind := false

				for _, branch := range branches {
					branch = strings.TrimSpace(branch)
					if branch == "" || branch == "main" || branch == "master" {
						continue
					}

					behind, err := getCommitCount(branch, defaultBranch)
					if err != nil || behind == 0 {
						continue
					}

					fmt.Printf("  - %s (behind by %d commits)\n", branch, behind)
					foundBehind = true
				}

				if !foundBehind {
					fmt.Println("  All feature branches are up-to-date!")
				} else {
					fmt.Println()
					fmt.Println("To sync a branch, run:")
					fmt.Println("  git checkout <branch-name>")
					fmt.Println("  dotclaude sync")
				}

				fmt.Println()
				return nil
			}

			// We're on a feature branch
			fmt.Printf("Current branch: %s\n", currentBranch)
			fmt.Println()

			// Check if branch is behind main
			behind, err := getCommitCount(currentBranch, defaultBranch)
			if err != nil {
				return fmt.Errorf("failed to check branch status: %w", err)
			}

			ahead, err := getCommitCount(defaultBranch, currentBranch)
			if err != nil {
				return fmt.Errorf("failed to check branch status: %w", err)
			}

			fmt.Printf("Status: %d commits ahead, %d commits behind %s\n", ahead, behind, defaultBranch)

			if behind == 0 {
				fmt.Println("✓ Branch is already up-to-date with", defaultBranch)
				return nil
			}

			// Check for uncommitted changes
			statusCmd := exec.Command("git", "diff-index", "--quiet", "HEAD", "--")
			if err := statusCmd.Run(); err != nil {
				fmt.Println("Error: You have uncommitted changes. Commit or stash them first.")
				fmt.Println()
				statusOutput := exec.Command("git", "status", "--short")
				statusOutput.Stdout = os.Stdout
				statusOutput.Run()
				return fmt.Errorf("uncommitted changes present")
			}

			fmt.Println()
			fmt.Printf("Branch is %d commits behind %s\n", behind, defaultBranch)
			fmt.Println()
			fmt.Println("Choose sync method:")
			fmt.Println("  1) Rebase (cleaner history, requires force push)")
			fmt.Println("  2) Merge (preserves history, no force push needed)")
			fmt.Println("  3) Cancel")
			fmt.Println()

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Selection (1/2/3): ")
			choice, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			choice = strings.TrimSpace(choice)

			switch choice {
			case "1":
				return syncRebase(currentBranch, defaultBranch, reader)
			case "2":
				return syncMerge(currentBranch, defaultBranch, reader)
			case "3":
				fmt.Println("Cancelled")
				return nil
			default:
				fmt.Println("Invalid selection")
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&defaultBranch, "base", "b", "main", "base branch to sync with")

	return cmd
}

func syncRebase(currentBranch, defaultBranch string, reader *bufio.Reader) error {
	fmt.Printf("Rebasing %s onto %s...\n", currentBranch, defaultBranch)
	fmt.Println()

	// Fetch latest
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Stdout = os.Stdout
	fetchCmd.Stderr = os.Stderr
	if err := fetchCmd.Run(); err != nil {
		return fmt.Errorf("git fetch failed: %w", err)
	}

	// Rebase
	rebaseCmd := exec.Command("git", "rebase", "origin/"+defaultBranch)
	rebaseCmd.Stdout = os.Stdout
	rebaseCmd.Stderr = os.Stderr

	if err := rebaseCmd.Run(); err != nil {
		fmt.Println()
		fmt.Println("Rebase failed. Resolve conflicts and run:")
		fmt.Println("  git rebase --continue")
		fmt.Println("  git push --force-with-lease")
		return fmt.Errorf("rebase failed")
	}

	fmt.Println()
	fmt.Println("✓ Rebase successful")
	fmt.Println()
	fmt.Println("To push changes:")
	fmt.Println("  git push --force-with-lease")
	fmt.Println()
	fmt.Print("Push now? (y/N): ")

	confirm, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm == "y" || confirm == "yes" {
		pushCmd := exec.Command("git", "push", "--force-with-lease")
		pushCmd.Stdout = os.Stdout
		pushCmd.Stderr = os.Stderr
		if err := pushCmd.Run(); err != nil {
			return fmt.Errorf("git push failed: %w", err)
		}
		fmt.Println("✓ Branch synced and pushed")
	}

	return nil
}

func syncMerge(currentBranch, defaultBranch string, reader *bufio.Reader) error {
	fmt.Printf("Merging %s into %s...\n", defaultBranch, currentBranch)
	fmt.Println()

	// Fetch latest
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Stdout = os.Stdout
	fetchCmd.Stderr = os.Stderr
	if err := fetchCmd.Run(); err != nil {
		return fmt.Errorf("git fetch failed: %w", err)
	}

	// Merge
	mergeMsg := fmt.Sprintf("Merge %s into %s", defaultBranch, currentBranch)
	mergeCmd := exec.Command("git", "merge", "origin/"+defaultBranch, "-m", mergeMsg)
	mergeCmd.Stdout = os.Stdout
	mergeCmd.Stderr = os.Stderr

	if err := mergeCmd.Run(); err != nil {
		fmt.Println()
		fmt.Println("Merge failed. Resolve conflicts and run:")
		fmt.Println("  git merge --continue")
		fmt.Println("  git push")
		return fmt.Errorf("merge failed")
	}

	fmt.Println()
	fmt.Println("✓ Merge successful")
	fmt.Println()
	fmt.Println("To push changes:")
	fmt.Println("  git push")
	fmt.Println()
	fmt.Print("Push now? (y/N): ")

	confirm, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm == "y" || confirm == "yes" {
		pushCmd := exec.Command("git", "push")
		pushCmd.Stdout = os.Stdout
		pushCmd.Stderr = os.Stderr
		if err := pushCmd.Run(); err != nil {
			return fmt.Errorf("git push failed: %w", err)
		}
		fmt.Println("✓ Branch synced and pushed")
	}

	return nil
}

// Reuse getCommitCount from check_branches.go
func getCommitCountForSync(from, to string) (int, error) {
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
