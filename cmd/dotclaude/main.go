// Package main provides the entry point for the dotclaude CLI.
package main

import (
	"os"

	"github.com/blackwell-systems/dotclaude/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
