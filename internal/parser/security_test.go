package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
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

func TestDeleteSection_ConfigProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := "[config]\nauthor = \"Sentinel\"\n"
	err = os.WriteFile(configPath, []byte(initialContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File:        configPath,
		Content:     []byte(initialContent),
		ContentRead: true,
	}

	p.DeleteSection("  CONFIG  ")

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var finalConfig map[string]interface{}
	err = toml.Unmarshal(data, &finalConfig)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := finalConfig["config"]; !ok {
		t.Error("The 'config' section was deleted!")
	}
}
