package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite_Symlink(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	targetFile := filepath.Join(tmpDir, "sensitive_file")
	err = os.WriteFile(targetFile, []byte("SENSITIVE DATA"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create a symlink from configPath to targetFile
	err = os.Symlink(targetFile, configPath)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File: configPath,
	}

	tmpl := &Template{
		Name: "test-template",
		Repo: "https://github.com/test/repo",
	}

	// This should NOT overwrite targetFile if we have symlink protection
	p.Write(tmpl)

	// Check if targetFile was modified
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "SENSITIVE DATA" {
		t.Errorf("Target file was modified through symlink! Content: %s", string(content))
	}
}
