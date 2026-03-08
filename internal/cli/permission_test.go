package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPermissionPreservation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-perm-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")

	// Create source directory with specific permissions
	err = os.MkdirAll(srcDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	// Create a file with specific permissions
	srcFile := filepath.Join(srcDir, "testfile")
	err = os.WriteFile(srcFile, []byte("test content"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	// Test copyFile
	dstFile := filepath.Join(tmpDir, "dstfile")
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	info, err := os.Stat(dstFile)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("copyFile: expected permissions 0600, got %o", info.Mode().Perm())
	}

	// Test copyDir
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatalf("copyDir failed: %v", err)
	}

	// Check dstDir permissions
	info, err = os.Stat(dstDir)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("copyDir (directory): expected permissions 0700, got %o", info.Mode().Perm())
	}

	// Check dstDir/testfile permissions
	info, err = os.Stat(filepath.Join(dstDir, "testfile"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("copyDir (file): expected permissions 0600, got %o", info.Mode().Perm())
	}
}
