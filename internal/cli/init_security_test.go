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
	tmpDir, err := os.MkdirTemp("", "glyph-test-copyfile-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src")
	os.WriteFile(srcFile, []byte("New Content"), 0644)

	targetFile := filepath.Join(tmpDir, "target")
	os.WriteFile(targetFile, []byte("Original Content"), 0644)

	symlinkFile := filepath.Join(tmpDir, "symlink")
	os.Symlink(targetFile, symlinkFile)

	err = copyFile(srcFile, symlinkFile)
	if err == nil {
		t.Errorf("copyFile should have failed when destination is a symlink")
	}

	content, _ := os.ReadFile(targetFile)
	if string(content) != "Original Content" {
		t.Errorf("Target file was modified through symlink!")
	}
}

func TestCopyDir_Symlink(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-copydir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	os.Mkdir(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "file1"), []byte("New Content"), 0644)

	dstDir := filepath.Join(tmpDir, "dst")
	os.Mkdir(dstDir, 0755)

	targetFile := filepath.Join(tmpDir, "target")
	os.WriteFile(targetFile, []byte("Original Content"), 0644)

	symlinkFile := filepath.Join(dstDir, "file1")
	os.Symlink(targetFile, symlinkFile)

	err = copyDir(srcDir, dstDir)
	if err == nil {
		t.Errorf("copyDir should have failed when destination contains a symlink")
	}

	content, _ := os.ReadFile(targetFile)
	if string(content) != "Original Content" {
		t.Errorf("Target file was modified through symlink!")
	}
}
