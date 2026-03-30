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

	executableFile := filepath.Join(srcDir, "script.sh")
	err = os.WriteFile(executableFile, []byte("#!/bin/sh\necho hello"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")

	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	dstExecutable := filepath.Join(dstDir, "script.sh")
	info, err := os.Lstat(dstExecutable)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode() != 0700 {
		t.Errorf("Permissions not preserved for executable file! Expected 0700, got %o", info.Mode())
	}
}

func TestReplaceInFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "script.sh")
	err = os.WriteFile(testFile, []byte("echo {{.Name}}"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"Name": "Sentinel"}
	replaceInFile(testFile, replacements)

	info, err := os.Lstat(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode() != 0700 {
		t.Errorf("Permissions not preserved after replacement! Expected 0700, got %o", info.Mode())
	}
}

func TestCopyFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src_script.sh")
	err = os.WriteFile(srcFile, []byte("#!/bin/sh"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "dst_script.sh")
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Lstat(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode() != 0700 {
		t.Errorf("Permissions not preserved by copyFile! Expected 0700, got %o", info.Mode())
	}
}
