package cli

import (
	"os"
)

var (
	// ColorEnabled indicates if terminal colors are supported
	ColorEnabled = true

	// Color codes (empty if colors disabled)
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[1;32m"
	ColorYellow = "\033[1;33m"
	ColorRed    = "\033[1;31m"
	ColorCyan   = "\033[1;36m"
	ColorBold   = "\033[1m"
)

func init() {
	initTerminalCommon()
	initTerminalPlatform()
}

// initTerminalCommon handles cross-platform terminal checks
func initTerminalCommon() {
	// Check if output is a terminal
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		// Not a terminal (piped), disable colors
		disableColors()
		return
	}

	// Check NO_COLOR environment variable (https://no-color.org/)
	if os.Getenv("NO_COLOR") != "" {
		disableColors()
		return
	}

	// Check TERM environment variable
	if os.Getenv("TERM") == "dumb" {
		disableColors()
		return
	}
}

// disableColors turns off all color output
func disableColors() {
	ColorEnabled = false
	ColorReset = ""
	ColorGreen = ""
	ColorYellow = ""
	ColorRed = ""
	ColorCyan = ""
	ColorBold = ""
}

// Green returns text wrapped in green color codes
func Green(s string) string {
	return ColorGreen + s + ColorReset
}

// Yellow returns text wrapped in yellow color codes
func Yellow(s string) string {
	return ColorYellow + s + ColorReset
}

// Red returns text wrapped in red color codes
func Red(s string) string {
	return ColorRed + s + ColorReset
}

// Cyan returns text wrapped in cyan color codes
func Cyan(s string) string {
	return ColorCyan + s + ColorReset
}

// Bold returns text wrapped in bold codes
func Bold(s string) string {
	return ColorBold + s + ColorReset
}
