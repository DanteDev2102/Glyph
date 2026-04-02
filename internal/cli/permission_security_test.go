package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplaceInFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "testfile")
	// Create file with restricted permissions
	err = os.WriteFile(testFile, []byte("Hello {{.ProjectName}}"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"ProjectName": "Test"}
	replaceInFile(testFile, replacements)

	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("Permissions not preserved! Expected 0600, got %o", info.Mode().Perm())
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "Hello Test" {
		t.Errorf("Content not replaced! Expected 'Hello Test', got '%s'", string(content))
	}
}
