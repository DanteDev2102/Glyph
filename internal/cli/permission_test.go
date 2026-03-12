package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPermissionPreservation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src")
	err = os.WriteFile(srcFile, []byte("test"), 0700)
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

	if info.Mode().Perm() != 0700 {
		t.Errorf("copyFile: Expected permission 0700, got %o", info.Mode().Perm())
	}

	// Test copyDir
	srcDir := filepath.Join(tmpDir, "srcdir")
	err = os.Mkdir(srcDir, 0750)
	if err != nil {
		t.Fatal(err)
	}
	srcFileInDir := filepath.Join(srcDir, "file")
	err = os.WriteFile(srcFileInDir, []byte("test"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dstdir")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	info, err = os.Stat(dstDir)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0750 {
		t.Errorf("copyDir: Expected dir permission 0750, got %o", info.Mode().Perm())
	}

	info, err = os.Stat(filepath.Join(dstDir, "file"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("copyDir: Expected file permission 0700, got %o", info.Mode().Perm())
	}
}
