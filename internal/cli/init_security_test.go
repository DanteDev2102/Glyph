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

func TestCopyDir_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	err = os.Mkdir(srcDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	execFile := filepath.Join(srcDir, "exec.sh")
	err = os.WriteFile(execFile, []byte("#!/bin/sh"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	copiedFile := filepath.Join(dstDir, "exec.sh")
	info, err := os.Stat(copiedFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0755 {
		t.Errorf("Expected permissions 0755, got %o", info.Mode().Perm())
	}
}

func TestReplaceInFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	execFile := filepath.Join(tmpDir, "exec.sh")
	err = os.WriteFile(execFile, []byte("#!/bin/sh\necho {{.ProjectName}}"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"ProjectName": "TestApp"}
	replaceInFile(execFile, replacements)

	info, err := os.Stat(execFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0700 {
		t.Errorf("Expected permissions 0700, got %o", info.Mode().Perm())
	}

	content, err := os.ReadFile(execFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "TestApp") {
		t.Errorf("Replacement failed, content: %s", string(content))
	}
}
