package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPermissions_FileAndDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-perm-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configDir := filepath.Join(tmpDir, ".config", "Glyph")
	configPath := filepath.Join(configDir, "repositories.toml")

	p := &Parser{
		File: configPath,
	}

	tmpl := &Template{
		Name: "test-template",
		Repo: "https://github.com/test/repo",
	}

	// 1. Test Directory Permissions
	p.Write(tmpl)

	dirInfo, err := os.Stat(configDir)
	if err != nil {
		t.Fatal(err)
	}

	// Expecting 0700 (drwx------)
	if dirInfo.Mode().Perm() != 0700 {
		t.Errorf("Config directory has insecure permissions: %v, expected 0700", dirInfo.Mode().Perm())
	}

	// 2. Test File Permissions Hardening (Existing File)
	// Reset: Delete file and create with loose permissions
	os.Remove(configPath)
	err = os.WriteFile(configPath, []byte("existing data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Verify it is 0644
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo.Mode().Perm() != 0644 {
		t.Fatalf("Setup failed: file is not 0644, it is %v", fileInfo.Mode().Perm())
	}

	// Call Write (which calls safeWrite)
	p.Content = []byte("") // reset content to avoid unmarshal error or use Refresh/ensureContentRead
	p.ContentRead = false
	p.Write(tmpl)

	// Check permissions again
	fileInfo, err = os.Stat(configPath)
	if err != nil {
		t.Fatal(err)
	}

	// Expecting 0600 (rw-------)
	if fileInfo.Mode().Perm() != 0600 {
		t.Errorf("Config file permissions were not hardened: %v, expected 0600", fileInfo.Mode().Perm())
	}
}
