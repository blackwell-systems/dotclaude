// Package cli implements the dotclaude command-line interface.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	// Version is the current dotclaude version
	Version = "1.0.0-alpha.1"
)

var (
	// RepoDir is the dotclaude repository directory
	RepoDir string
	// ClaudeDir is the Claude Code configuration directory
	ClaudeDir string
	// ProfilesDir is the directory containing profiles
	ProfilesDir string
	// Verbose enables debug output
	Verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotclaude",
	Short: "The definitive profile management system for Claude Code",
	Long: `dotclaude - Profile management for Claude Code

Manage your Claude Code configuration as layered, version-controlled profiles.
Switch between work contexts (OSS, client, employer) with one command.

Examples:
  dotclaude create my-project     Create a new profile
  dotclaude activate my-project   Activate a profile
  dotclaude list                  List all profiles
  dotclaude show                  Show active profile`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	// Set defaults
	if RepoDir == "" {
		if dir := os.Getenv("DOTCLAUDE_REPO_DIR"); dir != "" {
			RepoDir = dir
		} else {
			home, _ := os.UserHomeDir()
			RepoDir = home + "/code/dotclaude"
		}
	}

	if ClaudeDir == "" {
		if dir := os.Getenv("CLAUDE_DIR"); dir != "" {
			ClaudeDir = dir
		} else {
			home, _ := os.UserHomeDir()
			ClaudeDir = home + "/.claude"
		}
	}

	ProfilesDir = RepoDir + "/profiles"

	// Add subcommands
	rootCmd.AddCommand(
		newVersionCmd(),
		newListCmd(),
		newShowCmd(),
		newCreateCmd(),
		newDeleteCmd(),
		newEditCmd(),
		newActivateCmd(),
		newRestoreCmd(),
		newCheckBranchesCmd(),
	)
}

func initConfig() {
	// Configuration initialization if needed
	if Verbose {
		fmt.Fprintf(os.Stderr, "RepoDir: %s\n", RepoDir)
		fmt.Fprintf(os.Stderr, "ClaudeDir: %s\n", ClaudeDir)
		fmt.Fprintf(os.Stderr, "ProfilesDir: %s\n", ProfilesDir)
	}
}
