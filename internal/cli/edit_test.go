package cli

import (
	"os"
	"runtime"
	"testing"
)

func TestGetEditor(t *testing.T) {
	// Save original environment variables
	origEditor := os.Getenv("EDITOR")
	origVisual := os.Getenv("VISUAL")

	defer func() {
		if origEditor != "" {
			os.Setenv("EDITOR", origEditor)
		} else {
			os.Unsetenv("EDITOR")
		}
		if origVisual != "" {
			os.Setenv("VISUAL", origVisual)
		} else {
			os.Unsetenv("VISUAL")
		}
	}()

	t.Run("uses EDITOR env var", func(t *testing.T) {
		os.Setenv("EDITOR", "myeditor")
		os.Unsetenv("VISUAL")

		editor := getEditor()
		if editor != "myeditor" {
			t.Errorf("getEditor() = %q, want %q", editor, "myeditor")
		}
	})

	t.Run("uses VISUAL env var when EDITOR not set", func(t *testing.T) {
		os.Unsetenv("EDITOR")
		os.Setenv("VISUAL", "myvisual")

		editor := getEditor()
		if editor != "myvisual" {
			t.Errorf("getEditor() = %q, want %q", editor, "myvisual")
		}
	})

	t.Run("EDITOR takes precedence over VISUAL", func(t *testing.T) {
		os.Setenv("EDITOR", "myeditor")
		os.Setenv("VISUAL", "myvisual")

		editor := getEditor()
		if editor != "myeditor" {
			t.Errorf("getEditor() = %q, want %q", editor, "myeditor")
		}
	})

	t.Run("returns default when no env vars set", func(t *testing.T) {
		os.Unsetenv("EDITOR")
		os.Unsetenv("VISUAL")

		editor := getEditor()
		// Should return some default (platform-specific)
		if editor == "" {
			t.Error("getEditor() returned empty string, expected a default")
		}

		// On Windows, should be notepad if nothing else available
		if runtime.GOOS == "windows" {
			// Accept any non-empty result on Windows
			if editor == "" {
				t.Error("Expected non-empty editor on Windows")
			}
		}
	})

	t.Run("handles editor with arguments", func(t *testing.T) {
		os.Setenv("EDITOR", "code --wait")
		os.Unsetenv("VISUAL")

		editor := getEditor()
		if editor != "code --wait" {
			t.Errorf("getEditor() = %q, want %q", editor, "code --wait")
		}
	})
}
