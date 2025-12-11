package profile

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestRepo creates a temp directory with the expected repo structure
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
	if err != nil {
		t.Fatal(err)
	}

	// Create base directory
	baseDir := filepath.Join(tmpDir, "base")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(baseDir, "CLAUDE.md"), []byte("# Base Config\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(baseDir, "settings.json"), []byte(`{"key": "value"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create sample-profile template
	templateDir := filepath.Join(tmpDir, "examples", "sample-profile")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "CLAUDE.md"), []byte("# Sample Profile\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .claude directory
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestCreate(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	claudeDir := filepath.Join(tmpDir, ".claude")
	mgr := NewManager(tmpDir, claudeDir)

	t.Run("create new profile", func(t *testing.T) {
		err := mgr.Create("my-new-profile")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Verify profile directory exists
		profileDir := filepath.Join(tmpDir, "profiles", "my-new-profile")
		if _, err := os.Stat(profileDir); os.IsNotExist(err) {
			t.Error("Profile directory was not created")
		}

		// Verify CLAUDE.md was copied
		claudeMd := filepath.Join(profileDir, "CLAUDE.md")
		if _, err := os.Stat(claudeMd); os.IsNotExist(err) {
			t.Error("CLAUDE.md was not copied")
		}
	})

	t.Run("create duplicate profile", func(t *testing.T) {
		err := mgr.Create("my-new-profile")
		if err == nil {
			t.Error("Create() should error for duplicate profile")
		}
	})

	t.Run("create with invalid name", func(t *testing.T) {
		err := mgr.Create("invalid/name")
		if err == nil {
			t.Error("Create() should error for invalid name")
		}
	})

	t.Run("create with empty name", func(t *testing.T) {
		err := mgr.Create("")
		if err == nil {
			t.Error("Create() should error for empty name")
		}
	})
}

func TestCopyDir(t *testing.T) {
	// Create source directory with structure
	srcDir, err := os.MkdirTemp("", "dotclaude-src-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	// Create files and subdirectories
	if err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "file2.txt"), []byte("content2"), 0600); err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(srcDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Copy to destination
	dstDir, err := os.MkdirTemp("", "dotclaude-dst-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	targetDir := filepath.Join(dstDir, "copied")
	if err := copyDir(srcDir, targetDir); err != nil {
		t.Fatalf("copyDir() error = %v", err)
	}

	// Verify files were copied
	t.Run("files copied", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(targetDir, "file1.txt"))
		if err != nil {
			t.Fatalf("Failed to read copied file: %v", err)
		}
		if string(content) != "content1" {
			t.Errorf("Content = %q, want %q", string(content), "content1")
		}
	})

	t.Run("subdirectory copied", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(targetDir, "subdir", "nested.txt"))
		if err != nil {
			t.Fatalf("Failed to read nested file: %v", err)
		}
		if string(content) != "nested content" {
			t.Errorf("Content = %q, want %q", string(content), "nested content")
		}
	})
}

func TestCopyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "dest.txt")

	// Create source file
	content := "test content"
	if err := os.WriteFile(srcPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("copy file", func(t *testing.T) {
		if err := copyFile(srcPath, dstPath); err != nil {
			t.Fatalf("copyFile() error = %v", err)
		}

		copied, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("Failed to read copied file: %v", err)
		}

		if string(copied) != content {
			t.Errorf("Content = %q, want %q", string(copied), content)
		}
	})

	t.Run("copy non-existent file", func(t *testing.T) {
		err := copyFile("/non/existent/file", dstPath)
		if err == nil {
			t.Error("copyFile() should error for non-existent source")
		}
	})
}

func TestCreate_ErrorPaths(t *testing.T) {
	t.Run("missing template directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create minimal structure WITHOUT template
		claudeDir := filepath.Join(tmpDir, ".claude")
		if err := os.MkdirAll(claudeDir, 0755); err != nil {
			t.Fatal(err)
		}

		mgr := NewManager(tmpDir, claudeDir)
		err = mgr.Create("test-profile")
		if err == nil {
			t.Error("Create() should error when template directory is missing")
		}
	})

	t.Run("create ensures profiles directory exists", func(t *testing.T) {
		tmpDir, cleanup := setupTestRepo(t)
		defer cleanup()

		claudeDir := filepath.Join(tmpDir, ".claude")
		mgr := NewManager(tmpDir, claudeDir)

		// Remove profiles directory if it exists
		profilesDir := filepath.Join(tmpDir, "profiles")
		os.RemoveAll(profilesDir)

		err := mgr.Create("new-profile")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Verify profiles directory was created
		if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
			t.Error("Create() should have created profiles directory")
		}
	})
}

func TestCopyDir_ErrorPaths(t *testing.T) {
	t.Run("copy non-existent source", func(t *testing.T) {
		dstDir, err := os.MkdirTemp("", "dotclaude-dst-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dstDir)

		err = copyDir("/non/existent/source", filepath.Join(dstDir, "target"))
		if err == nil {
			t.Error("copyDir() should error for non-existent source")
		}
	})
}

func TestCopyFile_ErrorPaths(t *testing.T) {
	t.Run("copy to invalid destination", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "dotclaude-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		srcPath := filepath.Join(tmpDir, "source.txt")
		if err := os.WriteFile(srcPath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		// Try to copy to a path where parent doesn't exist
		err = copyFile(srcPath, "/non/existent/path/dest.txt")
		if err == nil {
			t.Error("copyFile() should error for invalid destination path")
		}
	})
}
