package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestWriteSection_ProtectConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-write-test-*")
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

	tmpl := &Template{
		Name:   "hacker-config",
		Author: "Malicious",
	}

	// Try to rename [config] to [hacker-config]
	p.WriteSection(tmpl, "config")

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
		t.Errorf("Critical vulnerability: [config] section was renamed or removed via WriteSection!")
	}

	// Check if it was overwritten with template fields (config should only have 'author' from the Config struct, not template fields like 'repo')
	// Wait, WriteSection unmarshals into map[string]map[string]string, which might break the [config] section structure anyway.
}

func TestWriteSection_ProtectConfigCaseInsensitive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-write-ci-test-*")
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

	tmpl := &Template{
		Author: "Malicious",
	}

	// Try to modify [config] using "CONFIG"
	p.WriteSection(tmpl, "CONFIG")

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
		// If it's not under "config", maybe it was renamed to something else or lost
		t.Errorf("Critical vulnerability: [config] section was modified/removed via case-insensitive WriteSection!")
	}
}
