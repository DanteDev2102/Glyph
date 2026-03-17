package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplaceInFile_PreservesPermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "executable_script.sh")
	content := []byte("#!/bin/bash\necho Hello {{.Name}}")
	// Create with executable permissions
	err = os.WriteFile(testFile, content, 0755)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"Name": "Sentinel"}
	replaceInFile(testFile, replacements)

	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0755).Perm()
	actualPerm := info.Mode().Perm()

	if actualPerm != expectedPerm {
		t.Errorf("Permissions not preserved after replaceInFile. Expected %v, got %v", expectedPerm, actualPerm)
	}
}

func TestCopyFile_PreservesPermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src_executable")
	err = os.WriteFile(srcFile, []byte("content"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "dst_executable")
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0700).Perm()
	actualPerm := info.Mode().Perm()

	if actualPerm != expectedPerm {
		t.Errorf("Permissions not preserved after copyFile. Expected %v, got %v", expectedPerm, actualPerm)
	}
}

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

	srcFile := filepath.Join(srcDir, "exec")
	err = os.WriteFile(srcFile, []byte("content"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(dstDir, "exec")
	info, err := os.Stat(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0755).Perm()
	actualPerm := info.Mode().Perm()

	if actualPerm != expectedPerm {
		t.Errorf("Permissions not preserved after copyDir. Expected %v, got %v", expectedPerm, actualPerm)
	}
}
