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

func newRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore from backup interactively",
		Long:  "Restore CLAUDE.md or settings.json from a backup file.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := profile.NewManager(RepoDir, ClaudeDir)

			// Header
			fmt.Println()
			fmt.Println("╭─────────────────────────────────────────────────────────────╮")
			fmt.Println("│  Backup Restoration                                         │")
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()

			// List backups
			backups, err := mgr.ListBackups()
			if err != nil {
				return err
			}

			if len(backups) == 0 {
				fmt.Println("  No backups found")
				fmt.Println()
				fmt.Println("  Backups are created automatically when switching profiles.")
				fmt.Println()
				fmt.Println("  Tip: Use 'dotclaude activate' to switch profiles")
				fmt.Println()
				return nil
			}

			fmt.Println("  Available backups:")
			fmt.Println()

			// Group backups by type
			claudeBackups := []*profile.Backup{}
			settingsBackups := []*profile.Backup{}

			for _, backup := range backups {
				if backup.Type == "CLAUDE.md" {
					claudeBackups = append(claudeBackups, backup)
				} else {
					settingsBackups = append(settingsBackups, backup)
				}
			}

			allBackups := []*profile.Backup{}
			index := 1

			// List CLAUDE.md backups
			if len(claudeBackups) > 0 {
				fmt.Println("  CLAUDE.md backups:")
				for _, backup := range claudeBackups {
					sizeKB := backup.Size / 1024
					fmt.Printf("    [%d] %s (%dK)\n", index, backup.Timestamp, sizeKB)
					allBackups = append(allBackups, backup)
					index++
				}
				fmt.Println()
			}

			// List settings.json backups
			if len(settingsBackups) > 0 {
				fmt.Println("  settings.json backups:")
				for _, backup := range settingsBackups {
					sizeKB := backup.Size / 1024
					fmt.Printf("    [%d] %s (%dK)\n", index, backup.Timestamp, sizeKB)
					allBackups = append(allBackups, backup)
					index++
				}
				fmt.Println()
			}

			// Prompt for selection
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("  Select backup to restore (or 'q' to quit): ")
			choice, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			choice = strings.TrimSpace(choice)

			if choice == "q" || choice == "Q" {
				fmt.Println()
				fmt.Println("  Cancelled")
				fmt.Println()
				return nil
			}

			// Parse selection
			selection, err := strconv.Atoi(choice)
			if err != nil || selection < 1 || selection > len(allBackups) {
				return fmt.Errorf("invalid selection")
			}

			selectedBackup := allBackups[selection-1]

			// Determine target file
			var targetFile string
			if selectedBackup.Type == "CLAUDE.md" {
				targetFile = ClaudeDir + "/CLAUDE.md"
			} else {
				targetFile = ClaudeDir + "/settings.json"
			}

			// Confirm overwrite
			fmt.Println()
			fmt.Println("  ⚠  This will overwrite:")
			fmt.Printf("    %s\n", targetFile)
			fmt.Println()
			fmt.Print("  Continue? (y/N): ")

			confirm, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			confirm = strings.TrimSpace(strings.ToLower(confirm))

			if confirm != "y" && confirm != "yes" {
				fmt.Println()
				fmt.Println("  Cancelled")
				fmt.Println()
				return nil
			}

			// Restore the backup
			if err := mgr.Restore(selectedBackup.Path); err != nil {
				return err
			}

			// Success message
			fmt.Println()
			fmt.Println("  [BACKUP] Current file backed up")
			fmt.Printf("  [RESTORE] Restored from: %s\n", selectedBackup.Filename)

			if selectedBackup.Type == "CLAUDE.md" {
				// Try to show the profile name if we updated it
				activeName := mgr.GetActiveProfileName()
				if activeName != "" {
					fmt.Printf("  [UPDATE] Active profile: %s\n", activeName)
				}
			}

			fmt.Println()
			fmt.Println("╭─────────────────────────────────────────────────────────────╮")
			fmt.Println("│  ✓ Backup Restored                                          │")
			fmt.Println("╰─────────────────────────────────────────────────────────────╯")
			fmt.Println()

			return nil
		},
	}

	return cmd
}
