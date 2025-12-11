package cli

import (
	"fmt"

	"github.com/blackwell-systems/dotclaude/internal/hooks"
	"github.com/spf13/cobra"
)

// newHookCmd returns the hook parent command
func newHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Manage and run hooks",
		Long: `Hook management for Claude Code integration.

Hooks are scripts that run at specific points in the Claude Code lifecycle.
They enable cross-platform automation and customization.

Built-in hooks provide core functionality (session info, profile mismatch detection).
Custom hooks can be added to the hooks directory (<claude-dir>/hooks/<hook-type>/).

Hook types:
  session-start    Runs when a Claude Code session starts
  post-tool-bash   Runs after Bash tool execution
  post-tool-edit   Runs after Edit tool execution
  pre-tool-bash    Runs before Bash tool execution
  pre-tool-edit    Runs before Edit tool execution

Hooks are executed in order by numeric prefix (00-first runs before 50-second).

Supported hook formats:
  Unix:    .sh, .bash (requires bash)
  Windows: .ps1, .cmd, .bat, .exe
  Any:     executable files without extension`,
	}

	cmd.AddCommand(
		newHookRunCmd(),
		newHookListCmd(),
		newHookInitCmd(),
	)

	return cmd
}

// newHookRunCmd returns the hook run command
func newHookRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <hook-type>",
		Short: "Run all hooks of a given type",
		Long: `Execute all hooks (built-in and custom) for the specified hook type.

This command is typically called from Claude Code's settings.json hooks configuration.
Hooks are run in order by numeric prefix.

Examples:
  dotclaude hook run session-start     Run session start hooks
  dotclaude hook run post-tool-bash    Run post-bash-tool hooks`,
		Args: cobra.ExactArgs(1),
		ValidArgs: []string{
			string(hooks.HookSessionStart),
			string(hooks.HookPostToolBash),
			string(hooks.HookPostToolEdit),
			string(hooks.HookPreToolBash),
			string(hooks.HookPreToolEdit),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			hookType := hooks.HookType(args[0])

			// Validate hook type
			validTypes := hooks.GetHookTypes()
			valid := false
			for _, t := range validTypes {
				if t == hookType {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("unknown hook type: %s\nValid types: session-start, post-tool-bash, post-tool-edit, pre-tool-bash, pre-tool-edit", args[0])
			}

			runner := hooks.NewRunner(ClaudeDir, RepoDir)
			return runner.Run(hookType)
		},
	}

	return cmd
}

// newHookListCmd returns the hook list command
func newHookListCmd() *cobra.Command {
	var hookType string

	cmd := &cobra.Command{
		Use:   "list [hook-type]",
		Short: "List available hooks",
		Long: `List all available hooks, optionally filtered by type.

Shows both built-in hooks and custom hooks with their priority and status.

Examples:
  dotclaude hook list                  List all hooks
  dotclaude hook list session-start    List only session-start hooks`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := hooks.NewRunner(ClaudeDir, RepoDir)

			var typesToList []hooks.HookType
			if len(args) > 0 {
				typesToList = []hooks.HookType{hooks.HookType(args[0])}
			} else {
				typesToList = hooks.GetHookTypes()
			}

			for _, ht := range typesToList {
				hookList := runner.List(ht)

				fmt.Printf("\n%s:\n", ht)
				if len(hookList) == 0 {
					fmt.Println("  (no hooks)")
					continue
				}

				for _, h := range hookList {
					status := "enabled"
					if !h.Enabled {
						status = "disabled"
					}

					if h.Type == "built-in" {
						fmt.Printf("  [%02d] %s (built-in, %s)\n", h.Priority, h.Name, status)
					} else {
						fmt.Printf("  [%02d] %s (%s)\n", h.Priority, h.Name, status)
					}
				}
			}

			fmt.Println()
			return nil
		},
	}

	cmd.Flags().StringVarP(&hookType, "type", "t", "", "Filter by hook type")

	return cmd
}

// newHookInitCmd returns the hook init command
func newHookInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize hooks directory structure",
		Long: `Create the hooks directory structure for custom hooks.

This creates directories for each hook type where custom hooks can be placed.
Custom hooks are executable files that run alongside built-in hooks.

Example hook structure after init:
  <claude-dir>/hooks/
  ├── session-start/
  ├── post-tool-bash/
  ├── post-tool-edit/
  ├── pre-tool-bash/
  └── pre-tool-edit/

To add a custom hook, place an executable script in the appropriate directory.
Use numeric prefixes to control execution order:
  Unix:    20-myhook.sh, 30-another.bash
  Windows: 20-myhook.ps1, 30-another.cmd`,
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := hooks.NewRunner(ClaudeDir, RepoDir)

			if err := runner.EnsureHooksDir(); err != nil {
				return err
			}

			fmt.Printf("Initialized hooks directory: %s/hooks/\n", ClaudeDir)
			fmt.Println("\nCreated directories:")
			for _, ht := range hooks.GetHookTypes() {
				fmt.Printf("  - %s/\n", ht)
			}

			fmt.Println("\nTo add custom hooks, place executable scripts in these directories.")
			fmt.Println("Use numeric prefixes to control order:")
			fmt.Println("  Unix:    20-myhook.sh")
			fmt.Println("  Windows: 20-myhook.ps1")

			return nil
		},
	}

	return cmd
}
