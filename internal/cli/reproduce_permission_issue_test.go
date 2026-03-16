package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDirPermissionPreservation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-copy-dir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	src := filepath.Join(tempDir, "src")
	dst := filepath.Join(tempDir, "dst")

	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(src, "exec.sh")
	if err := os.WriteFile(testFile, []byte("#!/bin/sh"), 0700); err != nil {
		t.Fatal(err)
	}

	if err := copyDir(src, dst); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Join(dst, "exec.sh"))
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0700)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("expected permissions %v, got %v", expectedPerm, info.Mode().Perm())
	}
}

func TestCopyFilePermissionPreservation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-copy-file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	src := filepath.Join(tempDir, "src.txt")
	dst := filepath.Join(tempDir, "dst.txt")

	if err := os.WriteFile(src, []byte("content"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := copyFile(src, dst); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0600)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("expected permissions %v, got %v", expectedPerm, info.Mode().Perm())
	}
}

func TestReplaceInFilePermissionPreservation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-replace")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello {{.Name}}"), 0600); err != nil {
		t.Fatal(err)
	}

	replaceInFile(testFile, map[string]string{"Name": "World"})

	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0600)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("expected permissions %v, got %v", expectedPerm, info.Mode().Perm())
	}
}
