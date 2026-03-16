package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDirectoryPermissions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-parser-perms")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	confFile := filepath.Join(tempDir, "config", "repositories.toml")
	p := &Parser{File: confFile}

	tmpl := &Template{
		Name: "test",
		Repo: "https://github.com/test/test",
	}

	p.Write(tmpl)

	dirInfo, err := os.Stat(filepath.Dir(confFile))
	if err != nil {
		t.Fatal(err)
	}

	expectedDirPerm := os.FileMode(0700)
	if dirInfo.Mode().Perm() != expectedDirPerm {
		t.Errorf("expected directory permissions %v, got %v", expectedDirPerm, dirInfo.Mode().Perm())
	}
}

func TestConfigFilePermissionsHardening(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-parser-file-perms")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	confFile := filepath.Join(tempDir, "repositories.toml")
	// Pre-create file with permissive permissions
	if err := os.WriteFile(confFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	p := &Parser{File: confFile}
	err = p.safeWrite([]byte("[config]\nauthor = \"test\""))
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(confFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0600)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("expected file permissions %v, got %v", expectedPerm, info.Mode().Perm())
	}
}
