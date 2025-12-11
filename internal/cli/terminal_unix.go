//go:build !windows

package cli

// initTerminalPlatform handles Unix-specific terminal setup
// Unix terminals generally support ANSI colors natively
func initTerminalPlatform() {
	// No special handling needed for Unix
	// ANSI colors work by default on most terminals
}
