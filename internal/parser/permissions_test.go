package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdirAll_ConfigSecurePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config", "repositories.toml")
	p := &Parser{
		File: configPath,
	}

	tmpl := &Template{
		Name: "test-template",
		Repo: "https://github.com/test/repo",
	}

	p.Write(tmpl)

	dirInfo, err := os.Stat(filepath.Join(tmpDir, "config"))
	if err != nil {
		t.Fatal(err)
	}

	if dirInfo.Mode().Perm() != 0700 {
		t.Errorf("Expected configuration directory permissions 0700, got %o", dirInfo.Mode().Perm())
	}
}
