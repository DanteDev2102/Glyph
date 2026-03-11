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

	// Create an executable file
	execFile := filepath.Join(srcDir, "exec_file")
	err = os.WriteFile(execFile, []byte("#!/bin/sh\necho hi"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create a read-only file
	readOnlyFile := filepath.Join(srcDir, "readonly_file")
	err = os.WriteFile(readOnlyFile, []byte("read only"), 0444)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	// Check permissions of executable file
	info, err := os.Stat(filepath.Join(dstDir, "exec_file"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0755 {
		t.Errorf("Expected 0755 perm for exec_file, got %o", info.Mode().Perm())
	}

	// Check permissions of read-only file
	info, err = os.Stat(filepath.Join(dstDir, "readonly_file"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0444 {
		t.Errorf("Expected 0444 perm for readonly_file, got %o", info.Mode().Perm())
	}
}

func TestCopyFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-fileperm-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src_exec")
	err = os.WriteFile(srcFile, []byte("exec content"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "dst_exec")
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dstFile)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("Expected 0700 perm for dst_exec, got %o", info.Mode().Perm())
	}
}

func TestReplaceInFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-replaceperm-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	targetFile := filepath.Join(tmpDir, "perm_check")
	err = os.WriteFile(targetFile, []byte("Hello {{.Name}}"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"Name": "Sentinel"}
	replaceInFile(targetFile, replacements)

	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "Hello Sentinel" {
		t.Errorf("Content mismatch: %s", string(content))
	}

	info, err := os.Stat(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected 0600 perm for perm_check, got %o", info.Mode().Perm())
	}
}
