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

	file600 := filepath.Join(tmpDir, "file600")
	err = os.WriteFile(file600, []byte("Hello {{.Name}}"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"Name": "World"}
	replaceInFile(file600, replacements)

	info, err := os.Stat(file600)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected permissions 0600, got %v", info.Mode().Perm())
	}
}

func TestCopyFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src")
	err = os.WriteFile(srcFile, []byte("content"), 0600)
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

	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected permissions 0600, got %v", info.Mode().Perm())
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

	srcFile := filepath.Join(srcDir, "file")
	err = os.WriteFile(srcFile, []byte("content"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Join(dstDir, "file"))
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected permissions 0600, got %v", info.Mode().Perm())
	}
}
