package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-perm-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	err = os.Mkdir(srcDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	exeFile := filepath.Join(srcDir, "script.sh")
	err = os.WriteFile(exeFile, []byte("#!/bin/sh\necho hello"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	copiedFile := filepath.Join(dstDir, "script.sh")
	info, err := os.Stat(copiedFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0755)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("Expected permissions %v, got %v", expectedPerm, info.Mode().Perm())
	}
}

func TestCopyFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-file-perm-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "script.sh")
	err = os.WriteFile(srcFile, []byte("#!/bin/sh\necho hello"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "script_copy.sh")
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0755)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("Expected permissions %v, got %v", expectedPerm, info.Mode().Perm())
	}
}
