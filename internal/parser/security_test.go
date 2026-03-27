package parser

import (
	"os"
	"path/filepath"
	"strings"
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

func TestDeleteSection_ConfigProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-delete-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	configContent := `[config]
author = "test-author"

[test-template]
repo = "https://github.com/test/repo"
`
	err = os.WriteFile(configPath, []byte(configContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File: configPath,
	}
	p.Refresh()

	// Attempt to delete the config section
	p.DeleteSection("config")

	// Verify the config section still exists
	p.Refresh()
	// Note: ExtractCommands handles [config] specially and puts it in p.Config,
	// but it shouldn't be in p.Commmands.
	// Let's check the raw content instead to be sure.

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "[config]") {
		t.Error("The [config] section was deleted despite protection!")
	}
}
