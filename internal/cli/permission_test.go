package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir_PreservesPermissions(t *testing.T) {
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

	restrictedFile := filepath.Join(srcDir, "restricted")
	expectedPerm := os.FileMode(0600)
	err = os.WriteFile(restrictedFile, []byte("restricted content"), expectedPerm)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	copiedFile := filepath.Join(dstDir, "restricted")
	info, err := os.Stat(copiedFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != expectedPerm {
		t.Errorf("Expected permission %v, got %v", expectedPerm, info.Mode().Perm())
	}
}

func TestCopyFile_PreservesPermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src")
	expectedPerm := os.FileMode(0600)
	err = os.WriteFile(srcFile, []byte("restricted content"), expectedPerm)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "dst")
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != expectedPerm {
		t.Errorf("Expected permission %v, got %v", expectedPerm, info.Mode().Perm())
	}
}

func TestReplaceInFile_PreservesPermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test")
	expectedPerm := os.FileMode(0600)
	err = os.WriteFile(testFile, []byte("Hello {{.Name}}"), expectedPerm)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"Name": "World"}
	replaceInFile(testFile, replacements)

	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != expectedPerm {
		t.Errorf("Expected permission %v, got %v", expectedPerm, info.Mode().Perm())
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "Hello World" {
		t.Errorf("Expected content 'Hello World', got '%s'", string(content))
	}
}
