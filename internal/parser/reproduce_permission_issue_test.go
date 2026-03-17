package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdirAll_RestrictedPermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	confDir := filepath.Join(tmpDir, ".config", "Glyph")
	p := &Parser{File: filepath.Join(confDir, "repositories.toml")}

	tmpl := &Template{
		Name: "test-tmpl",
		Repo: "https://github.com/example/repo",
	}

	// This will trigger MkdirAll
	p.Write(tmpl)

	info, err := os.Stat(confDir)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerm := os.FileMode(0700).Perm()
	actualPerm := info.Mode().Perm()

	if actualPerm != expectedPerm {
		t.Errorf("Config directory created with overly permissive permissions. Expected %v, got %v", expectedPerm, actualPerm)
	}
}
