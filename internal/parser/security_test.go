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
	tmpDir, err := os.MkdirTemp("", "glyph-vulnerability-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := `[config]
author = "Original Author"

[test-template]
repo = "https://github.com/test/repo"
`
	err = os.WriteFile(configPath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File:    configPath,
		Content: []byte(initialContent),
	}

	// Attempt to delete the [config] section
	p.DeleteSection("config")

	// Read back the file
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]interface{}
	err = toml.Unmarshal(content, &config)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := config["config"]; !ok {
		t.Error("config section was deleted! (Vulnerable)")
	}
}

func TestWriteSection_ConfigProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-vulnerability-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := `[config]
author = "Original Author"

[test-template]
repo = "https://github.com/test/repo"
`
	err = os.WriteFile(configPath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File:    configPath,
		Content: []byte(initialContent),
	}

	// Attempt to modify the [config] section
	tmpl := &Template{
		Author: "Hacker",
	}
	p.WriteSection(tmpl, "config")

	// Read back the file
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]map[string]string
	err = toml.Unmarshal(content, &config)
	if err != nil {
		t.Fatal(err)
	}

	if c, ok := config["config"]; ok {
		if c["author"] == "Hacker" {
			t.Error("config section was modified! (Vulnerable)")
		}
	} else {
		t.Error("config section was removed or renamed! (Vulnerable)")
	}
}

func TestMkdirPermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-perm-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configDir := filepath.Join(tmpDir, "config")
	configPath := filepath.Join(configDir, "repositories.toml")

	p := &Parser{
		File: configPath,
	}

	tmpl := &Template{
		Name: "test-template",
		Repo: "https://github.com/test/repo",
	}

	// This should call os.MkdirAll(dir, 0700)
	p.Write(tmpl)

	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatal(err)
	}

	mode := info.Mode().Perm()
	if mode != 0700 {
		t.Errorf("Config directory has incorrect permissions: %v, expected 0700", mode)
	}
}
