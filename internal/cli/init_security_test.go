package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplaceInFile_Symlink(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	targetFile := filepath.Join(tmpDir, "target")
	err = os.WriteFile(targetFile, []byte("Hello {{.Name}}"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	symlinkFile := filepath.Join(tmpDir, "symlink")
	err = os.Symlink(targetFile, symlinkFile)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"Name": "World"}
	replaceInFile(symlinkFile, replacements)

	// Check if targetFile was modified
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "Hello {{.Name}}" {
		t.Errorf("Target file was modified through symlink! Content: %s", string(content))
	}
}

func TestReplaceInFile_RegularFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	regularFile := filepath.Join(tmpDir, "regular")
	err = os.WriteFile(regularFile, []byte("Hello {{.Name}}"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"Name": "World"}
	replaceInFile(regularFile, replacements)

	content, err := os.ReadFile(regularFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "Hello World" {
		t.Errorf("Regular file was not modified! Content: %s", string(content))
	}
}

func TestDstPathValidation(t *testing.T) {
	cases := []struct {
		input    string
		expected bool // true if valid, false if invalid
	}{
		{"myproject", true},
		{"/", false},
		{".", false},
		{"..", false},
		{"/etc", true}, // We currently allow absolute paths that are not root
		{"   ", false},
		{"a", false},
	}

	for _, c := range cases {
		dstPath := filepath.Clean(c.input)
		valid := !(dstPath == "/" || dstPath == "." || dstPath == ".." || len(strings.TrimSpace(dstPath)) <= 1)
		if valid != c.expected {
			t.Errorf("Validation failed for %s: expected %v, got %v", c.input, c.expected, valid)
		}
	}
}

func TestCopyFile_Symlink(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	targetFile := filepath.Join(tmpDir, "target")
	err = os.WriteFile(targetFile, []byte("Original Content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "symlink_to_target")
	err = os.Symlink(targetFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}

	srcFile := filepath.Join(tmpDir, "source")
	err = os.WriteFile(srcFile, []byte("Malicious Content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(srcFile, dstFile)
	if err == nil {
		t.Error("Expected error when copying to a symlink destination, but got nil")
	} else if !strings.Contains(err.Error(), "destination is a symlink") {
		t.Errorf("Expected error containing 'destination is a symlink', but got: %v", err)
	}

	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) == "Malicious Content" {
		t.Errorf("Security Vulnerability: copyFile followed symlink and modified target file!")
	}
}

func TestCopyDir_Symlink(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source directory with a file
	srcDir := filepath.Join(tmpDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("Source Content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create destination directory with a symlink pointing elsewhere
	dstDir := filepath.Join(tmpDir, "dst")
	err = os.MkdirAll(dstDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	targetFile := filepath.Join(tmpDir, "secret")
	err = os.WriteFile(targetFile, []byte("Secret Content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create a symlink in dst that points to secret
	err = os.Symlink(targetFile, filepath.Join(dstDir, "file.txt"))
	if err != nil {
		t.Fatal(err)
	}

	// Try to copy src to dst
	err = copyDir(srcDir, dstDir)
	if err == nil {
		t.Error("Expected error when copying to a directory containing a symlink, but got nil")
	} else if !strings.Contains(err.Error(), "destination is a symlink") {
		t.Errorf("Expected error containing 'destination is a symlink', but got: %v", err)
	}

	// Check if secret file was NOT modified
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) == "Source Content" {
		t.Errorf("Security Vulnerability: copyDir followed symlink and modified target file!")
	}
}
