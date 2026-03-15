package cli

import (
	"os"
	"path/filepath"
	"testing"
)

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

	executableFile := filepath.Join(srcDir, "run.sh")
	err = os.WriteFile(executableFile, []byte("#!/bin/sh\necho hi"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Join(dstDir, "run.sh"))
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0700 {
		t.Errorf("copyDir failed to preserve permissions: expected %o, got %o", 0700, info.Mode().Perm())
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
		t.Errorf("copyFile failed to preserve permissions: expected %o, got %o", 0600, info.Mode().Perm())
	}
}

func TestReplaceInFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	file := filepath.Join(tmpDir, "README.md")
	err = os.WriteFile(file, []byte("Project: {{.ProjectName}}"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"ProjectName": "Test"}
	replaceInFile(file, replacements)

	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("replaceInFile failed to preserve permissions: expected %o, got %o", 0600, info.Mode().Perm())
	}
}
