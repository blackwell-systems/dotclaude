package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show dotclaude version",
		Long:  "Display the current dotclaude version and build information.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("dotclaude version %s (go)\n", Version)
			return nil
		},
	}
}
