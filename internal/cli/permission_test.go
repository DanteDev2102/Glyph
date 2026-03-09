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

	// 1. Test copyFile permission preservation
	srcFile := filepath.Join(tmpDir, "src_script.sh")
	// Create an executable file with restricted permissions
	err = os.WriteFile(srcFile, []byte("#!/bin/sh\necho hello"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "dst_script.sh")
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0700 {
		t.Errorf("copyFile: Expected permissions 0700, got %o", info.Mode().Perm())
	}

	// 2. Test replaceInFile permission preservation (behavioral check)
	readmeFile := filepath.Join(tmpDir, "README.md")
	err = os.WriteFile(readmeFile, []byte("Hello {{.Name}}"), 0400) // Read-only for owner
	if err != nil {
		t.Fatal(err)
	}

	replaceInFile(readmeFile, map[string]string{"Name": "Sentinel"})

	info, err = os.Stat(readmeFile)
	if err != nil {
		t.Fatal(err)
	}

	// On Unix, writing to an existing file should NOT change its permissions even if a different mode is passed to WriteFile.
	if info.Mode().Perm() != 0400 {
		t.Errorf("replaceInFile: Expected permissions 0400 to be preserved, got %o", info.Mode().Perm())
	}

	// 3. Test copyDir permission preservation
	srcDir := filepath.Join(tmpDir, "src_dir")
	err = os.Mkdir(srcDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	nestedFile := filepath.Join(srcDir, "nested.sh")
	err = os.WriteFile(nestedFile, []byte("echo nested"), 0711)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst_dir")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	info, err = os.Stat(filepath.Join(dstDir, "nested.sh"))
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0711 {
		t.Errorf("copyDir: Expected permissions 0711 for nested file, got %o", info.Mode().Perm())
	}
}
