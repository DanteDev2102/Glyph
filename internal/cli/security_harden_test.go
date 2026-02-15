package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyFile_SymlinkOverwrite(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-security-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	secretFile := filepath.Join(tmpDir, "secret")
	err = os.WriteFile(secretFile, []byte("TOP SECRET"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "dst")
	err = os.Symlink(secretFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}

	srcFile := filepath.Join(tmpDir, "src")
	err = os.WriteFile(srcFile, []byte("NOT SECRET"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Should return an error now
	err = copyFile(srcFile, dstFile)
	if err == nil {
		t.Error("Expected error when copying to a symlink destination, but got nil")
	} else if !strings.Contains(err.Error(), "symbolic link") {
		t.Errorf("Expected security error about symbolic link, got: %v", err)
	}

	// Check if secretFile was NOT overwritten
	content, err := os.ReadFile(secretFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "TOP SECRET" {
		t.Errorf("Secret file was overwritten through symlink! Content: %s", string(content))
	}
}

func TestCopyDir_SymlinkSource(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-security-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	secretFile := filepath.Join(tmpDir, "secret")
	err = os.WriteFile(secretFile, []byte("TOP SECRET"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	srcDir := filepath.Join(tmpDir, "srcDir")
	err = os.Mkdir(srcDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	symlinkFile := filepath.Join(srcDir, "leak")
	err = os.Symlink(secretFile, symlinkFile)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dstDir")

	// Should skip the symlink
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatalf("copyDir failed: %v", err)
	}

	// Check if leaked file does NOT exist in dstDir
	leakedFile := filepath.Join(dstDir, "leak")
	if _, err := os.Stat(leakedFile); !os.IsNotExist(err) {
		t.Errorf("Symlink was followed or copied to destination!")
	}
}

func TestCopyDir_SymlinkDestination(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-security-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	secretDir := filepath.Join(tmpDir, "secretDir")
	err = os.Mkdir(secretDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dstDir")
	err = os.Mkdir(dstDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create a symlink in dstDir pointing to secretDir
	symlinkSubDir := filepath.Join(dstDir, "subdir")
	err = os.Symlink(secretDir, symlinkSubDir)
	if err != nil {
		t.Fatal(err)
	}

	srcDir := filepath.Join(tmpDir, "srcDir")
	err = os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(srcDir, "subdir", "evil.txt"), []byte("EVIL"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Should return an error now
	err = copyDir(srcDir, dstDir)
	if err == nil {
		t.Error("Expected error when copying to a symlink destination, but got nil")
	} else if !strings.Contains(err.Error(), "symbolic link") {
		t.Errorf("Expected security error about symbolic link, got: %v", err)
	}

	// Check if evil.txt was NOT created in secretDir
	evilFile := filepath.Join(secretDir, "evil.txt")
	if _, err := os.Stat(evilFile); err == nil {
		t.Errorf("File was created in secret directory through symlink! %s", evilFile)
	}
}
