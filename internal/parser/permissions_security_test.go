package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPermissions_FileHardening(t *testing.T) {
	tmpDir := t.TempDir()
	confFile := filepath.Join(tmpDir, "repositories.toml")

	// Pre-create file with broad permissions
	if err := os.WriteFile(confFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	p := &Parser{File: confFile}
	if err := p.safeWrite([]byte("test content")); err != nil {
		t.Fatalf("safeWrite failed: %v", err)
	}

	info, err := os.Stat(confFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestPermissions_DirectoryHardening(t *testing.T) {
	tmpDir := t.TempDir()
	confDir := filepath.Join(tmpDir, ".config", "Glyph")
	confFile := filepath.Join(confDir, "repositories.toml")

	if err := os.MkdirAll(confDir, 0755); err != nil {
		t.Fatal(err)
	}

	p := &Parser{File: confFile}
	p.Write(&Template{Name: "test"})

	info, err := os.Stat(confDir)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0700 {
		t.Errorf("expected 0700, got %v", info.Mode().Perm())
	}
}

func TestPermissions_NoChmodOnCWD(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalWd)

	// Current directory should have 0755 or similar (usually restricted by umask)
	// We'll set it to 0755 explicitly to test
	if err := os.Chmod(".", 0755); err != nil {
		t.Fatal(err)
	}

	p := &Parser{File: "repositories.toml"}
	p.Write(&Template{Name: "test"})

	info, _ := os.Stat(".")
	if info.Mode().Perm() == 0700 {
		t.Error("safeWrite should not chmod the current directory")
	}
}
