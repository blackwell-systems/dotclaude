package cli

import (
	"os"
	"testing"
)

func TestColorFunctions(t *testing.T) {
	// Save original values
	origEnabled := ColorEnabled
	origGreen := ColorGreen
	origReset := ColorReset
	origYellow := ColorYellow
	origRed := ColorRed
	origCyan := ColorCyan
	origBold := ColorBold

	defer func() {
		ColorEnabled = origEnabled
		ColorGreen = origGreen
		ColorReset = origReset
		ColorYellow = origYellow
		ColorRed = origRed
		ColorCyan = origCyan
		ColorBold = origBold
	}()

	t.Run("Green with colors enabled", func(t *testing.T) {
		ColorEnabled = true
		ColorGreen = "\033[1;32m"
		ColorReset = "\033[0m"

		result := Green("test")
		expected := "\033[1;32mtest\033[0m"
		if result != expected {
			t.Errorf("Green() = %q, want %q", result, expected)
		}
	})

	t.Run("Green with colors disabled", func(t *testing.T) {
		ColorEnabled = false
		ColorGreen = ""
		ColorReset = ""

		result := Green("test")
		if result != "test" {
			t.Errorf("Green() with disabled colors = %q, want %q", result, "test")
		}
	})

	t.Run("Yellow with colors enabled", func(t *testing.T) {
		ColorEnabled = true
		ColorYellow = "\033[1;33m"
		ColorReset = "\033[0m"

		result := Yellow("test")
		expected := "\033[1;33mtest\033[0m"
		if result != expected {
			t.Errorf("Yellow() = %q, want %q", result, expected)
		}
	})

	t.Run("Red with colors enabled", func(t *testing.T) {
		ColorEnabled = true
		ColorRed = "\033[1;31m"
		ColorReset = "\033[0m"

		result := Red("test")
		expected := "\033[1;31mtest\033[0m"
		if result != expected {
			t.Errorf("Red() = %q, want %q", result, expected)
		}
	})

	t.Run("Cyan with colors enabled", func(t *testing.T) {
		ColorEnabled = true
		ColorCyan = "\033[1;36m"
		ColorReset = "\033[0m"

		result := Cyan("test")
		expected := "\033[1;36mtest\033[0m"
		if result != expected {
			t.Errorf("Cyan() = %q, want %q", result, expected)
		}
	})

	t.Run("Bold with colors enabled", func(t *testing.T) {
		ColorEnabled = true
		ColorBold = "\033[1m"
		ColorReset = "\033[0m"

		result := Bold("test")
		expected := "\033[1mtest\033[0m"
		if result != expected {
			t.Errorf("Bold() = %q, want %q", result, expected)
		}
	})
}

func TestDisableColors(t *testing.T) {
	// Save original values
	origEnabled := ColorEnabled
	origGreen := ColorGreen
	origReset := ColorReset
	origYellow := ColorYellow
	origRed := ColorRed
	origCyan := ColorCyan
	origBold := ColorBold

	defer func() {
		ColorEnabled = origEnabled
		ColorGreen = origGreen
		ColorReset = origReset
		ColorYellow = origYellow
		ColorRed = origRed
		ColorCyan = origCyan
		ColorBold = origBold
	}()

	// Set colors to non-empty
	ColorEnabled = true
	ColorGreen = "\033[1;32m"
	ColorReset = "\033[0m"
	ColorYellow = "\033[1;33m"
	ColorRed = "\033[1;31m"
	ColorCyan = "\033[1;36m"
	ColorBold = "\033[1m"

	disableColors()

	if ColorEnabled {
		t.Error("Expected ColorEnabled to be false after disableColors()")
	}
	if ColorGreen != "" {
		t.Errorf("Expected ColorGreen to be empty, got %q", ColorGreen)
	}
	if ColorReset != "" {
		t.Errorf("Expected ColorReset to be empty, got %q", ColorReset)
	}
	if ColorYellow != "" {
		t.Errorf("Expected ColorYellow to be empty, got %q", ColorYellow)
	}
	if ColorRed != "" {
		t.Errorf("Expected ColorRed to be empty, got %q", ColorRed)
	}
	if ColorCyan != "" {
		t.Errorf("Expected ColorCyan to be empty, got %q", ColorCyan)
	}
	if ColorBold != "" {
		t.Errorf("Expected ColorBold to be empty, got %q", ColorBold)
	}
}

func TestInitTerminalCommon(t *testing.T) {
	// Save original values
	origEnabled := ColorEnabled
	origGreen := ColorGreen
	origReset := ColorReset

	defer func() {
		ColorEnabled = origEnabled
		ColorGreen = origGreen
		ColorReset = origReset
	}()

	t.Run("NO_COLOR environment variable", func(t *testing.T) {
		ColorEnabled = true
		ColorGreen = "\033[1;32m"
		ColorReset = "\033[0m"

		// Set NO_COLOR
		oldVal := os.Getenv("NO_COLOR")
		os.Setenv("NO_COLOR", "1")
		defer func() {
			if oldVal == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", oldVal)
			}
		}()

		initTerminalCommon()

		if ColorEnabled {
			t.Error("Expected colors to be disabled with NO_COLOR set")
		}
	})

	t.Run("TERM=dumb environment variable", func(t *testing.T) {
		ColorEnabled = true
		ColorGreen = "\033[1;32m"
		ColorReset = "\033[0m"

		// Remove NO_COLOR and set TERM=dumb
		os.Unsetenv("NO_COLOR")
		oldTerm := os.Getenv("TERM")
		os.Setenv("TERM", "dumb")
		defer func() {
			if oldTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", oldTerm)
			}
		}()

		initTerminalCommon()

		if ColorEnabled {
			t.Error("Expected colors to be disabled with TERM=dumb")
		}
	})
}
