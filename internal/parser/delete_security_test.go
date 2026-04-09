package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestDeleteSection_ProtectConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-delete-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := `[config]
author = "Sentinel"

[test-template]
repo = "https://github.com/test/repo"
`
	err = os.WriteFile(configPath, []byte(initialContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File: configPath,
	}
	err = p.ExtractCommands()
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete the config section
	p.DeleteSection("config")

	// Verify if config section still exists
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	err = toml.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := data["config"]; !ok {
		t.Errorf("Critical vulnerability: [config] section was deleted!")
	}
}

func TestDeleteSection_ProtectConfigCaseInsensitive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-delete-ci-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := `[config]
author = "Sentinel"
`
	err = os.WriteFile(configPath, []byte(initialContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File: configPath,
	}
	err = p.ExtractCommands()
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete the config section with different casing
	p.DeleteSection("CONFIG")

	// Verify if config section still exists
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "[config]") && !strings.Contains(string(content), "[CONFIG]") {
		t.Errorf("Critical vulnerability: [config] section was deleted via case-insensitive match!")
	}
}
