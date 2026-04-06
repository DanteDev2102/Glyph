package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestDeleteSection_ConfigProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-delete-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := []byte("[config]\nauthor = \"test-author\"\n\n[my-template]\nrepo = \"some-repo\"\n")
	err = os.WriteFile(configPath, initialContent, 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File:    configPath,
		Content: initialContent,
	}

	// Try to delete the 'config' section (should currently succeed, but we want it to fail after hardening)
	p.DeleteSection("config")

	// Read back the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]interface{}
	err = toml.Unmarshal(data, &config)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := config["config"]; !ok {
		t.Errorf("Config section was deleted! Protection is missing.")
	}
}

func TestDeleteSection_CaseInsensitiveConfigProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-delete-case-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := []byte("[config]\nauthor = \"test-author\"\n")
	err = os.WriteFile(configPath, initialContent, 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File:    configPath,
		Content: initialContent,
	}

	// Try to delete 'CONFIG' (case-insensitive)
	p.DeleteSection("CONFIG")

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]interface{}
	toml.Unmarshal(data, &config)

	if _, ok := config["config"]; !ok {
		t.Errorf("Config section was deleted using uppercase 'CONFIG'! Case-insensitive protection is missing.")
	}
}
