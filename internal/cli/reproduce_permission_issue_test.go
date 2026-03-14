package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPermissionsPreservation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-perm-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")

	// 1. Test copyDir permissions preservation
	err = os.Mkdir(srcDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	restrictedFile := filepath.Join(srcDir, "restricted.sh")
	err = os.WriteFile(restrictedFile, []byte("#!/bin/sh\necho hi"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Join(dstDir, "restricted.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("copyDir did not preserve permissions: expected 0700, got %o", info.Mode().Perm())
	}

	// 2. Test copyFile permissions preservation
	srcLicense := filepath.Join(tmpDir, "LICENSE.test")
	dstLicense := filepath.Join(tmpDir, "LICENSE.dst")
	err = os.WriteFile(srcLicense, []byte("license content"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(srcLicense, dstLicense)
	if err != nil {
		t.Fatal(err)
	}

	info, err = os.Stat(dstLicense)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("copyFile did not preserve permissions: expected 0600, got %o", info.Mode().Perm())
	}

	// 3. Test replaceInFile permissions preservation
	readmePath := filepath.Join(tmpDir, "README.md")
	err = os.WriteFile(readmePath, []byte("Hello {{.Author}}"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"Author": "Sentinel"}
	replaceInFile(readmePath, replacements)

	info, err = os.Stat(readmePath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("replaceInFile did not preserve permissions: expected 0600, got %o", info.Mode().Perm())
	}
}
