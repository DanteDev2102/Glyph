package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir_PreservePermissions(t *testing.T) {
	tmpSrc, err := os.MkdirTemp("", "glyph-src-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpSrc)

	tmpDst, err := os.MkdirTemp("", "glyph-dst-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDst)

	// Create an executable file in the source
	exePath := filepath.Join(tmpSrc, "script.sh")
	err = os.WriteFile(exePath, []byte("#!/bin/sh\necho hello"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Copy the directory
	err = copyDir(tmpSrc, tmpDst)
	if err != nil {
		t.Errorf("copyDir failed: %v", err)
	}

	// Check permissions in the destination
	dstExePath := filepath.Join(tmpDst, "script.sh")
	info, err := os.Lstat(dstExePath)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode()&0111 == 0 {
		t.Errorf("Executable bit not preserved! Mode: %v", info.Mode())
	}
}

func TestReplaceInFile_PreservePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "README.md")
	// Create an executable README (unlikely but good for testing)
	err = os.WriteFile(filePath, []byte("Hello {{.ProjectName}}"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	replacements := map[string]string{"ProjectName": "TestProject"}
	replaceInFile(filePath, replacements)

	info, err := os.Lstat(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode()&0111 == 0 {
		t.Errorf("Executable bit not preserved after replacement! Mode: %v", info.Mode())
	}
}
