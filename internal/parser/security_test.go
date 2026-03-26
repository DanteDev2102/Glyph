package parser

import (
	"os"
	"path/filepath"
	"strings"
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
	tmpDir, err := os.MkdirTemp("", "glyph-delete-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := `[config]
author = "test-author"

[template1]
repo = "https://github.com/test/repo"`
	err = os.WriteFile(configPath, []byte(initialContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File: configPath,
	}
	// Populate p.Content
	_, err = p.Read()
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to delete the "config" section
	p.DeleteSection("config")

	// Read content back
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "[config]") {
		t.Error("The [config] section was deleted, but it should be protected!")
	}
}

func TestWriteSection_ConfigProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-write-section-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := `[config]
author = "test-author"

[template1]
repo = "https://github.com/test/repo"`
	err = os.WriteFile(configPath, []byte(initialContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File: configPath,
	}
	// Populate p.Content
	_, err = p.Read()
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to modify the "config" section via WriteSection
	tmpl := &Template{
		Author: "malicious-author",
	}
	p.WriteSection(tmpl, "config")

	// Read content back and check author
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]interface{}
	err = toml.Unmarshal(content, &config)
	if err != nil {
		t.Fatal(err)
	}

	if configMap, ok := config["config"].(map[string]interface{}); ok {
		if author, ok := configMap["author"].(string); ok && author == "malicious-author" {
			t.Error("The [config] section was modified, but it should be protected from direct updates!")
		}
	}
}
