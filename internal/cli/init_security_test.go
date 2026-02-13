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
