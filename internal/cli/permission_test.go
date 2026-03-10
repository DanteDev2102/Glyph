package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir_PreservePermissions(t *testing.T) {
	tmpSrc, err := os.MkdirTemp("", "glyph-test-src-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpSrc)

	tmpDst, err := os.MkdirTemp("", "glyph-test-dst-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDst)

	// Test restricted file (0600)
	restrictedFile := filepath.Join(tmpSrc, "restricted")
	if err := os.WriteFile(restrictedFile, []byte("secret"), 0600); err != nil {
		t.Fatal(err)
	}

	// Test executable file (0700)
	executableFile := filepath.Join(tmpSrc, "executable")
	if err := os.WriteFile(executableFile, []byte("#!/bin/sh"), 0700); err != nil {
		t.Fatal(err)
	}

	if err := copyDir(tmpSrc, tmpDst); err != nil {
		t.Fatalf("copyDir failed: %v", err)
	}

	// Check restricted file permissions
	info, err := os.Stat(filepath.Join(tmpDst, "restricted"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected permission 0600, got %o", info.Mode().Perm())
	}

	// Check executable file permissions
	info, err = os.Stat(filepath.Join(tmpDst, "executable"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("Expected permission 0700, got %o", info.Mode().Perm())
	}
}

func TestCopyFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := filepath.Join(tmpDir, "src")
	dst := filepath.Join(tmpDir, "dst")

	if err := os.WriteFile(src, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected permission 0600, got %o", info.Mode().Perm())
	}
}
